package service

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/seekincs/rate-limit/internal/errors"
	"github.com/seekincs/rate-limit/internal/model"
)

func isInList(itemSet map[string]bool, item string) bool {
	v, ok := itemSet[item]
	return v && ok
}

func RequestLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		byteBody, _ := ioutil.ReadAll(ctx.Request.Body)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(byteBody))
		log.Printf("%s|request body:%s", uuid.New().String(), string(byteBody))
		ctx.Next()
	}
}

func limit(f string, rlconfig *model.RL) bool {
	qps, ok := rlconfig.FuncQPS[f]
	if ok && model.EvalSha(f, qps, 1) {
		return true
	}
	return false
}

func HandleRateLimit(c *gin.Context) {
	var reqBody model.RateLimitRequest
	bindErr := c.BindJSON(&reqBody)
	rlconfig := model.RateLimitConfig.Load().(*model.RL)
	switch {
	case bindErr != nil:
		log.Printf("bindErr:%s", bindErr.Error())
		debug.PrintStack()
		c.JSON(http.StatusOK, errors.ErrInvalidParameters)
	case len(reqBody.IP) <= 0 && len(reqBody.Func) <= 0:
		c.JSON(http.StatusOK, errors.ResponseSuccess)
	case isInList(rlconfig.IPBlacklist, reqBody.IP):
		c.JSON(http.StatusOK, errors.ErrIPInBlacklist)
	case isInList(rlconfig.FuncBlacklist, reqBody.Func):
		c.JSON(http.StatusOK, errors.ErrFuncInBlacklist)
	case isInList(rlconfig.FuncWhitelist, reqBody.Func),
		isInList(rlconfig.IPWhitelist, reqBody.IP):
		c.JSON(http.StatusOK, errors.ResponseSuccess)
	case limit(reqBody.Func, rlconfig):
		c.JSON(http.StatusOK, errors.ErrFuncRateLimited)
	default:
		c.JSON(http.StatusOK, errors.ResponseSuccess)
	}
}
