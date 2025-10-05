--[[
  Dọn dẹp một lô kết nối từ một service đã chết.
  Được thiết kế để chạy nhiều lần cho đến khi tất cả các kết nối được dọn dẹp.

  KEYS:
    1: service_conns_key (của service đã chết)

  ARGV:
    1: batch_size (số lượng kết nối cần dọn dẹp trong một lần chạy)
--]]

local service_conns_key = KEYS[1]
local batch_size = tonumber(ARGV[1])

-- Lấy một lô connection_id từ tập hợp của service đã chết
local connections = redis.call('SRANDMEMBER', service_conns_key, batch_size)

if #connections == 0 then
  -- Không còn kết nối nào để dọn dẹp
  return 0
end

local cleaned_count = 0
for i, conn_id in ipairs(connections) do
  local conn_key = 'conn:' .. conn_id
  local user_id = redis.call('HGET', conn_key, 'userId')

  if user_id then
    local user_conns_key = 'user_conns:' .. user_id
    -- Xóa kết nối khỏi các tập hợp và hash
    redis.call('SREM', user_conns_key, conn_id)
    redis.call('SREM', service_conns_key, conn_id)
    redis.call('DEL', conn_key)
    cleaned_count = cleaned_count + 1
  else
    -- Nếu hash không tồn tại vì lý do nào đó, chỉ cần xóa khỏi tập hợp
    redis.call('SREM', service_conns_key, conn_id)
  end
end

return cleaned_count