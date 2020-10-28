package globals

import "time"

type config struct {
	Cookie           string  `json:"cookie"`
	IDsPerThread     int64   `json:"idsPerThread"`
	ThreadMultiplier int64   `json:"threadMultiplier"`
	Webhook          string  `json:"webhook"`
	WebhookMention   string  `json:"webhookMention"`
	ProfitPercent    float64 `json:"profitPercent"`
	AutoSell         bool    `json:"autoSell"`
	ProfileAvatar 	 string  `json:"ProfileAvatar"`
	Alias 			 string  `json:"YourAlias"`
	TrySnipe 		 bool  	 `json:"TrySnipePrint"`
}

var (
	Config            config
	PositiveQuotes []string
	NegativeQuotes []string
	ContextAssetIDs = make(map[string][]int64)
	CachedTokens      = make(map[string]string, 0)
	CachedProductIDs  = make(map[int64]int64, 0)
	CachedAssetNames = make(map[int64]string, 0)
	BlockedAssetIds = make(map[int64]int64, 0)
	PriceCheckCookies []string
)

// GetTime gets the current time from epoch in milliseconds
func GetTimeInMs() int64 {
	return int64(time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
}
