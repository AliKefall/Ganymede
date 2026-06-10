--@diagnostic disable: undefine-global


local key = KEYS[1]

local capacity = tonumber(ARGV[1])

local refill_rate = tonumber(ARGV[2])

local now =tonumber(ARGV[3])

local requested = tonumber(ARGV[4])

local data = redis.call("HMGET", key, "tokens", "last_refill")

local tokens = tonumber(data[1])
local last_refill = tonumber(data[2])

if tokens == nil then
    token = capacity
    last_refill = now
end

local delta = math.max(0, now - last_refill)

local refill = delta * refill_rate
tokens = math.min(capacity, tokens + refill)

local allowed = tokens >= requested

if allowed then
    tokens = tokens - requested
end

redis.call("HMSET", key, "tokens", "last_refill", now)

if allowed then
    return 1
else
    return 0
end
