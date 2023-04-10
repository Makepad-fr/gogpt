package gogpt

import "encoding/json"

type AccountPlan struct {
	IsPaidSubscriptionActive       bool   `json:"is_paid_subscription_active"`
	SubscriptionPlan               string `json:"subscription_plan"`
	AccountUserRole                string `json:"account_user_role"`
	WasPaidCustomer                bool   `json:"was_paid_customer"`
	HasCustomerObject              bool   `json:"has_customer_object"`
	SubscriptionExpiresAtTimestamp int64  `json:"subscription_expires_at_timestamp"`
}

type UserAccountInfo struct {
	AccountPlan AccountPlan `json:"account_plan"`
	UserCountry string      `json:"user_country"`
	Features    []string    `json:"features"`
}

func UnmarshalUserAccountInfo(jsonData []byte) (UserAccountInfo, error) {
	var userAccountInfo UserAccountInfo
	err := json.Unmarshal(jsonData, &userAccountInfo)
	return userAccountInfo, err
}
