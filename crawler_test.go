package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestCrawlerGuiaComercialBahia(t *testing.T) {
	crawler := NewCrawlerGuiaComercialBahia()
	genericTest(t, crawler, true)
}

func TestCrawlerAprofem(t *testing.T) {
	crawler := NewCrawlerAprofem()
	genericTest(t, crawler, true)
}

func TestCrawlerGuiaComercialBahiaWrite(t *testing.T) {
	crawler := NewCrawlerGuiaComercialBahia()
	crawler.GetContent()
	crawler.Write(os.Stdout)
}

func TestCrawlerAprofemWrite(t *testing.T) {
	crawler := NewCrawlerAprofem()
	crawler.GetContent()
	crawler.Write(os.Stdout)
}

func TestCrawlerTelelistaWrite(t *testing.T) {
	crawler := NewCrawlerTelelista()
	crawler.GetContent()
	crawler.Write(os.Stdout)
}

func TestCrawlerCrawlerAbcpfWrite(t *testing.T) {
	crawler := NewCrawlerAbcpf()
	crawler.GetContent()
	crawler.Write(os.Stdout)
}

func genericTest(t *testing.T, crawler crawlerItem, getContent bool) {
	crawler.Request()
	httpResponse := crawler.getResponse()
	if httpResponse == nil {
		t.Fatal("Request failed")
	}
	if httpResponse.StatusCode != http.StatusOK {
		t.Fatal("Request failed http code ", httpResponse.StatusCode)
	}
	if getContent {
		crawler.GetContent()
		return
	}
	content, err := ioutil.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s", content)
	fmt.Println("======================")
}
