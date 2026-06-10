local queue_key = KEYS[1]
local user_key = KEYS[2]

local user_id = ARGV[1]
local rating = tonumber(ARGV[2])
local joined_at = ARGV[3]
local time_control = ARGV[4]

if redis.call("EXISTS", user_key) == 1 then
    return {
        err = "already_queued"
    }
end

redis.call(
    "ZADD", -- sorted set
    queue_key,
    rating,
    user_id
)

redis.call(
    "HSET",
    user_key,
    "rating", rating,
    "joined_at", joined_at,
    "time_control", time_control
)

redis.call(
    "EXPIRE",
    user_key,
    300
)

return "ok"
