package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dadosjusbr/coletores/status"
)

// Inicializa um mapa com o formato da url para cada tipo de planilha
func initComplements(month, year int) map[string]string {
	return map[string]string{
		"remu": fmt.Sprint("https://servicos-portal.mpro.mp.br/plcVis/frameset?__report=..%2FROOT%2Frel%2Fcontracheque%2Fmembros%2FremuneracaoMembrosAtivos.rptdesign&anomes=", year, fmt.Sprintf("%02d", month), "&nome=&cargo=&lotacao="),
	}
}

// Inicializa o id de sessão para uma dada url
func seasonId(url string) string {

	resp, err := http.Get(url)
	if err != nil {
		status.ExitFromError(status.NewError(status.ConnectionError, fmt.Errorf("Was not possible to get a season id to the url: %s. %q", url, err)))
		os.Exit(1)
	}
	defer resp.Body.Close()

	page, err := ioutil.ReadAll(resp.Body)
	htmlCode := string(page)

	id := strings.Split(htmlCode, "Constants.viewingSessionId = \"")
	seasonId := id[1][0:19]

	return seasonId
}

func Crawl(month int, year int, outputPath string) []string {
	var paths []string
	complements := initComplements(month, year)

	for key, complement := range complements {
		var fileName = "file" + key + ".xls"
		seasonId := seasonId(complement)
		fmt.Println(fileName, seasonId)
	}

	return paths
}
