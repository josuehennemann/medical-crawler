package main

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	URL_TELELISTA     = "https://www.telelistas.net/{{estado}}/{{cidade}}/clinicas+medicas?pagina={{pagina}}&randsort=1851209712"
	MAX_PAGETELELISTA = 100
)

//lista de estado/cidade que vai fazer o replace na url da telelista
var (
	replaceTelelista = map[string][]string{"rs": []string{"porto+alegre", "santa+maria", "pelotas", "caxias+do+sul"},
		"sp": []string{"sao+paulo", "santos", "campinas", "ribeirao+preto"},
		"rj": []string{"rio+de+janeiro", "niteroi", "teresopolis", "macae"},
		"df": []string{"brasilia"},
		"pr": []string{"curitiba", "londrina", "cascavel", "maringa"},
		"sc": []string{"florianopolis", "blumenau", "chapeco", "itajai"},
		"mg": []string{"belo+horizonte", "uberlandia", "juiz+de+fora", "contagem"},
		"bh": []string{"salvador", "feira+de+santana", "vitoria+da+conquista", "juazeiro"},
	}
	regexpGetDomainTelelista = regexp.MustCompile(`window.location.href = '(http|https)://(www.|)(.*)';`)
)

func NewCrawlerTelelista() *CrawlerTelelistalista {
	c := &CrawlerTelelistalista{}
	c.url = URL_TELELISTA
	c.makeItems()
	c.page = 1
	return c
}

type CrawlerTelelistalista struct {
	StdCrawler
	page uint16
	doc  *goquery.Document
}

func (this *CrawlerTelelistalista) Request() {
	this.BasicRequest()
}
func (this *CrawlerTelelistalista) GetContent() {
	for uf, cities := range replaceTelelista {
		this.page = 1
		for _, city := range cities {
			urlbase := URL_TELELISTA
			urlbase = strings.Replace(urlbase, "{{estado}}", uf, -1)
			urlbase = strings.Replace(urlbase, "{{cidade}}", city, -1)
			for {
				this.url = strings.Replace(urlbase, "{{pagina}}", strconv.FormatUint(uint64(this.page), 10), -1)
				if this.makeRequestAndParse() {
					break
				}
			}
		}
	}
}

func (this *CrawlerTelelistalista) makeRequestAndParse() bool {
	//	println(this.url)
	this.Request()
	defer this.httpResponse.Body.Close()
	this.getDocument()
	if this.finish() {
		return true
	}
	this.parse()
	this.page++

	if this.page >= MAX_PAGETELELISTA {
		return true
	}

	return false
}

// analisa o html atual do site e separa os email que achou.
//Data da analise: 04/10/2018
func (this *CrawlerTelelistalista) parse() {

	content := this.doc.Find("#Content_Regs")
	if content == nil {
		return
	}

	content.Find("table").Each(func(a int, s *goquery.Selection) {

		s.Find("table").Each(func(a int, s1 *goquery.Selection) {
			name := s1.Find(".nome_resultado_ag").Text()

			if name == "" {
				name = s1.Find("tbody > tr > td").Eq(1).Text()
			}
			domain := ""
			s1.Find(".ib_ser > a").Each(func(a int, s2 *goquery.Selection) {
				value, exists := s2.Attr("title")

				if exists && strings.Contains(strings.TrimSpace(value), "Acesse o site e tenha") {
					link := ""
					if link, exists = s2.Attr("href"); !exists {
						return
					}

					domain = this.getDomain(link)
					return
				}

			})

			if domain == "" {
				return
			}

			name = strings.TrimSpace(name)
			for _, userEmail := range genericUserEmail {
				this.add(name, userEmail+"@"+domain)
			}
		})

	})
}
func (this *CrawlerTelelistalista) getDomain(url string) (domain string) {
	if url == "" {
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	destination := regexpGetDomainTelelista.FindStringSubmatch(string(content))
	if len(destination) != 4 {
		return
	}
	domain = destination[3]
	if idx := strings.Index(destination[3], "/"); idx > -1 {
		domain = destination[3][:idx]
	}
	return
}
func (this *CrawlerTelelistalista) getDocument() {
	var e error
	this.doc, e = goquery.NewDocumentFromReader(this.httpResponse.Body)
	checkError(e)
}

//valida se chegou na ultima pagina
func (this *CrawlerTelelistalista) finish() bool {
	selection := this.doc.Find(".tit_erro")
	if selection == nil {
		return false
	}
	return strings.Contains(selection.Text(), "Não foram encontrados")
}
