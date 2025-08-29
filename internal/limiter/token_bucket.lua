
-- user Id stored as key in hash
local key = KEYS[1]

-- args
local max_tokens = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local current_timestamp = tonumber(ARGV[3])


-- get current data
local data = redis.call('HMGET', key, 'tokens', 'ts')
local current_tokens = tonumber(data[1])
local last_refill_timestamp = tonumber(data[2])

-- no data for key, (1st req) init
if current_tokens == nil then
    current_tokens = max_tokens
    last_refill_timestamp = current_timestamp
end


local elapsed_time = current_timestamp - last_refill_timestamp
local new_tokens = elapsed_time * refill_rate

-- add new tokens
if new_tokens > 0 then
    current_tokens = current_tokens + new_tokens
    last_refill_timestamp = current_timestamp
end

-- max cond
if current_tokens > max_tokens then
    current_tokens = max_tokens
end

-- check if allowed and update
local allowed = 0
if current_tokens >= 1 then
    current_tokens = current_tokens - 1
    allowed = 1
end

-- set
redis.call('HMSET',key, 'tokens',current_tokens,'ts',last_refill_timestamp)
redis.call('EXPIRE', key, math.ceil(max_tokens/refill_rate)*2) -- expiry after 2 refills (clearup old users)

return {allowed,current_tokens}