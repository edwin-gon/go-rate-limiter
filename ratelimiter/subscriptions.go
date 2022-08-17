package ratelimiter

type Subscription interface {
	Name() string
	RequestLimit() int
	TimeFrame() int64 // Millisecond count
}

type BasicSubscription struct {
	name         string
	requestLimit int
	timeFrame    int64
}

type PremiumSubscription struct {
	name         string
	requestLimit int
	timeFrame    int64
}

const (
	basicName         = "Basic"
	basicRequestLimit = 5
	basicTimeFrame    = 60000

	premiumName         = "Premium"
	premiumRequestLimit = 20
	premiumTimeFrame    = 60000
)

func (sub BasicSubscription) Name() string {
	return basicName
}

func (sub BasicSubscription) RequestLimit() int {
	return basicRequestLimit
}

func (sub BasicSubscription) TimeFrame() int64 {
	return basicTimeFrame
}

func NewBasicSubscription() BasicSubscription {
	return BasicSubscription{basicName, basicRequestLimit, basicTimeFrame}
}

func (sub PremiumSubscription) Name() string {
	return premiumName
}

func (sub PremiumSubscription) RequestLimit() int {
	return premiumRequestLimit
}

func (sub PremiumSubscription) TimeFrame() int64 {
	return premiumTimeFrame
}

func NewPremiumSubscription() PremiumSubscription {
	return PremiumSubscription{premiumName, premiumRequestLimit, premiumTimeFrame}
}
