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

package eliona

import (
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"time"
)

// Input is a structure holds input data getting from the api endpoint. This structure corresponds
// to the input heap data in eliona.
type Input struct {
	Code string  `json:"code"`
	Rate float64 `json:"rate"`
}

// Info is a structure holds informational data getting from the api endpoint. This structure corresponds
// to the info heap data in eliona.
type Info struct {
	Daytime string `json:"daytime"`
}

// Status is a structure holds data getting from the api endpoint related to bitcoin rate.
// This structure corresponds to the status heap data in eliona.
type Status struct {
	Comment string `json:"comment"`
}

// InitAssetType creates asset type for bitcoin rate
func InitAssetType(db.Connection) error {
	err := asset.UpsertAssetType(api.AssetType{
		Name:        "bitcoin_rate",
		Custom:      common.Ptr(true),
		Vendor:      *api.NewNullableString(common.Ptr("coindesk")),
		Translation: *api.NewNullableTranslation(&api.Translation{De: common.Ptr("Bitcoin kurs"), En: common.Ptr("Bitcoin rate")}),
		Urldoc:      *api.NewNullableString(common.Ptr("https://api.coindesk.com/v1/bpi/currentprice.json")),
		Icon:        *api.NewNullableString(common.Ptr("bitcoin")),
		Attributes: []api.AssetTypeAttribute{
			{
				Type:        *api.NewNullableString(common.Ptr("bitcoin")),
				Name:        "code",
				Subtype:     api.SUBTYPE_INPUT,
				Translation: *api.NewNullableTranslation(&api.Translation{De: common.Ptr("Code"), En: common.Ptr("Code")}),
				Enable:      common.Ptr(true),
			},
			{
				Type:        *api.NewNullableString(common.Ptr("bitcoin")),
				Name:        "rate",
				Subtype:     api.SUBTYPE_INPUT,
				Translation: *api.NewNullableTranslation(&api.Translation{De: common.Ptr("Kurs"), En: common.Ptr("Rate")}),
				Enable:      common.Ptr(true),
			},
			{
				Type:        *api.NewNullableString(common.Ptr("bitcoin")),
				Name:        "daytime",
				Subtype:     api.SUBTYPE_INFO,
				Translation: *api.NewNullableTranslation(&api.Translation{De: common.Ptr("Tageszeit"), En: common.Ptr("Daytime")}),
				Enable:      common.Ptr(true),
			},
			{
				Type:        *api.NewNullableString(common.Ptr("bitcoin")),
				Name:        "comment",
				Subtype:     api.SUBTYPE_STATUS,
				Translation: *api.NewNullableTranslation(&api.Translation{De: common.Ptr("Kommentar"), En: common.Ptr("Comment")}),
				Enable:      common.Ptr(true),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func UpsertHeap(subtype api.DataSubtype, assetId int32, data any) {
	var statusHeap api.Data
	statusHeap.Subtype = subtype
	statusHeap.Timestamp = *api.NewNullableTime(common.Ptr(time.Now()))
	statusHeap.AssetId = assetId
	statusHeap.Data = common.StructToMap(data)
	err := asset.UpsertData(statusHeap)
	if err != nil {
		log.Error("Bitcoin rate", "Error during writing heap: %v", err)
		return
	}
}
