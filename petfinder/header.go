package petfinder

import (
	"time"
)

type header struct {
	Timestamp struct {
		T time.Time `json:"$t"`
	} `json:"timestamp"`
	Status struct {
		Message struct {
			T string `json:"$t"`
		} `json:"message"`
		Code struct {
			T string `json:"$t"`
		} `json:"code"`
	} `json:"status"`
	Version struct {
		T string `json:"$t"`
	} `json:"version"`
}
