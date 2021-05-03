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

type requestURLs struct {
	remunerationURL string
	benefitsURL     string
}

// Retorna as url para download de cada planilha em questão
func getRequestURLs(year, month int) (requestURLs, error) {
	remuIDURL := fmt.Sprint("https://servicos-portal.mpro.mp.br/plcVis/frameset?__report=..%2FROOT%2Frel%2Fcontracheque%2Fmembros%2FremuneracaoMembrosAtivos.rptdesign&anomes=", year, fmt.Sprintf("%02d", month), "&nome=&cargo=&lotacao=")
	remuSessionID, err := getSessionID(remuIDURL)
	if err != nil {
		return requestURLs{}, err
	}
	remuDownloadURL := fmt.Sprintf("%s&__sessionId=%s&__format=xls&__asattachment=true&__overwrite=false", remuIDURL, remuSessionID)

	viIDURL := fmt.Sprint("https://servicos-portal.mpro.mp.br/plcVis/frameset?__report=..%2FROOT%2Frel%2Fcontracheque%2Fmembros%2FverbasIndenizatoriasMembrosAtivos.rptdesign&anomes=", year, fmt.Sprintf("%02d", month))
	benefitsSessionID, err := getSessionID(viIDURL)
	if err != nil {
		return requestURLs{}, err
	}
	benefitsURL := fmt.Sprintf("%s&__sessionId=%s&__format=xls&__asattachment=true&__overwrite=false", viIDURL, benefitsSessionID)

	return requestURLs{remuDownloadURL, benefitsURL}, nil
}

// Inicializa o id de sessão para uma dada url
func getSessionID(url string) (string, error) {
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

	return id[1][0:19], err
}

func download(url string, filePath string, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return status.NewError(status.ConnectionError, fmt.Errorf("Problem doing GET on the URL(%s) to download the file(%s). Error: %q", url, filePath, err))
	}
	defer resp.Body.Close()

	_, err = os.Stat(outputPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(outputPath, 0755)
		if err != nil {
			return status.NewError(status.SystemError, fmt.Errorf("Error creating outputfolder (%s). Error: %q", outputPath, err))
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return status.NewError(status.DataUnavailable, fmt.Errorf("Error creating downloaded (%s) file(%s). Error: %q", url, filePath, err))
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return status.NewError(status.SystemError, fmt.Errorf("Was not possible to save the downloaded file: %s. The following mistake was teken: %q", filePath, err))
	}

	return nil
}

func Crawl(month int, year int, outputPath string) ([]string, error) {
	var paths []string

	request, err := getRequestURLs(year, month)
	if err != nil {
		return paths, err
	}

	for typ := 0; typ < 2; typ++ {
		switch typ {
		case remuURLType:
			var fileName = fmt.Sprintf("%d_%02d_remu.xls", year, month)
			var filePath = fmt.Sprint(outputPath, "/", fileName)

			err = download(request.remunerationURL, filePath, outputPath)
			if err != nil {
				return paths, err
			}

			paths = append(paths, filePath)
		case viURLType:
			var fileName = fmt.Sprintf("%d_%02d_vi.xls", year, month)
			var filePath = fmt.Sprint(outputPath, "/", fileName)

			err = download(request.benefitsURL, filePath, outputPath)
			if err != nil {
				return paths, err
			}

			paths = append(paths, filePath)
		}
	}
	return paths, nil
}
