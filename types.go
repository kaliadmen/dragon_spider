package dragonSpider

type initPaths struct {
	//root of application
	rootPath string
	//directories available to application
	dirNames []string
}

//cookieConfig holds data for cookie configuration
type cookieConfig struct {
	name       string
	lifetime   string
	persistent string
	Secure     string
	domain     string
}
