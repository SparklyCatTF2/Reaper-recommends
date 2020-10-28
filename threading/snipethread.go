package threading

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"strconv"

	"github.com/SparklyCatTF2/Reaper/globals"
	"github.com/SparklyCatTF2/Reaper/rblx"
)
var counter = 0
func SnipeThread(assetID int64, snipeChannel chan *rblx.PurchasePost) {
	emptystr := ""
	httpclient := &http.Client{}
	cachedPrices := make(map[int64]int64, 0)
	rand.Seed(time.Now().UnixNano())
	assetIDstr := strconv.Itoa(int(assetID))

	// Cache the token for the current and next price check cookie
	currentRobloxSession := &rblx.RBLXSession{Cookie: globals.PriceCheckCookies[rand.Intn(len(globals.PriceCheckCookies))], Client: httpclient, XCSRFToken: &emptystr}
	nextRobloxSession := &rblx.RBLXSession{Cookie: globals.PriceCheckCookies[rand.Intn(len(globals.PriceCheckCookies))], Client: httpclient, XCSRFToken: &emptystr}
	rblxsession := currentRobloxSession
	GrabToken(rblxsession, false)
	go GrabToken(nextRobloxSession, false)


	for {
		detailsResponse, err := rblxsession.GetRecommendations(8, globals.ContextAssetIDs[assetIDstr][0], 3)
		if err != nil {
			// Rate limit, change price check cookie
			switch err.Type {
			case rblx.TooManyRequests:
				currentRobloxSession = nextRobloxSession
				nextRobloxSession = &rblx.RBLXSession{Cookie: globals.PriceCheckCookies[rand.Intn(len(globals.PriceCheckCookies))], Client: httpclient, XCSRFToken: &emptystr}
				go GrabToken(nextRobloxSession, false)
			case rblx.AuthorizationDenied:
				fmt.Printf("[Reaper] Invalid price check cookie %s", rblxsession.Cookie)
			case rblx.TokenValidation:
				GrabToken(rblxsession, false)
			}
			continue
		}

		// Loop over the items & send the purchase details to main thread if snipe is profitable
		item := detailsResponse.Data[globals.ContextAssetIDs[assetIDstr][1]].Item
		if globals.Config.TrySnipe == true {
			counter++
			fmt.Printf("Counter: %d | LowestPrice: %d | Cached Price: %d \n", counter, item.Price, cachedPrices[assetID])
		}
		if item.Price == 0 {
			continue
		}
		if item.Price < cachedPrices[assetID] {
			if cachedPrices[assetID] == 0 {
				continue
			}
			getpercent := float64((30 * item.Price) / 100)
			oldPriceAfterTax := float64(cachedPrices[assetID])
			oldPriceAfterTax -= getpercent
			profitMargin := oldPriceAfterTax - float64(item.Price)
			profitPercent := profitMargin / float64(item.Price)
			if profitPercent >= globals.Config.ProfitPercent {
				purchaseStruct := &rblx.PurchasePost{AssetID: assetID, ExpectedCurrency: 1, ExpectedPrice: item.Price}
				snipeChannel <- purchaseStruct
			} else {
				cachedPrices[assetID] = item.Price
			}
		} else {
			cachedPrices[assetID] = item.Price
		}
	}
}
