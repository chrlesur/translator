package translation

import (
	"fmt"
	"path/filepath"
	"strings"

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

	translatedContent, err := t.TranslateText(content, targetLang)
	if err != nil {
		return fmt.Errorf("erreur lors de la traduction : %w", err)
	}

	outputFile := t.GenerateOutputFilename(sourceFile, targetLang)
	err = fileutils.WriteFile(outputFile, translatedContent)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier traduit : %w", err)
	}

	fmt.Printf("Traduction terminée. Fichier de sortie : %s\n", outputFile)
	return nil
}

func (t *Translator) TranslateText(text, targetLang string) (string, error) {
	if t.Debug {
		fmt.Printf("Traduction du texte vers %s\n", targetLang)
	}

	segments := t.MarkdownProcessor.ProcessMarkdown(text)

	for i, segment := range segments {
		if segment.Translatable && strings.TrimSpace(segment.Content) != "" {
			translated, err := t.APIClient.Translate(segment.Content, targetLang)
			if err != nil {
				return "", fmt.Errorf("erreur lors de la traduction du segment %d : %w", i, err)
			}
			segments[i].Content = translated
		}
	}

	translatedText := t.MarkdownProcessor.ReassembleMarkdown(segments)

	return translatedText, nil
}

func (t *Translator) GenerateOutputFilename(sourceFile, targetLang string) string {
	dir, file := filepath.Split(sourceFile)
	ext := filepath.Ext(file)
	baseName := strings.TrimSuffix(file, ext)
	return filepath.Join(dir, fmt.Sprintf("%s_%s%s", baseName, targetLang, ext))
}
