local key = KEYS[1]
local windowSize = tonumber(ARGV[1])
local threshold = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local uniqueId = ARGV[4]

local expired = now - windowSize
-- return key .. " " .. windowSize .. " " .. threshold .. " " .. now .. " " .. uniqueId .. " " .. expired
redis.call("zremrangebyscore", key, 0, expired)
redis.call("pexpire", key, windowSize)
local cnt = redis.call("zcount", key, expired, now)
if cnt < threshold then
    redis.call("zadd", key, now, uniqueId)
    return 1
else
    return 0
end
