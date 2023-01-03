package pankat

type config struct {
	InputPath        string
	OutputPath       string
	SiteURL          string
	SiteTitle        string
	MyMd5HashMapJson string
}

var instance config

func GetConfig() *config {
	return &instance
}
