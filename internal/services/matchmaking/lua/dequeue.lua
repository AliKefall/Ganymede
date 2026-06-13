local queue_key = KEYS[1]
local user_key = KEYS[2]

local user_id = ARGV[1]

-- Check if user exists in queue
if redis.call("ZSCORE", queue_key, user_id) == false then
    return redis.error_reply("not_in_queue")
end

-- Remove user from queue
redis.call("ZREM", queue_key, user_id)

-- Remove metadata
redis.call("DEL", user_key)

return "ok"
