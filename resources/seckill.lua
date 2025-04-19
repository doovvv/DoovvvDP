local voucherID = ARGV[1]
local userID = ARGV[2]
local orderID = ARGV[3]
local stockKey = "seckill:stock:" .. voucherID
local orderKey = "seckill:order:" .. voucherID
-- 判断库存是否充足
if (tonumber(redis.call('get', stockKey)) <= 0) then
    return 1
-- 判断是否重复下单
elseif (redis.call('sismember', orderKey, userID) == 1) then
    return 2
else
    -- 扣减库存
    redis.call('decr', stockKey)
    -- 下单（保存用户）
    redis.call('sadd', orderKey, userID)
    -- 保存订单到消息队列
    redis.call('xadd', 'stream.orders', '*', 'userID', userID, 'voucherID', voucherID, 'id', orderID)
    return 0
end