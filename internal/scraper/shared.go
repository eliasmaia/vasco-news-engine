package scraper

// removeDuplicateNews remove links duplicados de uma lista de notícias
func removeDuplicateNews(newsList []News) []News {
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