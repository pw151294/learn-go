package main

import "os"

func getModelName() string {
	return os.Getenv("DEEPSEEK_MODEL_NAME")
}

func getAPIKey() string {
	return os.Getenv("DEEPSEEK_API_KEY")
}

func getBaseURL() string {
	return os.Getenv("DEEPSEEK_BASE_URL")
}
