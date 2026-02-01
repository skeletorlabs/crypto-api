package models

type StandardResponse[T any] struct {
	Meta Meta `json:"meta"`
	Data T    `json:"data"`
}
