package aliyun

type UcodeEcsPrice struct {
	RegionId            string  `json:"RegionId"`            //区域ID
	InstanceType        string  `json:"InstanceType"`        //实例类型
	TradePrice          float64 `json:"TradePrice"`          //交易金额	越低越好
	AverageSpotDiscount int     `json:"AverageSpotDiscount"` //平均折扣率	越低越好
	InterruptionRate    float64 `json:"InterruptionRate"`    //释放率 越低越好
	ZoneId              string  `json:"ZoneId"`              //可用区
}

/*
*
积分算法，获取性价比、稳定度最佳的抢占式服务器
*/
func CompareUcodeEcsPrice(price1 *UcodeEcsPrice, price2 *UcodeEcsPrice) (result int) {
	if price1 == nil || price2 == nil {
		return 0
	}
	score1 := price2.TradePrice * ((100 - price1.InterruptionRate) / 100.0)
	score2 := price1.TradePrice * ((100 - price2.InterruptionRate) / 100.0)
	score1 = score1 + ((100.0 - float64(price1.AverageSpotDiscount)) / 5000.0)
	score2 = score2 + ((100.0 - float64(price2.AverageSpotDiscount)) / 5000.0)
	if score1 > score2 {
		return 1
	} else {
		return -1
	}
}
