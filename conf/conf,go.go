package conf

import "os"

type conf struct {
	OpenaiApiKey  string
	OpenaiBaseUrl string
}

var Conf conf

func init() {
	Conf.OpenaiApiKey = os.Getenv("OPENAI_API_KEY")
	Conf.OpenaiBaseUrl = os.Getenv("OPENAI_BASE_URL")
}
