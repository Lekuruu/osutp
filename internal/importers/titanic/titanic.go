package titanic

import "net/http"

type TitanicImporter struct {
	WebUrl string
	ApiUrl string
	client *http.Client
}

func NewTitanicImporter(webUrl, apiUrl string) *TitanicImporter {
	return &TitanicImporter{
		WebUrl: webUrl,
		ApiUrl: apiUrl,
		client: &http.Client{Timeout: defaultHttpTimeout},
	}
}
