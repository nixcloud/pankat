package pankat

type config struct {
	DocumentsPath    string
	SiteTitle        string
	MyMd5HashMapJson string
	Verbose          int
	Force            int
	ListenAndServe   string
}

var instance config

func Config() *config {
	return &instance
}
