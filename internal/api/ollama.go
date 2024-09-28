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

type OllamaClient struct {
	host  string
	port  string
	model string
	debug bool
}

func NewOllamaClient(host, port, model string, debug bool) *OllamaClient {
	if model == "" {
		model = "llama2" // Modèle par défaut
	}
	return &OllamaClient{
		host:  host,
		port:  port,
		model: model,
		debug: debug,
	}
}

func (c *OllamaClient) Translate(content, sourceLang, targetLang, additionalInstruction string) (string, error) {
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

	if c.debug {
		logger.Debug(fmt.Sprintf("Envoi de la requête à Ollama. Modèle : %s, Prompt : %s", c.model, prompt))
	}

	requestBody, _ := json.Marshal(map[string]string{
		"model":  c.model,
		"prompt": prompt,
	})

	maxRetries := 5
	baseTimeout := 10 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		timeout := baseTimeout + time.Duration(retry*20)*time.Second
		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Post(fmt.Sprintf("http://%s:%s/api/generate", c.host, c.port),
			"application/json", bytes.NewBuffer(requestBody))

		if err != nil {
			if retry < maxRetries-1 {
				if c.debug {
					logger.Debug(fmt.Sprintf("Tentative %d échouée. Timeout après %v. Nouvelle tentative...", retry+1, timeout))
				}
				continue
			}
			return "", fmt.Errorf("erreur lors de l'envoi de la requête à l'API Ollama après %d tentatives : %w", maxRetries, err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("erreur lors de la lecture de la réponse : %w", err)
		}

		if c.debug {
			logger.Debug(fmt.Sprintf("Réponse brute de Ollama : %s", string(body)))
		}

		if resp.StatusCode != http.StatusOK {
			if retry < maxRetries-1 {
				if c.debug {
					logger.Debug(fmt.Sprintf("Tentative %d échouée. Statut HTTP : %d. Nouvelle tentative...", retry+1, resp.StatusCode))
				}
				continue
			}
			return "", fmt.Errorf("erreur de l'API Ollama après %d tentatives. Dernier statut : %d, Corps : %s", maxRetries, resp.StatusCode, string(body))
		}

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return "", fmt.Errorf("erreur lors du décodage de la réponse JSON : %w", err)
		}

		translatedText, ok := result["response"].(string)
		if !ok {
			return "", fmt.Errorf("aucun contenu valide dans la réponse de l'API Ollama")
		}

		if c.debug {
			logger.Debug(fmt.Sprintf("Contenu de la réponse Ollama : %s", translatedText))
		}

		return translatedText, nil
	}

	return "", fmt.Errorf("échec de la traduction après %d tentatives", maxRetries)
}
