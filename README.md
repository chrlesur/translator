# Translator

Translator est un outil en ligne de commande puissant pour traduire des documents texte, principalement en format Markdown, en utilisant l'API Claude 3.5 Sonnet d'Anthropic.

## Fonctionnalités

- Traduction de fichiers texte vers différentes langues
- Préservation du format Markdown pendant la traduction
- Support pour les gros fichiers avec traitement par lots
- Traitement parallèle pour des performances optimales
- Mode interactif pour la traduction à la volée
- Estimation précise du nombre de tokens et du coût de traduction
- Gestion de la clé API via variable d'environnement
- Mode debug pour afficher des informations détaillées

## Installation

1. Assurez-vous d'avoir Go installé sur votre système (version 1.16 ou supérieure).
2. Clonez ce dépôt :
   ```
   git clone https://github.com/votreuser/translator.git
   ```
3. Naviguez dans le répertoire du projet :
   ```
   cd translator
   ```
4. Installez les dépendances
   ```
   go mod tidy
   ```
5. Construisez l'application :
   ```
   go build -o translator cmd/translator/main.go
   ```

## Configuration

1. Créez un fichier `.env` à la racine du projet.
2. Ajoutez votre clé API Claude dans le fichier `.env` :
   ```
   CLAUDE_API_KEY=votre_clé_api_ici
   ```

## Utilisation

### Traduction de fichier

```
./translator translate chemin/vers/fichier.md langue_cible
```

### Mode interactif

```
./translator translate -i
```

### Options disponibles

- `-d, --debug` : Active le mode debug
- `-i, --interactive` : Active le mode interactif
- `-b, --batch-size` : Définit la taille des lots pour les gros fichiers (défaut : 1000)
- `-t, --threads` : Définit le nombre de threads pour le traitement parallèle (défaut : 4)

### Test de l'API

```
./translator test-api
```

## Développement

### Structure du projet

```
translator/
├── cmd/
│   └── translator/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── claude.go
│   ├── translation/
│   │   └── translator.go
│   └── cli/
│       └── interactive.go
├── pkg/
│   ├── fileutils/
│   │   └── fileutils.go
│   ├── markdown/
│   │   └── processor.go
│   └── tokenizer/
│       └── tokenizer.go
├── .gitignore
├── go.mod
└── README.md
```

### Exécution des tests

```
go test ./...
```

## Contribution

Les contributions sont les bienvenues ! Veuillez suivre ces étapes :

1. Forkez le projet
2. Créez votre branche de fonctionnalité (`git checkout -b feature/AmazingFeature`)
3. Committez vos changements (`git commit -m 'Add some AmazingFeature'`)
4. Poussez vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrez une Pull Request

