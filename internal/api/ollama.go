package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chrlesur/translator/pkg/logger"
)

type OllamaClient struct {
	host        string
	port        string
	model       string
	debug       bool
	contextSize int
}

func NewOllamaClient(host, port, model string, debug bool, contextSize int) *OllamaClient {
	if model == "" {
		model = "llama2" // Modèle par défaut
	}
	if contextSize <= 0 {
		contextSize = 2000 // Valeur par défaut si non spécifié
	}
	return &OllamaClient{
		host:        host,
		port:        port,
		model:       model,
		debug:       debug,
		contextSize: contextSize,
	}
}

func (c *OllamaClient) Translate(content, sourceLang, targetLang, additionalInstruction string) (string, error) {
	prompt := fmt.Sprintf(`Translate the following text from %s to %s. Follow these rules strictly:
1. Translate accurately and professionally.
2. Preserve all formatting, including line breaks and spacing.
3. Do not translate URLs starting with http:// or https://.
4. Do not translate content within square brackets [ ].
5. Do not translate content between parentheses () immediately following a closing square bracket ].
6. Do not translate markdown filenames with .md extension.
7. Use academic language in your translation.
8. Do not add any comments, introductions, or explanations.
9. Your response must contain only the translated text, nothing else.
%s

Translate the above text now:

%s`, sourceLang, targetLang, additionalInstruction, content)

	if c.debug {
		logger.Debug(fmt.Sprintf("Envoi de la requête à Ollama. Prompt : %s", prompt))
	}

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":  c.model,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"num_ctx": c.contextSize,
		},
	})

	url := fmt.Sprintf("http://%s:%s/api/generate", c.host, c.port)
	if c.debug {
		logger.Debug(fmt.Sprintf("URL de l'API Ollama : %s", url))
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'envoi de la requête à l'API Ollama : %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la lecture de la réponse : %w", err)
	}

	if c.debug {
		logger.Debug(fmt.Sprintf("Réponse brute de Ollama : %s", string(body)))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("erreur lors du décodage de la réponse JSON : %w", err)
	}

	response, ok := result["response"].(string)
	if !ok {
		return "", fmt.Errorf("la réponse ne contient pas de champ 'response' valide")
	}

	if c.debug {
		logger.Debug(fmt.Sprintf("Contenu de la réponse Ollama : %s", response))
	}

	return response, nil
}
