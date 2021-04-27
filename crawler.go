package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dadosjusbr/coletores/status"
)

var urlFormats = map[string]string{
	"remu": "&__sessionId=%s&__format=xls&__asattachment=true&__overwrite=false",
}

// Inicializa um mapa com o formato da url complementar para cada tipo de planilha
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

func download(url string, filePath string) {

	resp, err := http.Get(url)
	if err != nil {
		status.ExitFromError(status.NewError(status.DataUnavailable, fmt.Errorf("Não foi possível fazer o download do arquivo: %s .O seguinte erro foi gerado: %q", filePath, err)))
		os.Exit(1)
	}

	file, err := os.Create(filePath)
	if err != nil {
		status.ExitFromError(status.NewError(status.DataUnavailable, fmt.Errorf("Não foi possível fazer o download do arquivo: %s .O seguinte erro foi gerado: %q", filePath, err)))
		os.Exit(1)
	}
	defer file.Close()

	io.Copy(file, resp.Body)
	defer resp.Body.Close()
}

func Crawl(month int, year int, outputPath string) []string {
	var paths []string
	complements := initComplements(month, year)

	for key, _ := range complements {
		var fileName = fmt.Sprint(year, "_", fmt.Sprintf("%02d", month), "_", key)
		var filePath = fmt.Sprint(outputPath, "/", fileName, ".xls")

		seasonId := seasonId(complements[key])
		url := fmt.Sprint(complements[key], fmt.Sprintf(urlFormats[key], seasonId))

		download(url, filePath)
		paths = append(paths, filePath)
	}

	return paths
}
