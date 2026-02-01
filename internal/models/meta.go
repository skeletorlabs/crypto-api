package models

import "time"

type Meta struct {
	UpdatedAt time.Time `json:"updatedAt"`
	Cached    bool      `json:"cached"`
}
