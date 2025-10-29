--[[
  Nguyên tử hóa việc tạo kết nối mới và kiểm tra giới hạn.
  Trả về 1 nếu thành công, 0 nếu vượt quá giới hạn kết nối.

  KEYS:
    1: device_conns_key (ví dụ: device_conns:device123)
    2: service_conns_key (ví dụ: service_conns:notif-A)

  ARGV:
    1: connection_id
    2: device_id
    3: service_id
    4: ip_address
    5: connected_at (timestamp)
    6: user_agent
--]]

local device_conns_key = KEYS[1]
local service_conns_key = KEYS[2]

local connection_id = ARGV[1]
local device_id = ARGV[2]
local service_id = ARGV[3]
local ip_address = ARGV[4]
local connected_at = ARGV[5]
local user_agent = ARGV[6]

-- Add device id to service
redis.call("SADD", service_conns_key, device_id)
-- Save device connection
redis.call("HSET", device_conns_key,
    "connection_id", connection_id,
    "service_id", service_id,
    "ip_address", ip_address,
    "connected_at", connected_at,
    "user_agent", user_agent
)

return 1