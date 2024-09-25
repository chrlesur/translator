package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chrlesur/translator/pkg/tokenizer"
)

const (
	ClaudeAPIURL = "https://api.anthropic.com/v1/chat"
)

type ClaudeClient struct {
	APIKey string
	Debug  bool
}

type ClaudeRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func NewClaudeClient(apiKey string, debug bool) *ClaudeClient {
	return &ClaudeClient{
		APIKey: apiKey,
		Debug:  debug,
	}
}

func (c *ClaudeClient) Translate(content, targetLang string) (string, error) {
	prompt := fmt.Sprintf("Translate the following text to %s:\n\n%s", targetLang, content)

	request := ClaudeRequest{
		Model: "claude-3-sonnet-20240229",
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création de la requête JSON : %w", err)
	}

	req, err := http.NewRequest("POST", ClaudeAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création de la requête HTTP : %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'envoi de la requête à l'API Claude : %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la lecture de la réponse : %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erreur de l'API Claude : %s", string(body))
	}

	var claudeResp ClaudeResponse
	err = json.Unmarshal(body, &claudeResp)
	if err != nil {
		return "", fmt.Errorf("erreur lors du décodage de la réponse JSON : %w", err)
	}

	if len(claudeResp.Choices) == 0 {
		return "", fmt.Errorf("aucune traduction reçue de l'API Claude")
	}

	return claudeResp.Choices[0].Message.Content, nil
}

func (c *ClaudeClient) EstimateTranslationCost(content string) (float64, int, int) {
	inputTokens, outputTokens := tokenizer.EstimateTokens(content)
	totalTokens := inputTokens + outputTokens

	// Prix pour Claude 3 Sonnet (à ajuster selon les tarifs réels)
	const pricePerMillionTokens = 15.0 // $15 par million de tokens
	cost := float64(totalTokens) * pricePerMillionTokens / 1000000.0

	return cost, inputTokens, outputTokens
}
