--[[
  Thực hiện ngắt kết nối một cách nguyên tử.

  KEYS:
    1: user_conns_key
    2: service_conns_key
    3: conn_key

  ARGV:
    1: connection_id
--]]

local user_conns_key = KEYS[1]
local service_conns_key = KEYS[2]
local conn_key = KEYS[3]
local connection_id = ARGV[1]

redis.call('SREM', user_conns_key, connection_id)
redis.call('SREM', service_conns_key, connection_id)
redis.call('DEL', conn_key)

return 1