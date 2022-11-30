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

package conf

import (
	"github.com/eliona-smart-building-assistant/go-utils/db"
)

// InitConfiguration initialize the configuration of the app
func InitConfiguration(connection db.Connection) error {
	err := Set(connection, "endpoint", "https://api.coindesk.com/v1/bpi/currentprice.json")
	if err != nil {
		return err
	}
	err = Set(connection, "polling_interval", "10")
	if err != nil {
		return err
	}
	return nil
}

// InitCurrencies adds example assets for currencies USD, GBP and EUR.
// This should be made editable by eliona frontend.
func InitCurrencies(connection db.Connection) error {

	err := InsertCurrency(connection, Currency{
		Code:        "USD",
		Description: "United States Dollar",
	})
	if err != nil {
		return err
	}

	err = InsertCurrency(connection, Currency{
		Code:        "GBP",
		Description: "British Pound Sterling",
	})
	if err != nil {
		return err
	}

	err = InsertCurrency(connection, Currency{
		Code:        "EUR",
		Description: "Euro",
	})
	if err != nil {
		return err
	}

	return nil
}
