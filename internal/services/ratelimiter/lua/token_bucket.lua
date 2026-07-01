-- token_bucket.lua
-- Atomic token-bucket check-and-consume.
--
-- KEYS[1] = redis key for this bucket
-- ARGV[1] = capacity        (max tokens the bucket can hold)
-- ARGV[2] = refill_rate     (tokens added per second)
-- ARGV[3] = now             (current time in unix milliseconds)
-- ARGV[4] = ttl             (seconds; how long an idle bucket may sit before expiring)

local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now_ms = tonumber(ARGV[3])
local ttl = tonumber(ARGV[4])

local data = redis.call("HMGET", key, "tokens", "last_refill_ms")
local tokens = tonumber(data[1])
local last_refill_ms = tonumber(data[2])

-- init state on first request for this key
if tokens == nil or last_refill_ms == nil then
	tokens = capacity
	last_refill_ms = now_ms
end

-- refill based on elapsed time
local delta_seconds = math.max(0, now_ms - last_refill_ms) / 1000
tokens = math.min(capacity, tokens + delta_seconds * refill_rate)

-- check and consume
local allowed = tokens >= 1
if allowed then
	tokens = tokens - 1
end

-- persist state; refresh TTL so idle buckets get cleaned up automatically
redis.call("HSET", key, "tokens", tokens, "last_refill_ms", now_ms)
redis.call("EXPIRE", key, ttl)

if allowed then
	return 1
else
	return 0
end

