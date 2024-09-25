package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

type MarkdownProcessor struct {
	placeholders map[string]string
	counter      int
}

func NewMarkdownProcessor() *MarkdownProcessor {
	return &MarkdownProcessor{
		placeholders: make(map[string]string),
		counter:      0,
	}
}

func (mp *MarkdownProcessor) PreprocessMarkdown(content string) string {
	// Préserver les blocs de code
	codeBlockRegex := regexp.MustCompile("```[\\s\\S]*?```")
	content = codeBlockRegex.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := mp.createPlaceholder()
		mp.placeholders[placeholder] = match
		return placeholder
	})

	// Préserver les liens
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	content = linkRegex.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := mp.createPlaceholder()
		mp.placeholders[placeholder] = match
		return placeholder
	})

	// Préserver les éléments en ligne (gras, italique, code inline)
	inlineRegex := regexp.MustCompile(`(\*\*|__|_|\*|` + "`" + `)(.+?)\1`)
	content = inlineRegex.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := mp.createPlaceholder()
		mp.placeholders[placeholder] = match
		return placeholder
	})

	return content
}

func (mp *MarkdownProcessor) PostprocessMarkdown(content string) string {
	for placeholder, original := range mp.placeholders {
		content = strings.Replace(content, placeholder, original, 1)
	}
	return content
}

func (mp *MarkdownProcessor) createPlaceholder() string {
	mp.counter++
	return fmt.Sprintf("__PLACEHOLDER_%d__", mp.counter)
}
