package v1

import (
	"context"
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
	messagequeue "doovvvDP/messageQueue"
	"doovvvDP/utils"

	"github.com/segmentio/kafka-go"
)

var (
	luaStockScript string
	MyIdWorker     *utils.IDWorker
	// 阻塞队列，已弃用
	// blockQueue *blockqueue.BlockingQueue

	// kafka消息队列
	kafkaQueue *messagequeue.KafkaService
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

	// 初始化消息队列
	kafkaQueue = messagequeue.NewKafkaService()
	kafkaQueue.KafkaInit()

	go func() {
		for {
			// 从消息队列中取出订单
			msg, err := kafkaQueue.OrderConsumer.ReadMessage(context.Background())
			// 如果失败，继续循环
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("接收到消息：", string(msg.Value))
			var voucherId, userId, orderId uint64
			n, err := fmt.Sscanf(string(msg.Value), "voucherID:%d,userID:%d,id:%d", &voucherId, &userId, &orderId)
			if err != nil || n != 3 {
				fmt.Println(err.Error())
				continue
			}
			voucherOrder := &model.VoucherOrder{
				ID:        orderId,
				UserID:    userId,
				VoucherID: voucherId,
			}
			handleVoucherOrder(voucherOrder)
			// 如果成功就下单，并确认
			err = kafkaQueue.OrderConsumer.CommitMessages(context.Background(), msg)
			if err != nil {
				continue
			}

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
	// 数据库查询秒杀券信息
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
	// lua脚本检测订单
	// fmt.Println("voucherID:", voucherId, " userID:", userId, " id:", id)
	tmp_r, err := redis.RDB.Eval(redis.RCtx, luaStockScript,
		[]string{}, []string{fmt.Sprint(voucherId), fmt.Sprint(userId)}).Result()
	if err != nil {
		fmt.Println(err.Error())
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
	// 往消息队列中添加订单
	err = kafkaQueue.OrderProducer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(strconv.FormatUint(id, 10)),
		Value: []byte(fmt.Sprintf("voucherID:%d,userID:%d,id:%d", voucherId, userId, id)),
	})
	if err != nil {
		return result.Fail("系统繁忙:" + err.Error())
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
