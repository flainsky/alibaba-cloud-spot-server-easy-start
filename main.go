package main

import (
	cloudConfig "CloudBuild/internal/conf"
	"CloudBuild/internal/model/aliyun"
	"CloudBuild/internal/service"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func main() {
	//地域ID
	regionId := "cn-zhangjiakou"
	//初始密码
	myPassword := "xxxxxxxxxxxxxx"
	//自定义镜像名
	myImageName := "xxxxxxxxxxxxxx"
	//cpu核心数
	myCpuCount = 8
	//内存大小
	myMemSize = 64
	//公网带宽
	myNewWidth = 10
	//系统盘大小
	myDiskSize = 40
	

	cloudConf, err := cloudConfig.GetCloudConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error conf file: %s \n", err))
	}
	fmt.Println(fmt.Sprintf("载入配置: %v", cloudConf))

	var isCreate = true
	if isCreate {
		ecsPriceArray, err := service.GetSpotAdvice(&cloudConf, regionId, myCpuCount, myMemSize)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		priceSize := len(ecsPriceArray)
		if priceSize <= 0 {
			fmt.Println("无有效的实例类型")
			return
		}

		sort.Slice(ecsPriceArray, func(i, j int) bool {
			return aliyun.CompareUcodeEcsPrice(&ecsPriceArray[i], &ecsPriceArray[j]) > 0
		})
		fmt.Println("查询实例价格数量：", priceSize)
		for i := 0; i < priceSize; i++ {
			fmt.Println("型号 " + ecsPriceArray[i].InstanceType + "	折扣系数 " + strconv.Itoa(ecsPriceArray[i].AverageSpotDiscount) + "	价格 " + fmt.Sprintf("%f", ecsPriceArray[i].TradePrice) + " 释放率" + fmt.Sprintf("%f", ecsPriceArray[i].InterruptionRate) + " 可用区 " + ecsPriceArray[i].ZoneId)
		}

		//---------------
		fmt.Println("-加载自定义镜像-")
		imageArray, err := service.GetDescribeImages(&cloudConf, regionId)
		if err != nil {
			fmt.Println("查询自定义镜像失败：" + err.Error())
			return
		}
		if len(imageArray) <= 0 {
			fmt.Println("无有效的镜像")
			return
		}
		fmt.Println(fmt.Sprintf("镜像数：%d", len(imageArray)))
		var myImage aliyun.Image
		for _, image := range imageArray {
			if strings.Compare(image.ImageName, myImageName) == 0 {
				myImage = image
			}
		}
		fmt.Println(myImage.ImageName)
		createInstanceResponse, err := service.CreateInstance(&cloudConf, regionId, ecsPriceArray[0], myImage, myNewWidth, myDiskSize, myPassword)
		if err != nil {
			fmt.Println("创建实例失败：" + err.Error())
			return
		}
		_, err = service.AllocatePublicIpAddress(&cloudConf, regionId, createInstanceResponse.InstanceID)
		if err != nil {
			fmt.Println("分配公网IP失败：" + err.Error())
			return
		}
		_, err = service.StartInstance(&cloudConf, regionId, createInstanceResponse.InstanceID)
		if err != nil {
			fmt.Println("启动实例失败：" + err.Error())
			return
		}
	} else {
		//删除实例ID对应的云服务器
		_, err = service.DeleteInstance(&cloudConf, regionId, "i-8vb7nlet29z7a47gqbb0")
		if err != nil {
			fmt.Println("删除实例失败：" + err.Error())
			return
		}
	}

}
