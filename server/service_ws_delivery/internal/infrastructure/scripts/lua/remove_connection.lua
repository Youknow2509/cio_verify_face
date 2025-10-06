--[[
  Thực hiện ngắt kết nối một cách nguyên tử.

  KEYS:
    1: device_conns_key
    2: service_conns_key

  ARGV:
    1: device_id
--]]

local device_conns_key = KEYS[1]
local service_conns_key = KEYS[2]
local device_id = ARGV[1]

redis.call("DEL", device_conns_key)
redis.call("SREM", service_conns_key, device_id)

return 1