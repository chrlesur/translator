# Module Translation

Ce module contient le cœur de la logique de traduction du projet Translator.

## Fichiers

### translator.go

Ce fichier contient la structure principale `Translator` et ses méthodes associées.

Principales fonctionnalités :
- Gestion du processus de traduction de fichiers entiers
- Division du contenu en lots pour optimiser la traduction
- Traitement parallèle des lots
- Gestion de la progression et des statistiques de traduction

### utils.go

Ce fichier contient des fonctions utilitaires utilisées dans le processus de traduction.

Principales fonctionnalités :
- Initialisation de l'encodeur pour le comptage de tokens
- Fonction `SplitIntoSentences` pour diviser le texte en phrases
- Fonction `CountTokens` pour compter le nombre de tokens dans un texte
- Fonction `FormatProgress` pour formater l'affichage de la progression

### language-codes.go

Ce fichier gère les codes de langue utilisés dans le projet.

Principales fonctionnalités :
- Mapping entre les noms de langues et leurs codes ISO 639-1
- Fonctions pour obtenir le code à partir du nom de la langue et vice versa

## Utilisation

Le module translation est utilisé comme suit :

1. Créez une instance de `Translator` en utilisant `NewTranslator`
2. Appelez la méthode `TranslateFile` pour traduire un fichier entier
3. Utilisez `TranslateText` pour traduire des portions de texte individuelles

Exemple :
```go
translator := translation.NewTranslator(client, batchSize, numThreads, debug, sourceLang, additionalInstruction)
err := translator.TranslateFile(sourceFile, targetLang)
```

## Personnalisation

Pour personnaliser le comportement de la traduction :

- Modifiez les paramètres de `TargetBatchTokens` et `MaxBatchTokens` dans `NewTranslator`
- Ajustez la logique de division en lots dans `splitIntoBatches`
- Modifiez la fonction `FormatProgress` pour changer l'affichage de la progression

