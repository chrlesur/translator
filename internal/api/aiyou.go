package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/chrlesur/translator/pkg/logger"
)

const AIYOUAPIURL = "https://ai.dragonflygroup.fr/api"

type AIYOUClient struct {
	Token       string
	AssistantID string
	Debug       bool
	Timeout     time.Duration
}

func NewAIYOUClient(assistantID string, debug bool) *AIYOUClient {
	return &AIYOUClient{
		AssistantID: assistantID,
		Debug:       debug,
		Timeout:     120 * time.Second,
	}
}

func (c *AIYOUClient) Login(email, password string) error {
	logger.Info("Tentative de connexion à AI.YOU")
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la création de la requête JSON : %v", err))
		return fmt.Errorf("erreur lors de la création de la requête JSON : %w", err)
	}

	resp, err := c.makeAPICall("/login", "POST", jsonData)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la connexion : %v", err))
		return err
	}

	var loginResp struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	}
	err = json.Unmarshal(resp, &loginResp)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors du décodage de la réponse JSON : %v", err))
		return fmt.Errorf("erreur lors du décodage de la réponse JSON : %w", err)
	}

	c.Token = loginResp.Token
	logger.Info("Connexion à AI.YOU réussie")
	return nil
}

func (c *AIYOUClient) Translate(content, sourceLang, targetLang, additionalInstruction string) (string, error) {
	logger.Debug(fmt.Sprintf("Démarrage de la traduction avec AI.YOU - Langue source: %s, Langue cible: %s", sourceLang, targetLang))

	threadID, err := c.createThread()
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la création du thread : %v", err))
		return "", err
	}
	logger.Debug(fmt.Sprintf("Thread créé avec l'ID : %s", threadID))

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
		logger.Debug(fmt.Sprintf("Prompt de traduction : %s", prompt))
	}

	err = c.addMessage(threadID, prompt)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de l'ajout du message au thread : %v", err))
		return "", err
	}
	logger.Debug("Message ajouté au thread avec succès")

	runID, err := c.createRun(threadID)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la création du run : %v", err))
		return "", err
	}
	logger.Debug(fmt.Sprintf("Run créé avec l'ID : %s", runID))

	completedRun, err := c.waitForCompletion(threadID, runID)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de l'attente de la complétion du run : %v", err))
		return "", err
	}
	logger.Debug("Run complété avec succès")

	response, ok := (*completedRun)["response"].(string)
	if !ok {
		logger.Error("La réponse n'a pas pu être extraite du run")
		return "", fmt.Errorf("la réponse n'a pas pu être extraite du run")
	}

	logger.Debug("Traduction AI.YOU terminée avec succès")
	if c.Debug {
		logger.Debug(fmt.Sprintf("Texte traduit : %s", response))
	}

	return response, nil
}

func (c *AIYOUClient) createThread() (string, error) {
	logger.Debug("Création d'un nouveau thread AI.YOU")
	resp, err := c.makeAPICall("/v1/threads", "POST", []byte("{}"))
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la création du thread : %v", err))
		return "", err
	}

	var threadResp struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(resp, &threadResp)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors du décodage de la réponse JSON : %v", err))
		return "", fmt.Errorf("erreur lors du décodage de la réponse JSON : %w", err)
	}

	if threadResp.ID == "" {
		logger.Error("L'ID du thread est vide dans la réponse")
		return "", fmt.Errorf("l'ID du thread est vide dans la réponse")
	}

	logger.Debug(fmt.Sprintf("Thread créé avec l'ID : %s", threadResp.ID))
	return threadResp.ID, nil
}

func (c *AIYOUClient) addMessage(threadID, content string) error {
	logger.Debug(fmt.Sprintf("Ajout d'un message au thread %s", threadID))
	messageData := map[string]string{
		"role":    "user",
		"content": content,
	}
	jsonData, err := json.Marshal(messageData)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la création de la requête JSON : %v", err))
		return fmt.Errorf("erreur lors de la création de la requête JSON : %w", err)
	}

	_, err = c.makeAPICall(fmt.Sprintf("/v1/threads/%s/messages", threadID), "POST", jsonData)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de l'ajout du message : %v", err))
		return err
	}
	logger.Debug("Message ajouté avec succès")
	return nil
}

func (c *AIYOUClient) createRun(threadID string) (string, error) {
	logger.Debug(fmt.Sprintf("Création d'un run pour le thread %s", threadID))
	runData := map[string]string{
		"assistantId": c.AssistantID,
	}
	jsonData, err := json.Marshal(runData)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la création de la requête JSON : %v", err))
		return "", fmt.Errorf("erreur lors de la création de la requête JSON : %w", err)
	}

	resp, err := c.makeAPICall(fmt.Sprintf("/v1/threads/%s/runs", threadID), "POST", jsonData)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la création du run : %v", err))
		return "", err
	}

	var runResp struct {
		ID string `json:"id"`
	}
	err = json.Unmarshal(resp, &runResp)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors du décodage de la réponse JSON : %v", err))
		return "", fmt.Errorf("erreur lors du décodage de la réponse JSON : %w", err)
	}

	logger.Debug(fmt.Sprintf("Run créé avec l'ID : %s", runResp.ID))
	return runResp.ID, nil
}

func (c *AIYOUClient) retrieveRun(threadID, runID string) (map[string]interface{}, error) {
	logger.Debug(fmt.Sprintf("Récupération du run %s pour le thread %s", runID, threadID))

	// Créons un corps de requête vide
	emptyBody := map[string]string{}
	jsonBody, err := json.Marshal(emptyBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la création du corps de la requête : %v", err))
		return nil, err
	}

	resp, err := c.makeAPICall(fmt.Sprintf("/v1/threads/%s/runs/%s", threadID, runID), "POST", jsonBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la récupération du run : %v", err))
		return nil, err
	}

	var runStatus map[string]interface{}
	err = json.Unmarshal(resp, &runStatus)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors du décodage de la réponse JSON : %v", err))
		return nil, fmt.Errorf("erreur lors du décodage de la réponse JSON : %w", err)
	}

	logger.Debug(fmt.Sprintf("Statut du run récupéré : %v", runStatus))
	return runStatus, nil
}

func (c *AIYOUClient) waitForCompletion(threadID, runID string) (*map[string]interface{}, error) {
	maxAttempts := 30
	delayBetweenAttempts := 2 * time.Second

	for i := 0; i < maxAttempts; i++ {
		logger.Debug(fmt.Sprintf("Tentative %d de récupération du statut du run", i+1))
		run, err := c.retrieveRun(threadID, runID)
		if err != nil {
			logger.Error(fmt.Sprintf("Erreur lors de la récupération du run : %v", err))
			return nil, err
		}

		status, ok := run["status"].(string)
		if !ok {
			logger.Error("Statut du run non trouvé ou invalide")
			return nil, fmt.Errorf("statut du run non trouvé ou invalide")
		}

		logger.Debug(fmt.Sprintf("Statut actuel du run : %s", status))

		if status == "completed" {
			logger.Debug("Run complété avec succès")
			return &run, nil
		}

		if status == "failed" || status == "cancelled" {
			logger.Error(fmt.Sprintf("Le run a échoué avec le statut : %s", status))
			return nil, fmt.Errorf("le run a échoué avec le statut: %s", status)
		}

		logger.Debug(fmt.Sprintf("En attente de la complétion du run. Pause de %v", delayBetweenAttempts))
		time.Sleep(delayBetweenAttempts)
	}

	logger.Error("Délai d'attente dépassé pour la complétion du run")
	return nil, fmt.Errorf("délai d'attente dépassé pour la complétion du run")
}

func (c *AIYOUClient) makeAPICall(endpoint, method string, data []byte) ([]byte, error) {
	url := AIYOUAPIURL + endpoint
	logger.Debug(fmt.Sprintf("Préparation de l'appel API vers %s avec la méthode %s", url, method))

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la création de la requête HTTP : %v", err))
		return nil, fmt.Errorf("erreur lors de la création de la requête HTTP : %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	// Log détaillé de la requête
	logHeaders := make(map[string]string)
	for k, v := range req.Header {
		logHeaders[k] = strings.Join(v, ", ")
	}
	logger.Debug(fmt.Sprintf("Détails de la requête : Method: %s, URL: %s, Headers: %v, Body: %s",
		method, url, logHeaders, string(data)))

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de l'envoi de la requête à l'API AI.YOU : %v", err))
		return nil, fmt.Errorf("erreur lors de l'envoi de la requête à l'API AI.YOU : %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Erreur lors de la lecture de la réponse : %v", err))
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse : %w", err)
	}

	// Log détaillé de la réponse
	logger.Debug(fmt.Sprintf("Réponse reçue : Status: %d, Body: %s", resp.StatusCode, string(body)))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger.Error(fmt.Sprintf("Erreur API (%d): %s", resp.StatusCode, string(body)))
		return nil, fmt.Errorf("erreur API (%d): %s", resp.StatusCode, string(body))
	}

	logger.Debug("Appel API réussi")
	return body, nil
}
