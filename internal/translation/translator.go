package translation

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chrlesur/translator/internal/api"
	"github.com/chrlesur/translator/pkg/fileutils"
	"github.com/chrlesur/translator/pkg/logger"
	"github.com/pkoukk/tiktoken-go"
)

const (
	targetBatchTokens = 1000 // Nombre cible de tokens par lot
	maxBatchTokens    = 2000 // Nombre maximum de tokens par lot
)

type Translator struct {
	APIClient             *api.ClaudeClient
	BatchSize             int
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

func NewTranslator(apiClient *api.ClaudeClient, batchSize, numThreads int, debug bool, sourceLang, additionalInstruction string) *Translator {
	if sourceLang == "" {
		sourceLang = "français" // Langue source par défaut
	}
	return &Translator{
		APIClient:             apiClient,
		BatchSize:             batchSize,
		NumThreads:            numThreads,
		Debug:                 debug,
		SourceLang:            sourceLang,
		AdditionalInstruction: additionalInstruction,
	}
}

func (t *Translator) TranslateFile(sourceFile, targetLang string) error {
	content, err := fileutils.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture du fichier : %w", err)
	}

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

		newOutput := t.formatProgress(batchStatuses)
		if newOutput != lastOutput {
			fmt.Print("\033[2K\r" + newOutput)
			lastOutput = newOutput
		}
	}

	// Mettre à jour le statut du premier lot avant de commencer
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

	// Démarrer une goroutine pour mettre à jour périodiquement la progression
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

	// Affichage final de la progression
	updateProgress()
	fmt.Println() // Ajoute un saut de ligne après l'affichage final

	t.printFinalStats(batchStatuses)

	return translatedBatches, nil
}

func (t *Translator) processSingleBatch(index int, batchContent, targetLang string, status *BatchStatus, result *string) {
	inputTokens := t.countTokens(batchContent)
	status.InputTokens = inputTokens
	status.Status = "Envoyé au LLM"

	translated, err := t.APIClient.Translate(batchContent, t.SourceLang, targetLang, t.AdditionalInstruction)

	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la traduction du lot %d : %v", index+1, err))
		status.Status = "Erreur"
		status.ErrorOccurred = true
		return
	}

	outputTokens := t.countTokens(translated)
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

	encoder, err := tiktoken.GetEncoding("cl100k_base")
	countTokens := func(s string) int {
		if err != nil {
			return len(strings.Fields(s))
		}
		return len(encoder.Encode(s, nil, nil))
	}

	addLineToBatch := func(line string) {
		if currentBatch.Len() > 0 {
			currentBatch.WriteString("\n")
		}
		currentBatch.WriteString(line)
		currentTokens += countTokens(line)
	}

	finalizeBatch := func() {
		if currentBatch.Len() > 0 {
			batches = append(batches, currentBatch.String())
			currentBatch.Reset()
			currentTokens = 0
		}
	}

	for i, line := range lines {
		lineTokens := countTokens(line)

		if currentTokens+lineTokens > maxBatchTokens {
			finalizeBatch()
		}

		if lineTokens > maxBatchTokens {
			// Split very long lines
			sentences := splitIntoSentences(line)
			for _, sentence := range sentences {
				sentenceTokens := countTokens(sentence)
				if currentTokens+sentenceTokens > maxBatchTokens {
					finalizeBatch()
				}
				addLineToBatch(sentence)
			}
		} else {
			addLineToBatch(line)
		}

		// Check if we should end the batch here
		if currentTokens >= targetBatchTokens {
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

func splitIntoSentences(text string) []string {
	var sentences []string
	var currentSentence strings.Builder
	inQuote := false

	for _, r := range text {
		currentSentence.WriteRune(r)

		if r == '"' {
			inQuote = !inQuote
		}

		if !inQuote && (r == '.' || r == '!' || r == '?') {
			// Check if it's really the end of a sentence (e.g., not "Mr." or "U.S.A.")
			if len(sentences) > 0 || !isAbbreviation(currentSentence.String()) {
				sentences = append(sentences, strings.TrimSpace(currentSentence.String()))
				currentSentence.Reset()
			}
		}
	}

	if currentSentence.Len() > 0 {
		sentences = append(sentences, strings.TrimSpace(currentSentence.String()))
	}

	return sentences
}

func isAbbreviation(s string) bool {
	// This is a simple check. You might want to expand this list or use a more sophisticated method
	abbreviations := []string{"Mr.", "Mrs.", "Dr.", "Prof.", "Sr.", "Jr.", "U.S.A.", "U.K.", "i.e.", "e.g."}
	for _, abbr := range abbreviations {
		if strings.HasSuffix(s, abbr) {
			return true
		}
	}
	return false
}

func (t *Translator) generateOutputFilename(sourceFile, targetLang string) string {
	dir, file := filepath.Split(sourceFile)
	ext := filepath.Ext(file)
	baseName := strings.TrimSuffix(file, ext)

	targetCode := GetCodeForLanguage(targetLang)
	if targetCode == "" {
		logger.Warning(fmt.Sprintf("Code pays non trouvé pour la langue : %s. Utilisation de la langue comme code.", targetLang))
		targetCode = targetLang
	}

	return filepath.Join(dir, fmt.Sprintf("%s_%s%s", baseName, targetCode, ext))
}

func (t *Translator) formatProgress(statuses []BatchStatus) string {
	completedCount := 0
	inProgressCount := 0
	for _, status := range statuses {
		if status.Status == "Terminé" {
			completedCount++
		} else if status.Status == "Envoyé au LLM" {
			inProgressCount++
		}
	}
	progress := float64(completedCount) / float64(len(statuses)) * 100

	const progressBarWidth = 20
	completedWidth := int(progress / 100 * progressBarWidth)
	progressBar := strings.Repeat("█", completedWidth) + strings.Repeat("░", progressBarWidth-completedWidth)

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Progression : [%s] %.2f%% | Terminés : %d/%d | En cours : %d",
		progressBar, progress, completedCount, len(statuses), inProgressCount))

	lastInProgress := []string{}
	for i := len(statuses) - 1; i >= 0 && len(lastInProgress) < 3; i-- {
		if statuses[i].Status == "Envoyé au LLM" {
			lastInProgress = append([]string{fmt.Sprintf("Lot %d (%d tokens)", statuses[i].ID, statuses[i].InputTokens)}, lastInProgress...)
		}
	}
	if len(lastInProgress) > 0 {
		buffer.WriteString(" | En traitement : " + strings.Join(lastInProgress, ", "))
	}

	return buffer.String()
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

func (t *Translator) countTokens(text string) int {
	encoding := "cl100k_base" // Encodage utilisé par GPT-3.5 et GPT-4
	tkm, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de l'initialisation de l'encodeur : %v", err))
		return 0
	}
	tokens := tkm.Encode(text, nil, nil)
	return len(tokens)
}

func (t *Translator) TranslateText(text, targetLang string) (string, error) {
	if t.Debug {
		logger.Debug(fmt.Sprintf("Traduction du texte : %s", text))
	}

	translated, err := t.APIClient.Translate(text, t.SourceLang, targetLang, t.AdditionalInstruction)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la traduction : %w", err)
	}

	return translated, nil
}
