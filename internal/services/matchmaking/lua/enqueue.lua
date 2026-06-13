local queue_key = KEYS[1]
local user_key = KEYS[2]

local user_id = ARGV[1]
local username = ARGV[2]
local rating = tonumber(ARGV[3])
local joined_at = tonumber(ARGV[4])
local time_control = ARGV[5]
local ttl = tonumber(ARGV[6])

-- User is already in queue
if redis.call("EXISTS", user_key) == 1 then
    return redis.error_reply("already_queued")
end

-- Add player to rating-sorted matchmaking queue
redis.call("ZADD", queue_key, rating, user_id)

-- Store player metadata
redis.call("HSET", user_key,
    "user_id", user_id,
    "username", username,
    "rating", rating,
    "joined_at", joined_at,
    "time_control", time_control
)

-- Expire player metadata after TTL
redis.call("EXPIRE", user_key, ttl)

-- Optional: only use this if queue is temporary
-- redis.call("EXPIRE", queue_key, ttl)

return "ok"
