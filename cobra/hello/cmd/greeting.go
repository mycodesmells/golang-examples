package cmd

func greeting(lang string) string {
	switch lang {
	case "pl":
		return "Cześć"
	case "es":
		return "Hola"
	default: // "en"
		return "Hello"
	}
}
