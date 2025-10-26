package dto

type SubscriptionDto struct {
	ID          string  `json:id,omitempty`
	ServiceName string  `json:serviceName:required`
	UserId      string  `json:userId:required`
	Price       float64 `json:price:required`
	StartDate   string  `json:startDate:required`
}
