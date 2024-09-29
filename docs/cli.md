# Module CLI

Ce module gère l'interface en ligne de commande du Translator, en particulier le mode interactif.

## Fichiers

### interactive.go

Ce fichier implémente le mode interactif du Translator.

Principales fonctionnalités :
- Boucle d'interaction avec l'utilisateur pour des traductions à la volée
- Gestion des entrées/sorties pour la traduction interactive
- Commandes spéciales (comme '/quit' pour quitter)

## Utilisation

Le mode interactif est généralement lancé via la commande :

```
translator interactive [options]
```

Dans le code, il est utilisé comme suit :

```go
cli.RunInteractiveMode(translator)
```

## Personnalisation

Pour étendre les fonctionnalités du mode interactif :

1. Ajoutez de nouvelles commandes spéciales dans la boucle principale
2. Modifiez le formatage des entrées/sorties
3. Intégrez des fonctionnalités supplémentaires (par exemple, sauvegarde de l'historique)

## Intégration avec le reste du projet

Le mode interactif utilise l'instance de `Translator` créée dans `main.go` pour effectuer les traductions. Il s'appuie sur les fonctionnalités du module `translation` pour le traitement du texte.

