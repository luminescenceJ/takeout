package service

import "takeout/global"

type ShopService struct{}

const ShopStatus = "shop_status"

func (service *ShopService) SetShopStatus(status string) error {
	if err := global.RedisClient.Set(ShopStatus, status, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (service *ShopService) GetShopStatus() (string, error) {
	status, err := global.RedisClient.Get(ShopStatus).Result()
	if err != nil {
		return "", err
	}
	return status, nil
}
