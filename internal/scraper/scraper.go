package scraper

type News struct {
	Title  string
	Link   string
	Source string
}

type SiteScraper interface {
	Fetch() ([]News, error)
}
