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

type urlRequest struct {
	remuDownloadURL string
}

// Retorna as url para download de cada planilha em questão
func initRequests(month, year int) (urlRequest, error) {

	idURL := fmt.Sprint("https://servicos-portal.mpro.mp.br/plcVis/frameset?__report=..%2FROOT%2Frel%2Fcontracheque%2Fmembros%2FremuneracaoMembrosAtivos.rptdesign&anomes=", year, fmt.Sprintf("%02d", month), "&nome=&cargo=&lotacao=")
	sessionId, err := seasonId(idURL)
	if err != nil {
		return urlRequest{}, err
	}

	downloadURL := fmt.Sprint(idURL, fmt.Sprintf("&__sessionId=%s&__format=xls&__asattachment=true&__overwrite=false", sessionId))
	return urlRequest{downloadURL}, nil

}

// Inicializa o id de sessão para uma dada url
func seasonId(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", status.NewError(status.ConnectionError, fmt.Errorf("Was not possible to get a season id to the url: %s. %q", url, err))
	}
	defer resp.Body.Close()

	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", status.NewError(status.ConnectionError, fmt.Errorf("Was not possible to get a season id to the url: %s. %q", url, err))
	}

	htmlCode := string(page)
	id := strings.Split(htmlCode, "Constants.viewingSessionId = \"")
	seasonId := id[1][0:19]

	return seasonId, err
}

func download(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return status.NewError(status.DataUnavailable, fmt.Errorf("Was not possible download the file: %s .The following mistake was taken: %q", filePath, err))
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return status.NewError(status.DataUnavailable, fmt.Errorf("Was not possible download the file: %s .The following mistake was taken: %q", filePath, err))
	}
	defer file.Close()

	io.Copy(file, resp.Body)
	return nil
}

func Crawl(month int, year int, outputPath string) ([]string, error) {
	var paths []string

	var fileName = fmt.Sprint(year, "_", fmt.Sprintf("%02d", month), "_remu")
	var filePath = fmt.Sprint(outputPath, "/", fileName, ".xls")

	request, err := initRequests(year, month)
	if err != nil {
		return paths, err
	}

	download(request.remuDownloadURL, filePath)
	paths = append(paths, filePath)

	return paths, nil
}
