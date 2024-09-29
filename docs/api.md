# Module API

Ce module contient les implémentations des différents clients d'API de traduction utilisés par le Translator.

## Fichiers

### claude.go

Ce fichier implémente le client pour l'API Claude d'Anthropic.

Principales fonctionnalités :
- Initialisation du client avec la clé API
- Méthode `Translate` pour envoyer des requêtes de traduction à l'API Claude
- Gestion des erreurs et des tentatives de reconnexion

### openai.go

Ce fichier implémente le client pour l'API GPT d'OpenAI.

Principales fonctionnalités :
- Initialisation du client avec la clé API
- Méthode `Translate` pour envoyer des requêtes de traduction à l'API GPT
- Gestion des modèles spécifiques (GPT-3.5, GPT-4, etc.)

### ollama.go

Ce fichier implémente le client pour l'API Ollama, permettant l'utilisation de modèles de langue en local.

Principales fonctionnalités :
- Initialisation du client avec l'hôte et le port Ollama
- Méthode `Translate` pour envoyer des requêtes de traduction à Ollama
- Gestion des différents modèles disponibles localement

**Attention à ne pas utiliser des batchs trop gros avec ollama et des modeles comme llama3.2. Limitez vous à 200 tokens par batch environ.
**
## Utilisation

Chaque client implémente l'interface `TranslationClient` définie dans `translator.go`. Pour utiliser un client spécifique :

1. Initialisez le client avec les paramètres appropriés (clé API, hôte, etc.)
2. Passez le client à la fonction `NewTranslator` lors de la création d'une instance de `Translator`

Exemple :
```go
client := api.NewClaudeClient(apiKey, model, debug, contextSize)
translator := translation.NewTranslator(client, batchSize, numThreads, debug, sourceLang, additionalInstruction)
```

## Ajout d'un nouveau client

Pour ajouter un nouveau client d'API de traduction :

1. Créez un nouveau fichier (par exemple `newapi.go`) dans ce dossier
2. Implémentez la structure du client et la méthode `Translate`
3. Assurez-vous que le client implémente l'interface `TranslationClient`
4. Ajoutez la logique nécessaire dans `main.go` pour permettre la sélection de ce nouveau client

