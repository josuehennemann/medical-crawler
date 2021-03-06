//Busca em diversos sites listas de medicos, clinicas medicas, ou cursos relacionados a area da saude e gera um csv

package main

import (
	"fmt"
	"os"
	"sync"
)

//TODO: implementar file ser uma flag
func main() {

	file, err := os.OpenFile("file.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	checkError(err)
	wg := sync.WaitGroup{}
	for _, crawler := range crawlerList {
		v := crawler
		wg.Add(1)
		go func(v crawlerItem) {
			defer wg.Done()
			//println(v.getUrl())
			v.GetContent()
			v.Write(file)
		}(v)
	}
	wg.Wait()
}

//lista generica de usuarios de email
var genericUserEmail = []string{"contato"}// "sac", "faleconosco"}

//lista de crawlers a serem executados
var crawlerList = []crawlerItem{
	NewCrawlerGuiaComercialBahia(),
	NewCrawlerAprofem(),
	NewCrawlerTelelista(),
	NewCrawlerAbcpf(),	
}

//lock para garantir a escrita no arquivo
var lock sync.Mutex

func checkError(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "error [%s]", e)
		os.Exit(1)
	}
}
