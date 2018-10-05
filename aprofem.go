package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
)

const (
	URL_APROFEM = "http://portal.aprofem.com.br/Dados/hlDados.ashx"
)

func NewCrawlerAprofem() *CrawlerAprofem {
	c := &CrawlerAprofem{}
	c.url = URL_APROFEM
	c.makeItems()
	c.headers = map[string]string{
		"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:62.0) Gecko/20100101 Firefox/62.0",
		"Accept":           "text/plain, */*; q=0.01",
		"Accept-Language":  "pt-BR,pt;q=0.8,en-US;q=0.5,en;q=0.3",
		"Referer":          "http://portal.aprofem.com.br/especialidades-medicas",
		"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
		"X-Requested-With": "XMLHttpRequest",
		"Cookie":           "ASP.NET_SessionId=34gb0eloqtk2ogtiuxyubuvt; AWSELB=A771C1790C513CEC79357B2A4AAFF27B4E172A7547B524DD915402FA7A667CB88B976DE0805D07A45CE7A2CDF8CCD3C1A2E2538CD9347649EF0E512B977E646775AD8E3A3A; _ga=GA1.3.955225383.1538698399; _gid=GA1.3.179746040.1538698399; _gat=1",
		"Connection":       "keep-alive",
		"Pragma":           "no-cache",
		"Cache-Control":    "no-cache",
	}
	return c
}

type CrawlerAprofem struct {
	StdCrawler
	headers map[string]string
}

func (this *CrawlerAprofem) Request() {
	body := strings.NewReader("validarSessao=0&tipoParceiro=004&pagina=1&qtdePorPagina=200&idFiltro_1=&idOpcaoN1_1=&idOpcaoN2_1=&idFiltro_2=&idOpcaoN1_2=&idOpcaoN2_2=&idFiltro_3=&idOpcaoN1_3=&idOpcaoN2_3=&m=listarParceiros")
	this.AdvancedRequest(this.headers, body)
}

func (this *CrawlerAprofem) GetContent() {
	this.Request()
	defer this.httpResponse.Body.Close()
	content, err := ioutil.ReadAll(this.httpResponse.Body)
	checkError(err)
	this.parse(content)
}

// analisa o json de retorno do site e separa os email que achou.
//Data da analise: 04/10/2018
func (this *CrawlerAprofem) parse(body []byte) {
	body = body[3:] // remove 1[ , pois o json que retorna nessa url começa com isso
	magicByte := []byte("}]")
	tmp := bytes.Split(body, magicByte) //separa o lixo que vem depois do fechamentod o json
	body = append(tmp[0], magicByte...) //recoloca o fechamento do json

	list := []map[string]interface{}{}
	checkError(json.Unmarshal(body, &list))

	for _, v := range list {
		name, _ := v["nomeParceiro"].(string)
		email, _ := v["email"].(string)
		this.add(name, email)
	}

}
