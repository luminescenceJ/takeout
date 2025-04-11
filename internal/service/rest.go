package service

import (
	"context"
	"fmt"
	"log"
	"takeout/global"
	"takeout/internal/model"
	"takeout/repository/dao"
)

const stockDeductLua = `
  local stock = redis.call("GET", KEYS[1])
  if not stock then
    return -1 -- 没有这个商品库存
  end
  stock = tonumber(stock)
  local count = tonumber(ARGV[1])
  if stock < count then
    return 0 -- 库存不足
  else
    redis.call("DECRBY", KEYS[1], count)
    return 1 -- 扣减成功
  end
`

var ctx = context.Background()

// 预扣库存（使用Lua）
func PreDeductStock(dishId int64, count int) (bool, error) {
	key := fmt.Sprintf("dish_stock:%d", dishId)
	client := global.RedisClient

	result, err := client.Eval(stockDeductLua, []string{key}, count).Int()
	if err != nil {
		return false, err
	}

	if result == 1 {
		log.Printf("库存扣减成功，菜品ID：%d，数量：%d\n", dishId, count)
		return true, nil
	} else if result == 0 {
		log.Printf("库存不足，菜品ID：%d\n", dishId)
		return false, nil
	} else {
		log.Printf("菜品不存在，ID：%d\n", dishId)
		return false, fmt.Errorf("菜品不存在")
	}
}

// 业务成功后的库存回写逻辑
func ConfirmStockToDB(dishId int64, count int) error {
	err := models.ReduceDishStock(dishId, count)
	if err != nil {
		log.Printf("数据库库存扣减失败，菜品ID：%d\n", dishId)
		return err
	}
	log.Printf("数据库库存扣减成功，菜品ID：%d，数量：%d\n", dishId, count)
	return nil
}

func OrderDish(dishId int64, count int) error {
	ok, err := PreDeductStock(dishId, count)
	if err != nil || !ok {
		return fmt.Errorf("库存不足或预扣失败")
	}

	success := dao.DishDao{}.Insert(global.DB, &model.Dish{})
	if !success {
		// TODO: 可以放到延迟队列补回库存
		log.Println("订单失败，需补偿库存")
		return fmt.Errorf("订单创建失败")
	}

	// 3. 扣除数据库库存
	err = ConfirmStockToDB(dishId, count)
	if err != nil {
		return fmt.Errorf("数据库扣减失败")
	}

	return nil
}
