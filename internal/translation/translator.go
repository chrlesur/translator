package translation

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/chrlesur/translator/internal/api"
	"github.com/chrlesur/translator/pkg/fileutils"
	"github.com/chrlesur/translator/pkg/logger"
	"github.com/pkoukk/tiktoken-go"
)

type Translator struct {
	APIClient  *api.ClaudeClient
	BatchSize  int
	NumThreads int
	Debug      bool
	SourceLang string
}

type BatchStatus struct {
	ID            int
	InputTokens   int
	OutputTokens  int
	Status        string
	ErrorOccurred bool
}

func NewTranslator(apiClient *api.ClaudeClient, batchSize, numThreads int, debug bool, sourceLang string) *Translator {
	if sourceLang == "" {
		sourceLang = "français" // Langue source par défaut
	}
	return &Translator{
		APIClient:  apiClient,
		BatchSize:  batchSize,
		NumThreads: numThreads,
		Debug:      debug,
		SourceLang: sourceLang,
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

	translatedContent := strings.Join(translatedBatches, "")

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

	t.printProgress(batchStatuses)

	for i, batch := range batches {
		wg.Add(1)
		go func(index int, batchContent string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			t.processSingleBatch(index, batchContent, targetLang, &batchStatuses[index], &translatedBatches[index])
			atomic.AddInt32(&completedBatches, 1)
			t.printProgress(batchStatuses)
		}(i, batch)
	}

	wg.Wait()
	fmt.Println() // Ajoute un saut de ligne après l'affichage de la progression

	t.printFinalStats(batchStatuses)

	return translatedBatches, nil
}

func (t *Translator) processSingleBatch(index int, batchContent, targetLang string, status *BatchStatus, result *string) {
	inputTokens := t.countTokens(batchContent)
	status.InputTokens = inputTokens
	status.Status = "Envoyé au LLM"

	translated, err := t.APIClient.Translate(batchContent, t.SourceLang, targetLang)
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
	var batches []string
	for len(content) > 0 {
		batchSize := t.BatchSize
		if batchSize > len(content) {
			batchSize = len(content)
		}
		batches = append(batches, content[:batchSize])
		content = content[batchSize:]
	}
	return batches
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

func (t *Translator) printProgress(statuses []BatchStatus) {
	completedCount := 0
	for _, status := range statuses {
		if status.Status == "Terminé" {
			completedCount++
		}
	}
	progress := float64(completedCount) / float64(len(statuses)) * 100

	fmt.Printf("\033[2K\r") // Effacer la ligne courante
	fmt.Printf("Progression : %.2f%% ", progress)

	for _, status := range statuses {
		var statusStr string
		switch status.Status {
		case "En attente":
			statusStr = "⏳"
		case "Envoyé au LLM":
			statusStr = fmt.Sprintf("🚀 (%d tokens)", status.InputTokens)
		case "Terminé":
			statusStr = fmt.Sprintf("✅ (%d → %d tokens)", status.InputTokens, status.OutputTokens)
		case "Erreur":
			statusStr = "❌"
		}
		fmt.Printf("| Batch %d: %s ", status.ID, statusStr)
	}
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
	encoding := "cl100k_base"  // Encodage utilisé par GPT-3.5 et GPT-4
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

	translated, err := t.APIClient.Translate(text, t.SourceLang, targetLang)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la traduction : %w", err)
	}

	return translated, nil
}
