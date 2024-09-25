# Translator

Translator is a powerful command-line tool for translating text documents, primarily in Markdown format, using Anthropic's Claude 3.5 Sonnet API.

## Fonctionnalités

- Translation of text files into different languages
- Preservation of Markdown format during translation
Support for large files with batch processing
Parallel processing for optimal performance
- Interactive mode for on-the-fly translation
- Accurate estimation of token count and translation cost
- API key management via environment variable
- Debug mode to display detailed information

## Installation

1. Make sure you have Go installed on your system (version 1.16 or higher).
2. Clone this repository:
I apologize, but there is no text provided in your message for me to translate. If you'd like a translation, please provide the original text and I'll be happy to translate it to English for you.
git clone https://github.com/youruser/translator.git
I apologize, but there is no text provided in your message to translate. If you'd like me to translate something, please provide the original text and I'll be happy to help.
3. Navigate to the project directory:
I apologize, but there is no text provided in the message for me to translate. If you would like me to translate something, please provide the text and specify the source language, and I'll be happy to translate it to English for you.
change directory to translator
I apologize, but there is no text provided in your request to translate. If you'd like me to translate something, please provide the original text and I'll be happy to translate it to English for you.
4. Install the dependencies
I apologize, but there is no text provided in your message for me to translate. If you'd like me to translate something, please provide the original text and I'll be happy to assist you with the translation to English.
go mod tidy
I apologize, but there is no text provided in the message for me to translate. The message contains only empty code block markers (```). If you'd like me to translate something, please provide the text you want translated within the code block or as part of your message.
5. Build the application:
You haven't provided any text to translate. If you'd like me to translate something, please provide the text and specify the source language, and I'll be happy to translate it to English for you.
go build -o translator cmd/translator/main.go
I apologize, but there is no text provided in your message for me to translate. If you would like me to translate something, please provide the text you want translated and I'll be happy to assist you.

## Configuration

1. Create a file`.env`at the root of the project.
2. Add your Claude API key to the file`.env`:
I apologize, but there is no text provided in your message for me to translate. If you'd like me to translate something, please include the original text and specify the source language. I'll be happy to translate it to English for you once you provide the text.
CLAUDE_API_KEY=your_api_key_here
I apologize, but there is no text provided in the message for me to translate. The message contains only empty quotation marks. If you'd like me to translate something, please provide the text you want translated within the quotation marks or in your message, and I'll be happy to assist you.

## Utilisation

### Traduction de fichier

I apologize, but there is no text provided in your message to translate. If you'd like me to translate something, please provide the original text and specify the source language, and I'll be happy to translate it to English for you.
./translator translate path/to/file.md target_language
I apologize, but there is no text provided in your message for me to translate. If you'd like me to translate something, please provide the original text and I'll be happy to translate it to English for you.

### Mode interactif

I apologize, but there is no text provided in your request to translate. If you would like me to translate something, please provide the original text and specify which language it is in, and I'll be happy to translate it to English for you.
I apologize, but there is no text provided after the command "./translator translate -i" for me to translate. If you have a specific text you'd like translated to English, please provide it, and I'll be happy to assist you with the translation.
I apologize, but there is no text provided in your message for me to translate. If you would like me to translate something, please provide the original text and I'll be happy to translate it to English for you.

### Options disponibles

-`-d, --debug`Enable debug mode
-`-i, --interactive`Activate interactive mode
There is no text provided to translate. The input is just a single hyphen (-). If you'd like me to translate something, please provide the text and specify the source language.`-b, --batch-size`Defines the batch size for large files (default: 1000)
-`-t, --threads`Defines the number of threads for parallel processing (default: 4)

### Test de l'API

I apologize, but there is no text provided in your message for me to translate. If you would like me to translate something, please provide the original text and specify which language it is in, and I'll be happy to translate it to English for you.
./translator test API
I apologize, but there is no text provided in your message for me to translate. If you'd like me to translate something for you, please provide the original text and specify which language it's in, and I'll be happy to translate it to English for you.

## Développement

### Structure du projet

I apologize, but there is no text provided in your message for me to translate. If you would like me to translate something, please provide the original text and specify which language it is in. Then I'll be happy to translate it to English for you.
translator
├── cmd/
│   └── translator/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── claude.go
│   ├── translation/
│   │   └── translator.go
│   └── cli/
└── interactive.go
└── pkg/
│   ├── fileutils/
│   │   └── fileutils.go
│   ├── markdown/
│   │   └── processor.go
│   └── tokenizer/
└── tokenizer.go
└── .gitignore
go.mod
└── README.md
I apologize, but there is no text provided in your message for me to translate. If you would like me to translate something, please provide the text you want translated and I'll be happy to help.

### Exécution des tests

I apologize, but there is no text provided in your message for me to translate. If you'd like me to translate something, please provide the original text and I'll be happy to translate it to English for you.
go test ./...
I apologize, but there is no text provided in your request to translate. Could you please provide the text you would like me to translate to English? Once you provide the text, I'll be happy to translate it for you without any additional comments or explanations.

## Contribution

Contributions are welcome! Please follow these steps:

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Licence

[License information coming soon]

## Contact

Your Name - your@email.com

Project link:[https://github.com/votreuser/translator](https://github.com/votreuser/translator)

