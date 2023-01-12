package pankat

type config struct {
	DocumentsPath    string
	SiteURL          string
	SiteTitle        string
	MyMd5HashMapJson string
	Verbose          int
	Force            int
	ListenAndServe   string
}

var instance config

func GetConfig() *config {
	return &instance
}
