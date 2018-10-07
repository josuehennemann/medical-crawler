package main

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

const (
	MAGIC_STRING = ";" // sequencia de caracters que é utilizado para juntar strings
)

type crawlerItem interface {
	GetContent()
	Write(io.Writer)
	Request()
	getResponse() *http.Response
	getUrl() string
}

type StdCrawler struct {
	url          string
	items        map[string]struct{} //mapa para ir salvando os itens encontrados
	exists       bool                //flag que vai indicando se o registro existe no mapa ou nao
	httpResponse *http.Response
}

//inicializa o mapa
func (this *StdCrawler) makeItems() {
	this.items = map[string]struct{}{}
}

func (this *StdCrawler) getResponse() *http.Response {
	return this.httpResponse
}

func (this *StdCrawler) getUrl() string {
	return this.url
}

func (this *StdCrawler) add(name, email string) {
	if !strings.Contains(email, "@") {
		return
	}

	if _, this.exists = this.items[name+MAGIC_STRING+email]; !this.exists {
		this.items[name+MAGIC_STRING+email] = struct{}{}
	}
}

//faz uma requisiçao basica, utilizado apenas get
func (this *StdCrawler) BasicRequest() {
	resp, err := http.Get(this.url)
	checkError(err)
	this.httpResponse = resp
}

//faz uma requisiçao utilizando post, podendo especificar headers na requisiçao
func (this *StdCrawler) AdvancedRequest(headers map[string]string, body io.Reader) {

	req, er := http.NewRequest("POST", this.url, body)
	checkError(er)
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, er := http.DefaultClient.Do(req)
	this.httpResponse = resp
}

//implementa o GetContent generico, apenas para nao dar erro nas outras interfaces
func (this *StdCrawler) GetContent() {
}

//implementa o Request generico, apenas para nao dar erro nas outras interfaces
func (this *StdCrawler) Request() {
}

//implementa o write centralizado para todas as interfaces{}
func (this *StdCrawler) Write(output io.Writer) {

	lock.Lock()
	defer lock.Unlock()

	for item, _ := range this.items {

		if len(strings.Split(item, MAGIC_STRING)) != 2 {
			checkError(errors.New("Invalid line: " + item + " origem: " + this.url))
		}

		output.Write([]byte(item + "\n"))
	}
}
