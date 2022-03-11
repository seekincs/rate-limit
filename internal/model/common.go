package model

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type commonItemList []string

type RL struct {
	FuncWhitelist map[string]bool  `json:"funcWhitelist"`
	FuncBlacklist map[string]bool  `json:"funcBlacklist"`
	IPWhitelist   map[string]bool  `json:"ipWhitelist"`
	IPBlacklist   map[string]bool  `json:"ipBlacklist"`
	FuncQPS       map[string]int32 `json:"funcQps"`
}

var RateLimitConfig atomic.Value

func (rlconfig *RL) String() string {
	bytes, marshalErr := json.Marshal(rlconfig)
	if marshalErr != nil {
		return fmt.Sprintf("%#+v", rlconfig)
	}
	return string(bytes)
}

// RateLimitRequest ...
type RateLimitRequest struct {
	IP   string `json:"ip"`
	Func string `json:"func"`
}

// RateLimitResponse ...
type RateLimitResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
