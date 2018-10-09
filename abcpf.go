package main

import (

    //"net/http"
    "strings"
    "github.com/PuerkitoBio/goquery"
)

const (
    URL_ABCPF = "http://www.abcpf.org.br/medicos-associados/efetivo"
)

func NewCrawlerAbcpf() *CrawlerAbcpf {
    c := &CrawlerAbcpf{}
    c.url = URL_ABCPF
    c.makeItems()
    return c
}

type CrawlerAbcpf struct {
    StdCrawler
}

func (this *CrawlerAbcpf) Request() {
    this.BasicRequest()
}
func (this *CrawlerAbcpf) GetContent() {
    this.Request()
    defer this.httpResponse.Body.Close()

    this.getDocument()
    this.parse()
}

// analisa o html atual do site e separa os email que achou.
//Data da analise: 04/10/2018
func (this *CrawlerAbcpf) parse() {
    this.doc.Find("#content").Each(func(i int, s *goquery.Selection) {

        s.Find(".left > p").Each(func(a int, s1 *goquery.Selection) {
            tmp := strings.Split(strings.TrimSpace(s1.Eq(0).Text()), "\n")

            if len(tmp) <= 0 {
                return
            }

            name := strings.TrimSpace(tmp[0])
            email := ""
            for _, v := range tmp[1:] {
                if emailTmp := strings.Split(strings.TrimSpace(v), "E-mail: "); len(emailTmp) == 2 {
                    email = emailTmp[1]
                }  
            }

            if email == ""{
                return
            }

             this.add(name, email)
        })

    })
}

