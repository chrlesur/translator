package translation

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chrlesur/translator/internal/api"
	"github.com/chrlesur/translator/pkg/fileutils"
	"github.com/chrlesur/translator/pkg/markdown"
)

type Translator struct {
	APIClient         *api.ClaudeClient
	BatchSize         int
	NumThreads        int
	Debug             bool
	MarkdownProcessor *markdown.MarkdownProcessor
}

func NewTranslator(apiClient *api.ClaudeClient, batchSize, numThreads int, debug bool) *Translator {
	return &Translator{
		APIClient:         apiClient,
		BatchSize:         batchSize,
		NumThreads:        numThreads,
		Debug:             debug,
		MarkdownProcessor: markdown.NewMarkdownProcessor(),
	}
}

func (t *Translator) TranslateFile(sourceFile, targetLang string) error {
	content, err := fileutils.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture du fichier : %w", err)
	}

	preprocessed := t.MarkdownProcessor.PreprocessMarkdown(content)
	batches := t.splitIntoBatches(preprocessed)

	translatedBatches := make([]string, len(batches))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, t.NumThreads)

	for i, batch := range batches {
		wg.Add(1)
		go func(index int, batchContent string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			translated, err := t.APIClient.Translate(batchContent, targetLang)
			if err != nil {
				fmt.Printf("Erreur lors de la traduction du lot %d : %v\n", index, err)
				return
			}
			translatedBatches[index] = translated

			if t.Debug {
				fmt.Printf("Lot %d traduit avec succès\n", index)
			}
		}(i, batch)
	}

	wg.Wait()

	translatedContent := strings.Join(translatedBatches, "")
	postprocessed := t.MarkdownProcessor.PostprocessMarkdown(translatedContent)

	outputFile := t.generateOutputFilename(sourceFile, targetLang)
	err = fileutils.WriteFile(outputFile, postprocessed)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier traduit : %w", err)
	}

	fmt.Printf("Traduction terminée. Fichier de sortie : %s\n", outputFile)
	return nil
}

func (t *Translator) TranslateText(text, targetLang string) (string, error) {
	if t.Debug {
		fmt.Printf("Traduction du texte : %s\n", text)
	}

	preprocessed := t.MarkdownProcessor.PreprocessMarkdown(text)
	translated, err := t.APIClient.Translate(preprocessed, targetLang)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la traduction : %w", err)
	}

	postprocessed := t.MarkdownProcessor.PostprocessMarkdown(translated)
	return postprocessed, nil
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
	return filepath.Join(dir, fmt.Sprintf("%s_%s%s", baseName, targetLang, ext))
}
