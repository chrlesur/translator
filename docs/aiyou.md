# Module AI.YOU

Ce module contient l'implémentation du client AI.YOU pour le projet Translator.

## Fichier

### aiyou.go

Ce fichier implémente le client pour l'API AI.YOU.

Principales fonctionnalités :
- Initialisation du client avec l'ID de l'assistant et le mode debug
- Authentification avec email et mot de passe
- Méthode `Translate` pour envoyer des requêtes de traduction à l'API AI.YOU
- Gestion des threads et des runs pour le processus de traduction
- Gestion des erreurs et des tentatives de reconnexion

## Structure principale

### AIYOUClient

```go
type AIYOUClient struct {
    Token       string
    AssistantID string
    Debug       bool
    Timeout     time.Duration
}
```

- `Token` : Token d'authentification pour l'API AI.YOU
- `AssistantID` : ID de l'assistant AI.YOU à utiliser
- `Debug` : Active ou désactive le mode debug
- `Timeout` : Durée maximale d'attente pour les requêtes API

## Fonctions principales

### NewAIYOUClient

```go
func NewAIYOUClient(assistantID string, debug bool) *AIYOUClient
```

Crée une nouvelle instance du client AI.YOU.

### Login

```go
func (c *AIYOUClient) Login(email, password string) error
```

Authentifie le client avec l'email et le mot de passe fournis.

### Translate

```go
func (c *AIYOUClient) Translate(content, sourceLang, targetLang, additionalInstruction string) (string, error)
```

Traduit le contenu donné de la langue source vers la langue cible.

## Fonctions internes

### createThread

```go
func (c *AIYOUClient) createThread() (string, error)
```

Crée un nouveau thread de conversation.

### addMessage

```go
func (c *AIYOUClient) addMessage(threadID, content string) error
```

Ajoute un message au thread spécifié.

### createRun

```go
func (c *AIYOUClient) createRun(threadID string) (string, error)
```

Crée un nouveau run pour le thread spécifié.

### retrieveRun

```go
func (c *AIYOUClient) retrieveRun(threadID, runID string) (map[string]interface{}, error)
```

Récupère le statut d'un run en cours.

### waitForCompletion

```go
func (c *AIYOUClient) waitForCompletion(threadID, runID string) (*map[string]interface{}, error)
```

Attend la fin d'un run et retourne le résultat.

### makeAPICall

```go
func (c *AIYOUClient) makeAPICall(endpoint, method string, data []byte) ([]byte, error)
```

Effectue un appel à l'API AI.YOU avec gestion des erreurs et logging.

## Utilisation

Pour utiliser le client AI.YOU :

1. Créez une instance de `AIYOUClient` avec `NewAIYOUClient`
2. Appelez `Login` pour authentifier le client
3. Utilisez `Translate` pour effectuer des traductions

Exemple :

```go
client := api.NewAIYOUClient("assistant_id", true)
err := client.Login("email@example.com", "password")
if err != nil {
    log.Fatal(err)
}
translation, err := client.Translate("Bonjour", "français", "anglais", "")
if err != nil {
    log.Fatal(err)
}
fmt.Println(translation)
```

Note : Assurez-vous d'avoir configuré les variables d'environnement `AIYOU_EMAIL` et `AIYOU_PASSWORD` dans votre fichier `.env` avant d'utiliser ce client.
