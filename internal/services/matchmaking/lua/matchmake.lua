local queue_key = KEYS[1]

local now = tonumber(ARGV[1])
local base_rating = tonumber(ARGV[2])
local rating_speed = tonumber(ARGV[3])

local min_rating = base_rating - rating_speed
local max_rating = base_rating + rating_speed

local candidates = redis.call("ZRANGEBYSCORE", queue_key, min_rating, max_rating, "WITHSCORES")

if #candidates < 4 then
	return nil
end

local p1_id = candidates[1]
local p1_rating = candidates[2]

local p2_id = candidates[3]
local p2_rating = candidates[4]

if p1_id == p2_id then
	return nil
end

redis.call("ZREM", queue_key, p1_id)
redis.call("ZREM", queue_key, p2_id)

redis.call("DEL", "queue:user:" .. p1_id)
redis.call("DEL", "queue:user:" .. p2_id)

local match_id = p1_id .. ":" .. p2_id

return {
	match_id,
	p1_id,
	p2_id,
	p1_rating,
	p2_rating,
}
