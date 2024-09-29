# Translator

Translator est un outil en ligne de commande puissant pour traduire des documents texte, principalement en format Markdown, en utilisant différents moteurs d'IA tels que Claude, GPT-4, ou Ollama.

## Caractéristiques

- Support de multiples moteurs d'IA : Anthropic Claude, OpenAI GPT, et Ollama
- Traduction de fichiers entiers ou mode interactif pour des traductions rapides
- Préservation du formatage, y compris les sauts de ligne et l'espacement
- Gestion intelligente des lots pour optimiser les performances et respecter les limites des API
- Options de personnalisation avancées pour chaque traduction

## Installation

### Prérequis

Avant d'installer Translator, assurez-vous d'avoir Go installé sur votre système. Translator nécessite Go 1.16 ou une version ultérieure.

Pour vérifier votre version de Go, exécutez :
```
go version
```

Si Go n'est pas installé, vous pouvez le télécharger et l'installer depuis [le site officiel de Go](https://golang.org/dl/).

### Installation depuis les sources

Pour installer Translator, suivez ces étapes :

1. Clonez le dépôt :
   ```
   git clone https://github.com/chrlesur/translator.git
   ```

2. Naviguez dans le répertoire du projet :
   ```
   cd translator
   ```

3. Compilez et installez le projet :
   ```
   go install ./cmd/translator
   ```

Cette commande compilera le projet et placera l'exécutable dans votre `GOPATH/bin`.

### Configuration spécifique à chaque système d'exploitation

#### Windows

Assurez-vous que votre `GOPATH/bin` est dans votre PATH. Vous pouvez l'ajouter en exécutant cette commande dans PowerShell :
```
$env:Path += ";$env:GOPATH\bin"
```
Pour rendre ce changement permanent, ajoutez-le à votre profil PowerShell ou aux variables d'environnement du système.

#### Linux et macOS

Ajoutez la ligne suivante à votre fichier `.bashrc`, `.zshrc` ou équivalent :
```
export PATH=$PATH:$(go env GOPATH)/bin
```
Puis rechargez votre configuration de shell :
```
source ~/.bashrc  # ou ~/.zshrc, selon votre shell
```

### Vérification de l'installation

Pour vérifier que Translator est correctement installé, exécutez :
```
translator --version
```

Cela devrait afficher la version actuelle de Translator.

### Mise à jour

Pour mettre à jour Translator vers la dernière version, naviguez dans le répertoire du projet et exécutez :
```
git pull
go install ./cmd/translator
```

### Utilisation sans installation

Si vous préférez ne pas installer Translator globalement, vous pouvez le compiler et l'exécuter directement depuis le répertoire du projet :

1. Naviguez dans le répertoire du projet
2. Compilez le projet :
   ```
   go build -o translator ./cmd/translator
   ```
3. Exécutez le programme compilé :
   ```
   ./translator [commandes et options]
   ```

Cette méthode créera un exécutable nommé `translator` (ou `translator.exe` sur Windows) dans le répertoire courant.

## Utilisation

### Traduction de fichier

```
translator translate input.md EN --engine anthropic --model claude-3-sonnet-20240229
```

### Mode interactif

```
translator interactive --engine openai --model "gpt-4"
```

### Test de l'API

```
translator test-api --engine anthropic --model "claude-3-5-sonnet-20240620"
```

## Options

- `-d, --debug` : Active le mode debug
- `-b, --batch-size` : Nombre cible de tokens par lot pour le traitement (défaut: 1000)
- `-t, --threads` : Nombre de threads pour le traitement parallèle (défaut: 4)
- `-s, --source-lang` : Langue source du texte à traduire (défaut: français)
- `-i, --instruction` : Instruction complémentaire pour la traduction
- `-e, --engine` : Moteur de traduction (anthropic, openai, ollama)
- `-m, --model` : Modèle spécifique à utiliser pour le moteur choisi
- `--ollama-host` : Hôte Ollama (défaut: localhost)
- `--ollama-port` : Port Ollama (défaut: 11434)

## Configuration

Créez un fichier `.env` à la racine du projet avec les clés API nécessaires :

```
CLAUDE_API_KEY=votre_clé_claude
OPENAI_API_KEY=votre_clé_openai
```

## Développement

Pour contribuer au projet :

1. Forkez le dépôt
2. Créez votre branche de fonctionnalité (`git checkout -b feature/AmazingFeature`)
3. Committez vos changements (`git commit -m 'Add some AmazingFeature'`)
4. Poussez vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrez une Pull Request

## Structure du projet

```
translator/
│
├── cmd/
│   └── translator/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   ├── claude.go
│   │   ├── openai.go
│   │   └── ollama.go
│   │
│   ├── cli/
│   │   └── interactive.go
│   │
│   └── translation/
│       ├── translator.go
│       ├── utils.go
│       └── language-codes.go
│
├── pkg/
│   ├── fileutils/
│   │   └── fileutils.go
│   │
│   └── logger/
│       └── logger.go
│
├── .env
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Fonctionnement interne

Le traducteur fonctionne en plusieurs étapes :

1. Lecture du fichier d'entrée
2. Segmentation du contenu en lots gérables
3. Traduction parallèle des lots
4. Assemblage des traductions
5. Écriture du fichier de sortie

Le processus utilise un système de tokens pour optimiser l'utilisation des API de traduction et respecter leurs limites.

## Licence

Distribué sous la licence GPL-3.0. Voir `LICENSE` pour plus d'informations.

## Contact

Christophe Lesur - christophe.lesur@cloud-temple.com

Lien du projet : [https://github.com/chrlesur/translator](https://github.com/chrlesur/translator)
