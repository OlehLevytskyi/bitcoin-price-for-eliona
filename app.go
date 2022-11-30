//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"fmt"
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"hailo/apiserver"
	"hailo/apiservices"
	"hailo/conf"
	"hailo/eliona"
	"hailo/rate"
	"net/http"
)

// CollectData reads the defined currencies rate from configuration and writes the data as eliona heap
func CollectData() {

	currencies := make(chan conf.Currency)
	go func() {
		_ = conf.ReadCurrencies(currencies)
	}()
	for currency := range currencies {

		// Check if eliona project is defined; if not, ignore this currency
		if len(currency.ProjectId) == 0 {
			log.Warn("Bitcoin rate", "Ignoring currency, because no project id is defined: %s", currency.Code)
			continue
		}

		// Create or update asset and get ID
		assetId, err := asset.UpsertAsset(api.Asset{
			ProjectId:             currency.ProjectId,
			GlobalAssetIdentifier: fmt.Sprintf("%s %s", currency.Code, currency.Description),
			Description:           *api.NewNullableString(&currency.Description),
			Name:                  *api.NewNullableString(&currency.Code),
			AssetType:             "bitcoin_rate",
		})
		log.Debug("Bitcoin rate", "Determining asset id %d for currency '%s'", *assetId, currency.Code)

		// Reads the current rate for bitcoin
		coinRate, err := rate.Today(currency)
		if err != nil {
			log.Error("Bitcoin rate", "Error during requesting API endpoint: %v", err)
			return
		}
		log.Debug("Bitcoin rate", "New rate for currency '%s' found: %f", coinRate.Code, coinRate.Rate)

		// Writes input data as heap
		eliona.UpsertHeap(api.SUBTYPE_INPUT, *assetId, eliona.Input{
			Code: coinRate.Code,
			Rate: coinRate.Rate,
		})

		// Writes info data as heap
		eliona.UpsertHeap(api.SUBTYPE_INFO, *assetId, eliona.Info{
			Daytime: coinRate.Daytime,
		})

		// Writes status data as heap
		eliona.UpsertHeap(api.SUBTYPE_STATUS, *assetId, eliona.Status{
			Comment: coinRate.Comment,
		})
	}
}

// listenApiRequests starts an API server and listen for API requests
// The API endpoints are defined in the openapi.yaml file
func listenApiRequests() {
	err := http.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", "3000"), apiserver.NewRouter(
		apiserver.NewConfigurationApiController(apiservices.NewConfigurationApiService()),
	))
	log.Fatal("Hailo", "Error in API Server: %v", err)
}
