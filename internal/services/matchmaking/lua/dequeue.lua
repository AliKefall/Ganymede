local queue_key = KEYS[1]
local user_key = KEYS[2]

local user_id = ARGV[1]

local removed = redis.call(
    "ZREM",
    queue_key,
    user_id
)

redis.call(
    "DEL",
    user_key
)

return removed
