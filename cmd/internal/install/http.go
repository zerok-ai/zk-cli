package install

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var myClient = &http.Client{Timeout: 20 * time.Second}

func GetHTTPGETResponse(url string, target interface{}) error {
	response, err := myClient.Get(url)
	var resErr error = nil
	if response != nil {
		if response.StatusCode < 200 || response.StatusCode >= 299 {
			resErr = errors.New(fmt.Sprintf("Error response code = %d", response.StatusCode))
		}
	}
	if err != nil {
		if resErr != nil {
			return fmt.Errorf(resErr.Error(), err)
		}
		return err
	} else if resErr != nil {
		return resErr
	}

	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}

func GetHTTPPOSTResponse(url string, body interface{}, target interface{}) error {

	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error:", err)
		return errors.New("Cluster metadata is not valid")
	}

	log.Println("url=%s", url)
	// Create an io.Reader from the JSON data
	bodyReader := strings.NewReader(string(jsonData))

	response, err := myClient.Post(url, "application/json", bodyReader)
	var resErr error = nil
	if response != nil {
		if response.StatusCode < 200 || response.StatusCode >= 299 {
			resErr = errors.New(fmt.Sprintf("Error response code = %d", response.StatusCode))
		}
	}
	if err != nil {
		if resErr != nil {
			return fmt.Errorf(resErr.Error(), err)
		}
		return err
	} else if resErr != nil {
		return resErr
	}

	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}
