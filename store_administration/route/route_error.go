package route

func CreateErrorMessage(message string) map[string]string {
    return createJSONMessage("error", message)
}

func CreateMessage(message string) map[string]string {
    return createJSONMessage("message", message)
}

func createJSONMessage(key, value string) map[string]string {
    return map[string]string{key: value}
}
