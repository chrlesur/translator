package translation

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chrlesur/translator/pkg/fileutils"
	"github.com/chrlesur/translator/pkg/logger"
)

type TranslationClient interface {
	Translate(content, sourceLang, targetLang, additionalInstruction string) (string, error)
}

type Translator struct {
	Client                TranslationClient
	TargetBatchTokens     int
	MaxBatchTokens        int
	NumThreads            int
	Debug                 bool
	SourceLang            string
	AdditionalInstruction string
}

type BatchStatus struct {
	ID            int
	InputTokens   int
	OutputTokens  int
	Status        string
	ErrorOccurred bool
}

func NewTranslator(client TranslationClient, batchSize, numThreads int, debug bool, sourceLang, additionalInstruction string) *Translator {
	if sourceLang == "" {
		sourceLang = "français"
	}
	return &Translator{
		Client:                client,
		TargetBatchTokens:     batchSize,
		MaxBatchTokens:        batchSize * 2,
		NumThreads:            numThreads,
		Debug:                 debug,
		SourceLang:            sourceLang,
		AdditionalInstruction: additionalInstruction,
	}
}

func (t *Translator) TranslateFile(sourceFile, targetLang string) error {
	if err := InitializeEncoder(); err != nil {
		return fmt.Errorf("erreur lors de l'initialisation de l'encodeur : %w", err)
	}

	content, err := fileutils.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture du fichier : %w", err)
	}

	// Logging des informations de traduction
	sourceFileName := filepath.Base(sourceFile)
	sourceLanguageCode := GetCodeForLanguage(t.SourceLang)
	targetLanguageCode := GetCodeForLanguage(targetLang)
	logger.Info(fmt.Sprintf("Démarrage de la traduction - Fichier source: %s | Langue source: %s (%s) | Langue cible: %s (%s)",
		sourceFileName, t.SourceLang, sourceLanguageCode, targetLang, targetLanguageCode))

	batches := t.splitIntoBatches(content)
	logger.Info(fmt.Sprintf("Fichier divisé en %d lots", len(batches)))

	translatedBatches, err := t.processBatches(batches, targetLang)
	if err != nil {
		return err
	}

	translatedContent := strings.Join(translatedBatches, "\n\n")

	outputFile := t.generateOutputFilename(sourceFile, targetLang)
	err = fileutils.WriteFile(outputFile, translatedContent)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier traduit : %w", err)
	}

	logger.Info(fmt.Sprintf("Traduction terminée. Fichier de sortie : %s", outputFile))
	return nil
}

func (t *Translator) processBatches(batches []string, targetLang string) ([]string, error) {
	translatedBatches := make([]string, len(batches))
	batchStatuses := make([]BatchStatus, len(batches))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, t.NumThreads)
	var completedBatches int32

	for i := range batchStatuses {
		batchStatuses[i] = BatchStatus{ID: i + 1, Status: "En attente"}
	}

	var lastOutput string
	var outputMutex sync.Mutex

	updateProgress := func() {
		outputMutex.Lock()
		defer outputMutex.Unlock()

		newOutput := FormatProgress(batchStatuses)
		if newOutput != lastOutput {
			fmt.Print("\033[2K\r" + newOutput)
			lastOutput = newOutput
		}
	}

	if len(batchStatuses) > 0 {
		batchStatuses[0].Status = "Envoyé au LLM"
		updateProgress()
	}

	for i, batch := range batches {
		wg.Add(1)
		go func(index int, batchContent string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if index > 0 {
				batchStatuses[index].Status = "Envoyé au LLM"
				updateProgress()
			}

			t.processSingleBatch(index, batchContent, targetLang, &batchStatuses[index], &translatedBatches[index])
			atomic.AddInt32(&completedBatches, 1)
			updateProgress()
		}(i, batch)
	}

	stopTicker := make(chan struct{})
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				updateProgress()
			case <-stopTicker:
				return
			}
		}
	}()

	wg.Wait()
	close(stopTicker)

	updateProgress()
	fmt.Println()

	t.printFinalStats(batchStatuses)

	return translatedBatches, nil
}

func (t *Translator) processSingleBatch(index int, batchContent, targetLang string, status *BatchStatus, result *string) {
	inputTokens := CountTokens(batchContent)
	status.InputTokens = inputTokens
	status.Status = "Envoyé au LLM"

	translated, err := t.Client.Translate(batchContent, t.SourceLang, targetLang, t.AdditionalInstruction)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la traduction du lot %d : %v", index+1, err))
		status.Status = "Erreur"
		status.ErrorOccurred = true
		return
	}

	outputTokens := CountTokens(translated)
	*result = translated
	status.OutputTokens = outputTokens
	status.Status = "Terminé"

	if t.Debug {
		logger.Debug(fmt.Sprintf("Lot %d traduit avec succès", index+1))
	}
}

func (t *Translator) splitIntoBatches(content string) []string {
	lines := strings.Split(content, "\n")
	var batches []string
	var currentBatch strings.Builder
	currentTokens := 0

	addLineToBatch := func(line string) {
		if currentBatch.Len() > 0 {
			currentBatch.WriteString("\n")
		}
		currentBatch.WriteString(line)
		currentTokens += CountTokens(line)
	}

	finalizeBatch := func() {
		if currentBatch.Len() > 0 {
			batches = append(batches, currentBatch.String())
			currentBatch.Reset()
			currentTokens = 0
		}
	}

	for i, line := range lines {
		lineTokens := CountTokens(line)

		if currentTokens+lineTokens > t.MaxBatchTokens {
			finalizeBatch()
		}

		if lineTokens > t.MaxBatchTokens {
			sentences := SplitIntoSentences(line)
			for _, sentence := range sentences {
				sentenceTokens := CountTokens(sentence)
				if currentTokens+sentenceTokens > t.MaxBatchTokens {
					finalizeBatch()
				}
				addLineToBatch(sentence)
			}
		} else {
			addLineToBatch(line)
		}

		if currentTokens >= t.TargetBatchTokens {
			nextLine := ""
			if i+1 < len(lines) {
				nextLine = lines[i+1]
			}
			if strings.TrimSpace(nextLine) == "" || strings.HasPrefix(nextLine, "#") {
				finalizeBatch()
			}
		}
	}

	finalizeBatch()
	return batches
}

func (t *Translator) generateOutputFilename(sourceFile, targetLang string) string {
	dir, file := filepath.Split(sourceFile)
	ext := filepath.Ext(file)
	baseName := strings.TrimSuffix(file, ext)

	targetCode := GetCodeForLanguage(targetLang)
	if targetCode == "" {
		logger.Warning(fmt.Sprintf("Code de langue non trouvé pour : %s. Utilisation du nom de la langue comme code.", targetLang))
		targetCode = targetLang
	}

	return filepath.Join(dir, fmt.Sprintf("%s_%s%s", baseName, targetCode, ext))
}

func (t *Translator) printFinalStats(statuses []BatchStatus) {
	var totalInputTokens, totalOutputTokens int
	for _, status := range statuses {
		totalInputTokens += status.InputTokens
		totalOutputTokens += status.OutputTokens
	}
	logger.Info(fmt.Sprintf("Statistiques finales :"))
	logger.Info(fmt.Sprintf("Total des tokens en entrée : %d", totalInputTokens))
	logger.Info(fmt.Sprintf("Total des tokens en sortie : %d", totalOutputTokens))
}

func (t *Translator) TranslateText(text, targetLang string) (string, error) {
	if t.Debug {
		logger.Debug(fmt.Sprintf("Traduction du texte : %s", text))
	}

	translated, err := t.Client.Translate(text, t.SourceLang, targetLang, t.AdditionalInstruction)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la traduction : %w", err)
	}

	return translated, nil
}
