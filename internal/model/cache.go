package model

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func readConfig(table string, field string) commonItemList {
	stmt, prepareErr := conn.MySQL.Prepare(
		fmt.Sprintf("select %s from %s", field, table),
	)
	if prepareErr != nil {
		log.Printf("prepareErr: %s", prepareErr.Error())
		debug.PrintStack()
		return make(commonItemList, 0)
	}
	rows, queryErr := stmt.Query()
	if queryErr != nil {
		log.Printf("queryErr: %s", queryErr.Error())
		debug.PrintStack()
		return make(commonItemList, 0)
	}
	var (
		itemList commonItemList
		column   string
	)
	for rows.Next() {
		scanErr := rows.Scan(&column)
		if scanErr != nil {
			log.Printf("scanErr: %s", scanErr.Error())
			debug.PrintStack()
			continue
		}
		itemList = append(itemList, column)
	}
	return itemList
}

func list2Set(list commonItemList) map[string]bool {
	itemSet := make(map[string]bool, len(list))
	for _, item := range list {
		itemSet[item] = true
	}
	return itemSet
}

func CacheRateLimitConfig() {
	qpsRows := readConfig("t_func_qps", "concat(f_func, '#', f_qps)")
	var funcQps = make(map[string]int32)
	for _, row := range qpsRows {
		strArr := strings.Split(row, "#")
		if len(strArr) == 2 {
			qps, _ := strconv.Atoi(strArr[1])
			funcQps[strArr[0]] = int32(qps)
		}
	}
	rlconfig := &RL{
		FuncWhitelist: list2Set(readConfig("t_func_whitelist", "f_func")),
		FuncBlacklist: list2Set(readConfig("t_func_blacklist", "f_func")),
		IPWhitelist:   list2Set(readConfig("t_ip_whitelist", "f_ip")),
		IPBlacklist:   list2Set(readConfig("t_ip_blacklist", "f_ip")),
		FuncQPS:       funcQps,
	}

	go func() {
		for f, qps := range funcQps {
			conn.Redis.Set(context.Background(), fmt.Sprintf("%s:token", f), qps, 3*time.Second)
			conn.Redis.Set(context.Background(), fmt.Sprintf("%s:timestamp", f), qps, 0)
		}
	}()

	log.Printf("rlconfig:%s\n", rlconfig.String())
	RateLimitConfig.Store(rlconfig)
}
