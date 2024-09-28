package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/chrlesur/translator/pkg/logger"
)

const (
	ClaudeAPIURL = "https://api.anthropic.com/v1/messages"
)

type ClaudeClient struct {
	APIKey      string
	Model       string
	Debug       bool
	ContextSize int
}

type ClaudeRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func NewClaudeClient(apiKey, model string, debug bool, contextSize int) *ClaudeClient {
	if model == "" {
		model = "claude-3-5-sonnet-20240620" // Modèle par défaut
	}
	if contextSize <= 0 {
		contextSize = 4000 // Valeur par défaut si non spécifié
	}
	return &ClaudeClient{
		APIKey:      apiKey,
		Model:       model,
		Debug:       debug,
		ContextSize: contextSize,
	}
}

func (c *ClaudeClient) Translate(content, sourceLang, targetLang, additionalInstruction string) (string, error) {
	prompt := fmt.Sprintf(`You are a translation AI specializing in information technology content. Your task is to translate the following text from %s to %s.

Important instructions:
1. Translate the content accurately and professionally.
2. Preserve all formatting, including line breaks and spacing.
3. Do not translate URLs starting with http:// or https://.
4. Do not translate content within square brackets [ ].
5. Do not translate content between parentheses () immediately following a closing square bracket ].
6. Do not translate markdown filenames with .md extension.
7. Use academic language in your translation.
8. Do not add any comments, introductions, or explanations to your translation.
9. Provide only the translated text in your response, nothing else.
%s

Here's the text to translate:

%s

Remember, your response should contain only the translated text, with no additional comments or explanations.`, sourceLang, targetLang, additionalInstruction, content)

	if c.Debug {
		logger.Debug(fmt.Sprintf("Envoi de la requête à Claude. Modèle : %s, Prompt : %s", c.Model, prompt))
	}

	request := ClaudeRequest{
		Model: c.Model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 8000,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création de la requête JSON : %w", err)
	}

	baseTimeout := 10 * time.Second
	maxRetries := 5

	for retry := 0; retry < maxRetries; retry++ {
		timeout := baseTimeout + time.Duration(retry*20)*time.Second
		client := &http.Client{
			Timeout: timeout,
		}

		req, err := http.NewRequest("POST", ClaudeAPIURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", fmt.Errorf("erreur lors de la création de la requête HTTP : %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", c.APIKey)
		req.Header.Set("anthropic-version", "2023-06-01")

		resp, err := client.Do(req)
		if err != nil {
			if retry < maxRetries-1 {
				if c.Debug {
					logger.Debug(fmt.Sprintf("Tentative %d échouée. Timeout après %v. Nouvelle tentative...", retry+1, timeout))
				}
				continue
			}
			return "", fmt.Errorf("erreur lors de l'envoi de la requête à l'API Claude après %d tentatives : %w", maxRetries, err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("erreur lors de la lecture de la réponse : %w", err)
		}

		if c.Debug {
			logger.Debug(fmt.Sprintf("Réponse brute de l'API Claude : %s", string(body)))
		}

		if resp.StatusCode != http.StatusOK {
			if retry < maxRetries-1 {
				if c.Debug {
					logger.Debug(fmt.Sprintf("Tentative %d échouée. Statut HTTP : %d. Nouvelle tentative...", retry+1, resp.StatusCode))
				}
				continue
			}
			return "", fmt.Errorf("erreur de l'API Claude après %d tentatives. Dernier statut : %d, Corps : %s", maxRetries, resp.StatusCode, string(body))
		}

		var claudeResp ClaudeResponse
		err = json.Unmarshal(body, &claudeResp)
		if err != nil {
			return "", fmt.Errorf("erreur lors du décodage de la réponse JSON : %w", err)
		}

		if len(claudeResp.Content) == 0 {
			return "", fmt.Errorf("aucun contenu dans la réponse de l'API")
		}

		if c.Debug {
			logger.Debug(fmt.Sprintf("Contenu de la réponse Claude : %s", claudeResp.Content[0].Text))
		}

		return claudeResp.Content[0].Text, nil
	}

	return "", fmt.Errorf("échec de la traduction après %d tentatives", maxRetries)
}
