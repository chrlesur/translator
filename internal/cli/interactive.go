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

	fmt.Println("Mode interactif activé. Tapez '/quit' pour quitter.")
	fmt.Print("Langue cible : ")
	targetLang, _ := reader.ReadString('\n')
	targetLang = strings.TrimSpace(targetLang)

	for {
		fmt.Print("Texte à traduire : ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "/quit" {
			fmt.Println("Au revoir !")
			return
		}

		translated, err := translator.TranslateText(text, targetLang)
		if err != nil {
			fmt.Printf("Erreur lors de la traduction : %v\n", err)
			continue
		}

		fmt.Printf("Traduction : %s\n\n", translated)
	}
}
