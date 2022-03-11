package configs

const (
	// Do not hard-code database accounts and passwords in the production environment,
	// as there will be security risks, instead they should be placed in a place like the configuration center.
	DSN       = "root:passwd@tcp(localhost:3306)/db_rate_limit"
	RedisAddr = "localhost:6379"
	// By periodically reading the configuration into the cache,
	// the efficiency and performance of access will be improved.
	RefreshConfigCron = "@every 1m"
	// Borrow from https://github.com/wallstreetcn/rate/blob/master/redis/rate.go
	RateLimitEvalScript = `
	local token_bucket_key = KEYS[1]
	local timestamp_key = KEYS[2]

	local rate = tonumber(ARGV[1])
	local capacity = tonumber(ARGV[2])
	local current_timestamp = tonumber(ARGV[3])
	local requested_tokens = tonumber(ARGV[4])

	local last_tokens = tonumber(redis.call("get", token_bucket_key))
	if last_tokens == nil then
		last_tokens = capacity
	end

	local last_refreshed_timestamp = tonumber(redis.call("get", timestamp_key))
	if last_refreshed_timestamp == nil then
		last_refreshed_timestamp = 0
	end

	local delta = math.max(0, current_timestamp-last_refreshed_timestamp)
	local filled_tokens = math.min(capacity, last_tokens+(delta*rate))
	local allowed = filled_tokens >= requested_tokens
	local new_tokens = filled_tokens
	if allowed then
		new_tokens = filled_tokens - requested_tokens
	end

	redis.call("set", token_bucket_key, new_tokens)
	redis.call("set", timestamp_key, current_timestamp)

	return { allowed, new_tokens }
	`
)
