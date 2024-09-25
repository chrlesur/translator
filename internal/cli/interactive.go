package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chrlesur/translator/internal/translation"
)

func RunInteractiveMode(translator *translation.Translator) {
	reader := bufio.NewReader(os.Stdin)
	targetLang := "fr" // Langue par défaut

	fmt.Println("Mode interactif activé. Tapez '/help' pour voir les commandes disponibles.")

	for {
		fmt.Printf("Langue cible actuelle : %s\n", targetLang)
		fmt.Print("Texte à traduire : ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		switch text {
		case "/quit":
			fmt.Println("Au revoir !")
			return
		case "/help":
			printHelp()
		case "/lang":
			targetLang = changeLang(reader)
		default:
			translateAndPrint(translator, text, targetLang)
		}
	}
}

func printHelp() {
	fmt.Println("\nCommandes disponibles :")
	fmt.Println("/quit - Quitter le mode interactif")
	fmt.Println("/help - Afficher cette aide")
	fmt.Println("/lang - Changer la langue cible")
	fmt.Println("")
}

func changeLang(reader *bufio.Reader) string {
	fmt.Print("Nouvelle langue cible : ")
	newLang, _ := reader.ReadString('\n')
	return strings.TrimSpace(newLang)
}

func translateAndPrint(translator *translation.Translator, text, targetLang string) {
	cost, inputTokens, outputTokens := translator.APIClient.EstimateTranslationCost(text)
	fmt.Printf("Estimation : %d tokens en entrée, %d tokens en sortie, coût : $%.4f\n", inputTokens, outputTokens, cost)

	translated, err := translator.TranslateText(text, targetLang)
	if err != nil {
		fmt.Printf("Erreur lors de la traduction : %v\n", err)
		return
	}

	fmt.Printf("Traduction : %s\n\n", translated)
}
