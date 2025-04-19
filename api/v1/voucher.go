package v1

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"doovvvDP/config"
	"doovvvDP/dal/model"
	"doovvvDP/dal/mysql"
	"doovvvDP/dal/redis"
	"doovvvDP/dto"
	"doovvvDP/utils"

	goredis "github.com/redis/go-redis/v9"
)

var (
	luaStockScript string
	MyIdWorker     *utils.IDWorker
	// 阻塞队列，已弃用
	// blockQueue *blockqueue.BlockingQueue
)

func VoucherServiceInit() {
	luaScriptPath := "resources/seckill.lua"
	scriptBytes, err := os.ReadFile(luaScriptPath)
	if err != nil {
		panic(err)
	}
	luaStockScript = string(scriptBytes)

	MyIdWorker, err = utils.NewIdWorker(int64(config.MyConfig.MainConfig.WorkerId))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			// 从消息队列中取出订单
			list, err := redis.RDB.XReadGroup(redis.RCtx, &goredis.XReadGroupArgs{
				Group:    "g1",
				Consumer: "c1",
				Streams:  []string{"stream.orders", ">"},
				Count:    1,
				Block:    2 * time.Second,
			}).Result()
			// 如果失败，继续循环
			if err != nil || len(list) == 0 {
				continue
			}
			record := list[0].Messages[0]
			fmt.Println("从消息队列中取出订单：", record.Values)
			// 解析订单
			voucherId, _ := strconv.ParseUint(record.Values["voucherID"].(string), 10, 64)
			userId, _ := strconv.ParseUint(record.Values["userID"].(string), 10, 64)
			orderId, _ := strconv.ParseUint(record.Values["id"].(string), 10, 64)
			voucherOrder := &model.VoucherOrder{
				ID:        orderId,
				UserID:    userId,
				VoucherID: voucherId,
			}
			handleVoucherOrder(voucherOrder)
			// 如果成功就下单，并确认
			redis.RDB.XAck(redis.RCtx, "stream.orders", "g1", record.ID)

		}
	}()
}

func QueryVoucherByShopId(shopId int) *dto.Result {
	result := &dto.Result{}
	vouchers := model.QueryVoucherByShopId(shopId)
	return result.OkWithData(vouchers)
}

func AddSeckillVoucher(voucher model.DTOVoucher) *dto.Result {
	result := &dto.Result{}
	model.AddSeckillVoucher(voucher)
	return result.Ok()
}

func SeckillVoucher(voucherId uint64, userId uint64) *dto.Result {
	result := &dto.Result{}
	voucher, err := model.QueryVoucherById(voucherId)
	if err != nil {
		return result.Fail("查询失败")
	}
	// 判断秒杀是否开始
	if time.Now().Before(voucher.BeginTime) {
		return result.Fail("秒杀还未开始")
	}

	if time.Now().After(voucher.EndTime) {
		return result.Fail("秒杀已经结束")
	}
	id, err := MyIdWorker.Generate()
	if err != nil {
		return result.Fail("订单id生成错误")
	}
	// lua脚本向消息队列中添加订单
	tmp_r, err := redis.RDB.Eval(redis.RCtx, luaStockScript,
		[]string{}, []string{fmt.Sprint(voucherId), fmt.Sprint(userId), fmt.Sprint(id)}).Result()
	if err != nil {
		return result.Fail("系统繁忙:" + err.Error())
	}
	r := tmp_r.(int64)
	if r != 0 {
		errStr := ""
		switch r {
		case 1:
			errStr = "库存不足"
		case 2:
			errStr = "您已经购买过该券"
		}
		return result.Fail(errStr)
	}
	return result.OkWithData(id)
}

// func SeckillVoucher(voucherId uint64, userId uint64) *dto.Result {
// 	result := &dto.Result{}
// 	voucher, err := model.QueryVoucherById(voucherId)
// 	if err != nil {
// 		return result.Fail("查询失败")
// 	}
// 	// 判断秒杀是否开始
// 	if time.Now().Before(voucher.BeginTime) {
// 		return result.Fail("秒杀还未开始")
// 	}

// 	if time.Now().After(voucher.EndTime) {
// 		return result.Fail("秒杀已经结束")
// 	}

// 	// if voucher.Stock < 1 {
// 	// 	return result.Fail("秒杀券库存不足")
// 	// }
// 	// lock := redislock.NewRedisLock(fmt.Sprintf("lock:order:%d", userId))
// 	// isLock := lock.TryLock(time.Second * 60)
// 	// if !isLock {
// 	// 	return result.Fail("系统繁忙")
// 	// }
// 	// defer lock.Unlock()
// 	// return createVoucherOrder(voucherId, userId)

// 	tmp_r, err := redis.RDB.Eval(redis.RCtx, luaStockScript,
// 		[]string{fmt.Sprint(voucherId), fmt.Sprint(userId)}, []string{}).Result()
// 	if err != nil {
// 		return result.Fail("系统繁忙:" + err.Error())
// 	}
// 	r := tmp_r.(int64)
// 	if r != 0 {
// 		errStr := ""
// 		switch r {
// 		case 1:
// 			errStr = "库存不足"
// 		case 2:
// 			errStr = "您已经购买过该券"
// 		}
// 		return result.Fail(errStr)
// 	}
// 	// todo: 异步下单
// 	id, err := MyIdWorker.Generate()
// 	if err != nil {
// 		return result.Fail("订单id生成错误")
// 	}
// 	voucherOrder := model.VoucherOrder{
// 		ID:        id,
// 		UserID:    userId,
// 		VoucherID: voucherId,
// 	}
// 	blockQueue.Put(voucherOrder)
// 	return result.OkWithData(id)
// }

// 单机锁
var (
	userLocks    = make(map[uint64]*sync.Mutex)
	lockMapMutex sync.Mutex
)

// func getOrderCreatLock(userId uint64) *sync.Mutex {
// 	lockMapMutex.Lock()
// 	defer lockMapMutex.Unlock()
// 	if _, exists := userLocks[userId]; !exists {
// 		userLocks[userId] = &sync.Mutex{}
// 	}
// 	return userLocks[userId]
// }

func createVoucherOrder(voucherId uint64, userId uint64) *dto.Result {
	result := &dto.Result{}
	// 一人一单
	// lock := getOrderCreatLock(userId)
	// lock.Lock()
	// defer lock.Unlock()
	haveOrder := model.CheckVoucherOrder(userId, voucherId)
	if haveOrder {
		return result.Fail("您已经购买过该券")
	}
	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			result.Fail("秒杀失败")
		}
	}()
	// 删除库存
	err := model.DecreaseStock(tx, voucherId)
	if err != nil {
		tx.Rollback()
		return result.Fail("库存不足")
	}
	// 订单id
	id, err := MyIdWorker.Generate()
	if err != nil {
		tx.Rollback()
		return result.Fail("订单id生成错误")
	}
	voucherOrder := &model.VoucherOrder{
		ID:        id,
		UserID:    userId,
		VoucherID: voucherId,
	}

	err = model.AddVoucherOrder(tx, voucherOrder)
	if err != nil {
		tx.Rollback()
		return result.Fail("订单创建失败")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return result.Fail("订单创建失败")
	}
	// fmt.Println(voucherOrder.ID)
	return result.OkWithData(voucherOrder.ID)
}

func handleVoucherOrder(voucherOrder *model.VoucherOrder) {
	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("秒杀失败")
		}
	}()
	// 删除库存
	voucherId := voucherOrder.VoucherID
	err := model.DecreaseStock(tx, voucherId)
	if err != nil {
		tx.Rollback()
		fmt.Println("库存不足")
	}

	err = model.AddVoucherOrder(tx, voucherOrder)
	if err != nil {
		tx.Rollback()
		fmt.Println("订单创建失败")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		fmt.Println("订单创建失败")
	}
	// fmt.Println(voucherOrder.ID)
}
