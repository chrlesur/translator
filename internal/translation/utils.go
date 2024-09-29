package translation

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/chrlesur/translator/pkg/logger"
	"github.com/pkoukk/tiktoken-go"
)

var globalEncoder *tiktoken.Tiktoken

func InitializeEncoder() error {
	var err error
	logger.Info("Initialisation du tokenizer")
	globalEncoder, err = tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return fmt.Errorf("erreur lors de l'initialisation de l'encodeur : %v", err)
	}
	return nil
}

func SplitIntoSentences(text string) []string {
	var sentences []string
	var currentSentence strings.Builder
	inQuote := false

	for _, r := range text {
		currentSentence.WriteRune(r)

		if r == '"' {
			inQuote = !inQuote
		}

		if !inQuote && (r == '.' || r == '!' || r == '?') {
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
	abbreviations := []string{"Mr.", "Mrs.", "Dr.", "M.", "Prof.", "Sr.", "Jr.", "U.S.A.", "U.K.", "i.e.", "e.g."}
	for _, abbr := range abbreviations {
		if strings.HasSuffix(s, abbr) {
			return true
		}
	}
	return false
}

func CountTokens(text string) int {
	if globalEncoder == nil {
		logger.Error("L'encodeur n'a pas été initialisé")
		return len(strings.Fields(text))
	}
	return len(globalEncoder.Encode(text, nil, nil))
}

func FormatProgress(statuses []BatchStatus) string {
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
