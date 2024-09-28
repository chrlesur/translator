package main

import (
	"fmt"
	"os"

	"github.com/chrlesur/translator/internal/api"
	"github.com/chrlesur/translator/internal/cli"
	"github.com/chrlesur/translator/internal/translation"
	"github.com/chrlesur/translator/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const VERSION = "1.0.0"

var (
	debug                 bool
	batchSize             int
	numThreads            int
	sourceLang            string
	additionalInstruction string
	engine                string
	model                 string
	ollamaHost            string
	ollamaPort            string
	contextSize           int
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Active le mode debug")
	rootCmd.PersistentFlags().IntVarP(&batchSize, "batch-size", "b", 1000, "Nombre cible de tokens par lot pour le traitement")
	rootCmd.PersistentFlags().IntVarP(&numThreads, "threads", "t", 4, "Nombre de threads pour le traitement parallèle")
	rootCmd.PersistentFlags().StringVarP(&sourceLang, "source-lang", "s", "français", "Langue source du texte à traduire")
	rootCmd.PersistentFlags().StringVarP(&additionalInstruction, "instruction", "i", "", "Instruction complémentaire pour la traduction")
	rootCmd.PersistentFlags().StringVarP(&engine, "engine", "e", "anthropic", "Moteur de traduction (anthropic, openai, ollama)")
	rootCmd.PersistentFlags().StringVarP(&model, "model", "m", "claude-3-5-sonnet-20240620", "Modèle spécifique à utiliser pour le moteur choisi")
	rootCmd.PersistentFlags().StringVar(&ollamaHost, "ollama-host", "localhost", "Hôte Ollama")
	rootCmd.PersistentFlags().StringVar(&ollamaPort, "ollama-port", "11434", "Port Ollama")
	rootCmd.PersistentFlags().IntVarP(&contextSize, "context-size", "c", 0, "Taille du contexte pour le modèle (0 pour utiliser la valeur par défaut)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(translateCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.AddCommand(testAPICmd)
}

var rootCmd = &cobra.Command{
	Use:   "translator",
	Short: "Translator est un outil de traduction de documents utilisant divers moteurs d'IA",
	Long: `Translator est un outil en ligne de commande pour traduire des documents texte, 
principalement en format Markdown, en utilisant différents moteurs d'IA comme Claude, GPT-4, ou Ollama.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.SetDebugMode(debug)
		logger.Info(fmt.Sprintf("Translator version %s", VERSION))
		if debug {
			logger.Debug("Mode debug activé")
		}
		err := godotenv.Load()
		if err != nil {
			logger.Error("Erreur lors du chargement du fichier .env")
		}
	},
}

var translateCmd = &cobra.Command{
	Use:   "translate [fichier_source] [langue_cible]",
	Short: "Traduit un fichier vers la langue cible spécifiée",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sourceFile := args[0]
		targetLang := args[1]

		client := getTranslationClient()
		if client == nil {
			return
		}

		translator := translation.NewTranslator(client, batchSize, numThreads, debug, sourceLang, additionalInstruction)

		err := translator.TranslateFile(sourceFile, targetLang)
		if err != nil {
			logger.Error(fmt.Sprintf("Erreur lors de la traduction : %v", err))
			return
		}
	},
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Lance le mode interactif pour des traductions à la volée",
	Run: func(cmd *cobra.Command, args []string) {
		client := getTranslationClient()
		if client == nil {
			return
		}

		translator := translation.NewTranslator(client, batchSize, numThreads, debug, sourceLang, additionalInstruction)

		cli.RunInteractiveMode(translator)
	},
}

var testAPICmd = &cobra.Command{
	Use:   "test-api",
	Short: "Teste la connexion à l'API de traduction",
	Run: func(cmd *cobra.Command, args []string) {
		client := getTranslationClient()
		if client == nil {
			return
		}

		_, err := client.Translate("Ceci est un test.", sourceLang, "anglais", "")
		if err != nil {
			logger.Error(fmt.Sprintf("Erreur lors du test de l'API : %v", err))
			return
		}
		logger.Info(fmt.Sprintf("La connexion à l'API %s est opérationnelle.", engine))
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Affiche la version du logiciel",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Translator version %s\n", VERSION)
	},
}

func getTranslationClient() translation.TranslationClient {
	switch engine {
	case "anthropic":
		apiKey := os.Getenv("CLAUDE_API_KEY")
		if apiKey == "" {
			logger.Error("La clé API Claude n'est pas définie dans la variable d'environnement CLAUDE_API_KEY")
			return nil
		}
		logger.Debug(fmt.Sprintf("Utilisation du moteur Anthropic avec la clé API : %s", apiKey[:5]+"..."))
		return api.NewClaudeClient(apiKey, model, debug, getContextSize("anthropic"))
	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			logger.Error("La clé API OpenAI n'est pas définie dans la variable d'environnement OPENAI_API_KEY")
			return nil
		}
		logger.Debug(fmt.Sprintf("Utilisation du moteur OpenAI avec la clé API : %s", apiKey[:5]+"..."))
		logger.Debug(fmt.Sprintf("Modèle sélectionné : %s", model))
		return api.NewOpenAIClient(apiKey, model, debug, getContextSize("openai"))
	case "ollama":
		logger.Debug(fmt.Sprintf("Utilisation d'Ollama avec l'hôte %s et le port %s", ollamaHost, ollamaPort))
		return api.NewOllamaClient(ollamaHost, ollamaPort, model, debug, getContextSize("ollama"))
	default:
		logger.Error(fmt.Sprintf("Moteur non reconnu : %s", engine))
		return nil
	}
}

func getContextSize(engineType string) int {
	if contextSize > 0 {
		return contextSize
	}
	switch engineType {
	case "anthropic", "openai":
		return 4000
	case "ollama":
		return 2000
	default:
		return 2000
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(fmt.Sprintf("Erreur : %v", err))
		os.Exit(1)
	}
}
