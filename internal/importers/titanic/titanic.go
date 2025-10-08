package titanic

type TitanicImporter struct {
	WebUrl string
	ApiUrl string
}

func NewTitanicImporter(webUrl, apiUrl string) *TitanicImporter {
	return &TitanicImporter{
		WebUrl: webUrl,
		ApiUrl: apiUrl,
	}
}
