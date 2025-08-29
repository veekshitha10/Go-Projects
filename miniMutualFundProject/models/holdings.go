package models
type HoldingsTable struct {
	SchemeCode string `json:"scheme_code"`
	UserId uint `json:"user_id"`
	Units float64 `json:"units"`
}
