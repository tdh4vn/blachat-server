package models


type Pageable struct {
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	Data interface{} `json:"data"`
}