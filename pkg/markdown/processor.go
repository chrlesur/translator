package markdown

import (
	"regexp"
	"strings"
)

type Segment struct {
	Content      string
	Translatable bool
}

type MarkdownProcessor struct{}

func NewMarkdownProcessor() *MarkdownProcessor {
	return &MarkdownProcessor{}
}

func (mp *MarkdownProcessor) ProcessMarkdown(content string) []Segment {
	var segments []Segment

	// Regex pour identifier les éléments non traduisibles
	codeBlockRegex := regexp.MustCompile("(?s)^```[\\s\\S]*?^```")
	inlineCodeRegex := regexp.MustCompile("`[^`\n]+`")
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	headingRegex := regexp.MustCompile(`^#{1,6}\s.*$`)

	// Diviser le contenu en lignes
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if codeBlockRegex.MatchString(line) {
			// Traiter les blocs de code multi-lignes
			codeBlock := mp.extractCodeBlock(lines[i:])
			segments = append(segments, Segment{Content: codeBlock, Translatable: false})
			i += strings.Count(codeBlock, "\n")
		} else if headingRegex.MatchString(line) {
			// Traiter les titres
			segments = append(segments, Segment{Content: line + "\n", Translatable: false})
		} else if strings.TrimSpace(line) == "" {
			// Préserver les lignes vides
			segments = append(segments, Segment{Content: "\n", Translatable: false})
		} else {
			// Traiter le texte normal, le code inline et les liens
			parts := mp.splitLine(line, inlineCodeRegex, linkRegex)
			segments = append(segments, parts...)
			// Ajouter un retour à la ligne après chaque ligne non vide
			segments = append(segments, Segment{Content: "\n", Translatable: false})
		}
	}

	return segments
}

func (mp *MarkdownProcessor) extractCodeBlock(lines []string) string {
	var codeBlock strings.Builder
	for _, line := range lines {
		codeBlock.WriteString(line + "\n")
		if strings.TrimSpace(line) == "```" {
			break
		}
	}
	return codeBlock.String()
}

func (mp *MarkdownProcessor) splitLine(line string, regexes ...*regexp.Regexp) []Segment {
	var segments []Segment
	lastIndex := 0

	for i := 0; i < len(line); i++ {
		for _, regex := range regexes {
			if loc := regex.FindStringIndex(line[i:]); loc != nil {
				start, end := i+loc[0], i+loc[1]
				if start > lastIndex {
					segments = append(segments, Segment{Content: line[lastIndex:start], Translatable: true})
				}
				segments = append(segments, Segment{Content: line[start:end], Translatable: false})
				lastIndex = end
				i = end - 1
				break
			}
		}
	}

	if lastIndex < len(line) {
		segments = append(segments, Segment{Content: line[lastIndex:], Translatable: true})
	}

	return segments
}

func (mp *MarkdownProcessor) ReassembleMarkdown(segments []Segment) string {
	var result strings.Builder
	for _, segment := range segments {
		result.WriteString(segment.Content)
	}
	return result.String()
}
