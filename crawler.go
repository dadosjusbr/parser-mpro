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

const (
	viURLType   int = 0
	remuURLType int = 1
)

type urlRequests struct {
	remuDownloadURL string
	viDownloadURL   string
}

// Retorna as url para download de cada planilha em questão
func requestURL(year, month int) (urlRequests, error) {
	remuIDURL := fmt.Sprint("https://servicos-portal.mpro.mp.br/plcVis/frameset?__report=..%2FROOT%2Frel%2Fcontracheque%2Fmembros%2FremuneracaoMembrosAtivos.rptdesign&anomes=", year, fmt.Sprintf("%02d", month), "&nome=&cargo=&lotacao=")
	remuSessionID, err := seasonID(remuIDURL)
	if err != nil {
		return urlRequests{}, err
	}
	remuDownloadURL := fmt.Sprint(remuIDURL, fmt.Sprintf("&__sessionId=%s&__format=xls&__asattachment=true&__overwrite=false", remuSessionID))

	viIDURL := fmt.Sprint("https://servicos-portal.mpro.mp.br/plcVis/frameset?__report=..%2FROOT%2Frel%2Fcontracheque%2Fmembros%2FverbasIndenizatoriasMembrosAtivos.rptdesign&anomes=", year, fmt.Sprintf("%02d", month))
	viSessionID, err := seasonID(viIDURL)
	if err != nil {
		return urlRequests{}, err
	}
	viDownloadURL := fmt.Sprint(viIDURL, fmt.Sprint("&__sessionId=%s&__format=xls&__asattachment=true&__overwrite=false", viSessionID))

	return urlRequests{remuDownloadURL, viDownloadURL}, nil
}

// Inicializa o id de sessão para uma dada url
func seasonID(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", status.NewError(status.ConnectionError, fmt.Errorf("Was not possible to get a season id to the url: %s. %q", url, err))
	}
	defer resp.Body.Close()

	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", status.NewError(status.ConnectionError, fmt.Errorf("Was not possible to get a season id to the url: %s. %q", url, err))
	}

	id := strings.Split(string(page), "Constants.viewingSessionId = \"")
	seasonId := id[1][0:19]

	return seasonId, err
}

func download(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return status.NewError(status.ConnectionError, fmt.Errorf("Problem doing GET on the URL(%s) to download the file(%s). Error: %q", url, filePath, err))
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return status.NewError(status.DataUnavailable, fmt.Errorf("Error creating downloaded (%s) file(%s). Error: %q", url, filePath, err))
	}
	defer file.Close()

	_, erro := io.Copy(file, resp.Body)
	if erro != nil {
		return status.NewError(status.SystemError, fmt.Errorf("Was not possible to save the downloaded file: %s. The following mistake was teken: %q", filePath, erro))
	}
	return nil
}

func Crawl(month int, year int, outputPath string) ([]string, error) {
	var paths []string

	request, err := requestURL(year, month)
	if err != nil {
		return paths, err
	}

	for typ := 0; typ < 2; typ++ {
		switch typ {
		case remuURLType:
			var fileName = fmt.Sprint("%d", "_", "%02d", "_remu", year, month)
			var filePath = fmt.Sprint(fileName, ".xls")

			err = download(request.remuDownloadURL, filePath)
			if err != nil {
				return paths, err
			}

			paths = append(paths, filePath)
		case viURLType:
			var fileName = fmt.Sprintf("%d", "_", "%02d", "_vi", year, month)
			var filePath = fmt.Sprintf(fileName, ".xls")

			err = download(request.viDownloadURL, filePath)
			if err != nil {
				return paths, err
			}

			paths = append(paths, filePath)
		}
	}
	return paths, nil
}
