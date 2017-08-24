package uaa

import (
	"io/ioutil"
	"net/http"
	"github.com/jhamon/uaa-cli/utils"
	"encoding/json"
)

type Keys struct {
	Keys []JWK
}

func TokenKeys(context UaaContext) ([]JWK, error) {
	tokenKeysUrl, err := utils.BuildUrl(context.BaseUrl, "token_keys")
	if err != nil {
		return []JWK{}, err
	}
	url := tokenKeysUrl.String()

	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept","application/json")

	resp, err := httpClient.Do(req)
	if (resp.StatusCode != 200 || err != nil) {
		key, err := TokenKey(context)
		return []JWK{key}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []JWK{}, unknownError()
	}

	keys := Keys{}
	err = json.Unmarshal(body,&keys)
	if err != nil {
		return []JWK{}, parseError(url, body)
	}

	return keys.Keys, nil
}