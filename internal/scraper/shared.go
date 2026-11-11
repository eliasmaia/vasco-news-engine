package scraper

type News struct {
	Title	string
	Link	string
	Source	string
}

//defines the standard contract for all scrapers (ISP)
type SiteScraper interface {
		Fetch() ([]News, error)
}

// removeDuplicateNews remove duplicate links from a list of news
func RemoveDuplicateNews(newsList []News) []News {
	keys := make(map[string]bool)
	var list []News
	for _, item := range newsList {
		if _, value := keys[item.Link]; !value {
			keys[item.Link] = true
			list = append(list, item)
		}
	}
	return list
}