package errors

import (
	"github.com/seekincs/rate-limit/internal/model"
)

var (
	ResponseSuccess      = model.RateLimitResponse{Code: 0, Message: "OK"}
	ErrInvalidParameters = model.RateLimitResponse{Code: 40001, Message: "Invalid parameters"}
	ErrIPInBlacklist     = model.RateLimitResponse{Code: 40010, Message: "IP is in the blacklist"}
	ErrFuncInBlacklist   = model.RateLimitResponse{Code: 40011, Message: "Func is in the blacklist"}
	ErrFuncRateLimited   = model.RateLimitResponse{Code: 40012, Message: "Func request frequency exceeds limit"}
)
