package scraper

import {
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
}

type OfficialScraper struct {
	URL string
}

func (s *OfficialScraper) Fetch() ([]News, error){
	res, err := http.Get(s.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	var news []News

	doc.Find("article-post-item, .archive-noticias .item").Each(func(i int, sel *goquery.Selection){
		title := strings.TrimSpace(sel.Find("h2", "h3").Text())
		link, _ := sel.Find("a").Attr("href")

		if title != "" && link != "" {
			news = append(news, News{
				Title:	title,
				Link:	link,
				Source: "Site Oficial do Vasco",
			})
		}
	})

	return news, nil
}