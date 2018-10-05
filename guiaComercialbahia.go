package main

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

const (
	URL_GUIACOMERCIALBAHIA = "http://www.guiacomercialdabahia.com.br/Cl%C3%ADnicas+M%C3%A9dicas"
)

func NewCrawlerGuiaComercialBahia() *CrawlerGuiaComercialBahia {
	c := &CrawlerGuiaComercialBahia{}
	c.url = URL_GUIACOMERCIALBAHIA
	c.makeItems()
	return c
}

type CrawlerGuiaComercialBahia struct {
	StdCrawler
}

func (this *CrawlerGuiaComercialBahia) Request() {
	this.BasicRequest()
}
func (this *CrawlerGuiaComercialBahia) GetContent() {
	this.Request()
	defer this.httpResponse.Body.Close()
	doc, e := goquery.NewDocumentFromReader(this.httpResponse.Body)
	checkError(e)
	this.parse(doc)
}

// analisa o html atual do site e separa os email que achou.
//Data da analise: 04/10/2018
func (this *CrawlerGuiaComercialBahia) parse(doc *goquery.Document) {

	doc.Find(".col-xs-12").Each(func(i int, s *goquery.Selection) {

		s.Find(".row").Each(func(a int, s1 *goquery.Selection) {

			s.Find(".col-xs-12.col-sm-6.col-md-7.col-lg-8").Each(func(a int, s1 *goquery.Selection) {

				email := s1.Find("div").Eq(3).Text()
				name := s1.Find("div").Eq(0).Text()

				name = strings.TrimSpace(name)
				email = strings.Replace(email, "&nbsp;", "", -1)

				this.add(name, email)
				/*if !strings.Contains(email, "@") {
					return
				}
				if _, this.exists = this.items[name+MAGIC_STRING+email]; !this.exists {
					this.items[name+MAGIC_STRING+email] = struct{}{}
				}*/
			})

		})

	})
}
