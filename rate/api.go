package rate

import (
	"encoding/json"
	"fmt"
	"github.com/eliona-smart-building-assistant/go-utils/http"
	"hailo/conf"
	"time"
)

type CoinRate struct {
	Code    string  `json:"code"`
	Rate    float64 `json:"rate"`
	Comment string  `json:"comment"`
	Daytime string  `json:"daytime"`
}

func Today(currency conf.Currency) (CoinRate, error) {
	var coinRate CoinRate

	// Request the API to get current conditions conditions
	url, payload, err := request(currency)
	if err != nil {
		return coinRate, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(payload, &result)
	if err != nil {
		return coinRate, err
	}

	if result["chartName"].(string) != "Bitcoin" {
		return coinRate, fmt.Errorf("error requesting api %s: %s", url, result["disclaimer"].(string))
	}
	coinRate.Code, _ = result["bpi"].(map[string]interface{})[currency.Code].(map[string]interface{})["code"].(string)
	coinRate.Rate, _ = result["bpi"].(map[string]interface{})[currency.Code].(map[string]interface{})["rate_float"].(float64)
	coinRate.Comment, _ = result["bpi"].(map[string]interface{})[currency.Code].(map[string]interface{})["description"].(string)
	coinRate.Daytime, _ = result["time"].(map[string]interface{})["updated"].(string)

	return coinRate, nil
}

// request calls the api to get structured currency data
func request(currency conf.Currency) (string, []byte, error) {
	url := fmt.Sprintf(conf.Endpoint())
	request, err := http.NewRequest(url)
	if err != nil {
		return url, nil, err
	}
	payload, err := http.Do(request, time.Second*10, true)
	return url, payload, err
}
