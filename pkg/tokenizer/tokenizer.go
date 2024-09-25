package tokenizer

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8" // Ajout de cet import
)

var (
	numberRegex = regexp.MustCompile(`^[0-9]+`)
	wordRegex   = regexp.MustCompile(`^[a-zA-Z]+`)
)

// CountTokens compte le nombre de tokens dans un texte donné
func CountTokens(text string) int {
	tokens := 0
	for len(text) > 0 {
		text = strings.TrimLeft(text, " \n\t\r")
		if len(text) == 0 {
			break
		}

		if numberRegex.MatchString(text) {
			match := numberRegex.FindString(text)
			tokens += len(match)
			text = text[len(match):]
		} else if wordRegex.MatchString(text) {
			match := wordRegex.FindString(text)
			tokens += len(match)
			text = text[len(match):]
		} else {
			r, size := utf8.DecodeRuneInString(text)
			if unicode.IsPunct(r) || unicode.IsSymbol(r) {
				tokens++
			} else {
				tokens += size
			}
			text = text[size:]
		}
	}
	return tokens
}

// EstimateTokens estime le nombre de tokens pour l'entrée et la sortie de la traduction
func EstimateTokens(input string) (inputTokens, outputTokens int) {
	inputTokens = CountTokens(input)
	// Supposons que la traduction produise environ 30% plus de tokens que l'entrée
	outputTokens = int(float64(inputTokens) * 1.3)
	return inputTokens, outputTokens
}
