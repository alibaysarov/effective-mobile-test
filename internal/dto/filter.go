package dto

type Filter struct {
	StartDateFrom *string `json:"start_date_from,omitempty" form:"start_date_from"`
	StartDateTo   *string `json:"start_date_to,omitempty" form:"start_date_to"`
	UserId        *string `json:"userId,omitempty" form:"userId"`
	ServiceName   *string `json:"service_name,omitempty" form:"service_name"`
}
