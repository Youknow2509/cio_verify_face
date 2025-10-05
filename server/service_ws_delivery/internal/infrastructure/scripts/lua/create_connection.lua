--[[
  Nguyên tử hóa việc tạo kết nối mới và kiểm tra giới hạn.
  Trả về 1 nếu thành công, 0 nếu vượt quá giới hạn kết nối.

  KEYS:
    1: user_conns_key (ví dụ: user_conns:user123)
    2: service_conns_key (ví dụ: service_conns:notif-A)
    3: conn_key (ví dụ: conn:uuid-xyz)

  ARGV:
    1: connection_id
    2: user_id
    3: service_id
    4: ip_address
    5: connected_at (timestamp)
    6: user_agent
    7: max_conns_per_user (giới hạn số kết nối)
--]]

local user_conns_key = KEYS[1]
local service_conns_key = KEYS[2]
local conn_key = KEYS[3]

local connection_id = ARGV[1]
local user_id = ARGV[2]
local service_id = ARGV[3]
local ip_address = ARGV[4]
local connected_at = ARGV[5]
local user_agent = ARGV[6]
local max_conns_per_user = tonumber(ARGV[7])

-- Bước 1: Kiểm tra giới hạn kết nối của người dùng
local current_conns = redis.call('SCARD', user_conns_key)
if current_conns == nil then
  current_conns = 0
end
if current_conns >= max_conns_per_user then
  return 0
end

-- Bước 2: Thêm kết nối vào các tập hợp và tạo hash
redis.call('HSET', conn_key,
  'userId', user_id,
  'serviceId', service_id,
  'ipAddress', ip_address,
  'connectedAt', connected_at,
  'userAgent', user_agent
)
redis.call('SADD', user_conns_key, connection_id)
redis.call('SADD', service_conns_key, connection_id)

return 1