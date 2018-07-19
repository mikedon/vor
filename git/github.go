package git

import (
	"encoding/json"
	"github.com/trevor-atlas/vor/logger"
	"strings"
	"github.com/trevor-atlas/vor/utils"
	"io/ioutil"
	"fmt"
	"bytes"
	"net/http"
)

// func GeneratePRName(branchName string) string {
	// branchName
// }

func Post (url string, requestBody []byte) PullRequestResponse {
	githubAPIKey := utils.GetStringEnv("github.apikey")
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
    req.Header.Set("Authorization", "token " + githubAPIKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	parsed := PullRequestResponse{}
	logger.Debug("response Status:", resp.Status)
	logger.Debug("response Headers:", resp.Header)
	logger.Debug("response Body:", string(body))

	parseError := json.Unmarshal(body, &parsed)
	if parseError != nil {
		fmt.Printf("error parsing json\n %s", parseError)
		panic(parseError)
	}
	return parsed

}