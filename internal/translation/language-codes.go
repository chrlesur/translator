package translation

import "strings"

var LanguageToCode = map[string]string{
    "afrikaans": "af",
    "albanais": "sq",
    "allemand": "de",
    "amharique": "am",
    "anglais": "en",
    "arabe": "ar",
    "arménien": "hy",
    "assamais": "as",
    "aymara": "ay",
    "azerbaïdjanais": "az",
    "bambara": "bm",
    "basque": "eu",
    "bengali": "bn",
    "bhojpuri": "bho",
    "biélorusse": "be",
    "birman": "my",
    "bosnien": "bs",
    "bulgare": "bg",
    "catalan": "ca",
    "cebuano": "ceb",
    "chinois (simplifié)": "zh-CN",
    "chinois (traditionnel)": "zh-TW",
    "coréen": "ko",
    "corse": "co",
    "créole haïtien": "ht",
    "croate": "hr",
    "danois": "da",
    "dogri": "doi",
    "écossais": "gd",
    "espagnol": "es",
    "espéranto": "eo",
    "estonien": "et",
    "ewe": "ee",
    "finnois": "fi",
    "français": "fr",
    "frison": "fy",
    "galicien": "gl",
    "géorgien": "ka",
    "grec": "el",
    "guarani": "gn",
    "gujarati": "gu",
    "haoussa": "ha",
    "hawaïen": "haw",
    "hébreu": "he",
    "hindi": "hi",
    "hmong": "hmn",
    "hongrois": "hu",
    "igbo": "ig",
    "indonésien": "id",
    "irlandais": "ga",
    "islandais": "is",
    "italien": "it",
    "japonais": "ja",
    "javanais": "jv",
    "kannada": "kn",
    "kazakh": "kk",
    "khmer": "km",
    "kinyarwanda": "rw",
    "kirghiz": "ky",
    "konkani": "gom",
    "kurde": "ku",
    "laotien": "lo",
    "latin": "la",
    "letton": "lv",
    "lituanien": "lt",
    "luxembourgeois": "lb",
    "macédonien": "mk",
    "maithili": "mai",
    "malgache": "mg",
    "malais": "ms",
    "malayalam": "ml",
    "maltais": "mt",
    "maori": "mi",
    "marathi": "mr",
    "mizo": "lus",
    "mongol": "mn",
    "népalais": "ne",
    "norvégien": "no",
    "nyanja": "ny",
    "odia": "or",
    "oromo": "om",
    "ouzbek": "uz",
    "pachtô": "ps",
    "pendjabi": "pa",
    "persan": "fa",
    "polonais": "pl",
    "portugais": "pt",
    "quechua": "qu",
    "roumain": "ro",
    "russe": "ru",
    "samoan": "sm",
    "sanskrit": "sa",
    "serbe": "sr",
    "sesotho": "st",
    "shona": "sn",
    "sindhi": "sd",
    "singhalais": "si",
    "slovaque": "sk",
    "slovène": "sl",
    "somali": "so",
    "sotho du sud": "st",
    "soundanais": "su",
    "suédois": "sv",
    "swahili": "sw",
    "tadjik": "tg",
    "tamoul": "ta",
    "tatar": "tt",
    "tchèque": "cs",
    "telugu": "te",
    "thaï": "th",
    "tigrinya": "ti",
    "tsonga": "ts",
    "turc": "tr",
    "turkmène": "tk",
    "ukrainien": "uk",
    "urdu": "ur",
    "vietnamien": "vi",
    "xhosa": "xh",
    "yiddish": "yi",
    "yoruba": "yo",
    "zoulou": "zu",
}

var CodeToLanguage map[string]string

func init() {
    CodeToLanguage = make(map[string]string)
    for language, code := range LanguageToCode {
        CodeToLanguage[code] = language
    }
}

// GetCodeForLanguage retourne le code ISO 639-1 pour une langue donnée
func GetCodeForLanguage(language string) string {
    code, exists := LanguageToCode[strings.ToLower(language)]
    if exists {
        return code
    }
    return ""
}

// GetLanguageForCode retourne le nom de la langue pour un code ISO 639-1 donné
func GetLanguageForCode(code string) string {
    language, exists := CodeToLanguage[strings.ToLower(code)]
    if exists {
        return language
    }
    return ""
}
