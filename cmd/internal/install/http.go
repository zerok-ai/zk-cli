package install

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

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
