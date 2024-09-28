package api

import (
	"context"
	"fmt"
	"time"

	"github.com/chrlesur/translator/pkg/logger"
	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
	model  string
	debug  bool
}

func NewOpenAIClient(apiKey, model string, debug bool) *OpenAIClient {
	if model == "" {
		model = "gpt-4o-mini" // Modèle par défaut
	}
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
		model:  model,
		debug:  debug,
	}
}

func (c *OpenAIClient) Translate(content, sourceLang, targetLang, additionalInstruction string) (string, error) {
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
		logger.Debug(fmt.Sprintf("Envoi de la requête à OpenAI. Modèle : %s, Prompt : %s", c.model, prompt))
	}

	maxRetries := 5
	baseTimeout := 10 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		timeout := baseTimeout + time.Duration(retry*20)*time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		resp, err := c.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: c.model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)

		if err != nil {
			if retry < maxRetries-1 {
				if c.debug {
					logger.Debug(fmt.Sprintf("Tentative %d échouée. Timeout après %v. Nouvelle tentative...", retry+1, timeout))
				}
				continue
			}
			return "", fmt.Errorf("erreur lors de l'envoi de la requête à l'API OpenAI après %d tentatives : %w", maxRetries, err)
		}

		if c.debug {
			logger.Debug(fmt.Sprintf("Réponse brute de OpenAI : %+v", resp))
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("aucun contenu dans la réponse de l'API OpenAI")
		}

		translatedText := resp.Choices[0].Message.Content

		if c.debug {
			logger.Debug(fmt.Sprintf("Contenu de la réponse OpenAI : %s", translatedText))
		}

		return translatedText, nil
	}

	return "", fmt.Errorf("échec de la traduction après %d tentatives", maxRetries)
}
