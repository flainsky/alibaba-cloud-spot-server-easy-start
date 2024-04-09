package service

import (
	"CloudBuild/internal/model"
	"CloudBuild/internal/model/aliyun"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func getClientV2(conf *model.CloudConf, regionId string) (_err error, aliClient *sdk.Client) {
	client, err := sdk.NewClientWithAccessKey(regionId, conf.AliYunConfig.AccessKey, conf.AliYunConfig.AccessSecret)
	if err != nil {
		return errors.New("client init error"), nil
	}
	return nil, client
}

// 获取实例类型数据
func GetInstanceTypes(conf *model.CloudConf, regionId string, cpuNum int, memorySize int) (_err error) {
	_, client := getClientV2(conf, regionId)
	request := requests.NewCommonRequest()                            // 构造一个公共请求
	request.Method = "POST"                                           // 设置请求方式
	request.Product = "Ecs"                                           // 指定产品
	request.Domain = "ecs.aliyuncs.com"                               // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2014-05-26"                                    // 指定产品版本
	request.ApiName = "DescribeInstanceTypes"                         // 指定接口名
	request.QueryParams["MinimumCpuCoreCount"] = strconv.Itoa(cpuNum) // 设置参数值
	request.QueryParams["MaximumCpuCoreCount"] = strconv.Itoa(cpuNum) // 设置参数值
	request.QueryParams["MinimumMemorySize"] = strconv.Itoa(memorySize)
	request.QueryParams["MaximumMemorySize"] = strconv.Itoa(memorySize)
	request.QueryParams["CpuArchitecture"] = "X86"
	request.QueryParams["MaxResults"] = "1024"
	request.TransToAcsRequest() // 把公共请求转化为acs请求
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return errors.New("执行失败")
	}
	//var describeSpotAdviceResponse aliyun.DescribeSpotAdviceResponse
	//err = json.Unmarshal(response.GetHttpContentBytes(), &describeSpotAdviceResponse)
	fmt.Println(response.GetHttpContentString())
	return nil
}

// 获取抢占式类型列表---通过（查询指定地域下，抢占式实例近30天的实例平均释放率、平均折扣率等信息）接口获取
func GetSpotAdvice(conf *model.CloudConf, regionId string, cpuNum int, memorySize int) (resultList []aliyun.UcodeEcsPrice, _err error) {
	_, client := getClientV2(conf, regionId)
	request := requests.NewCommonRequest() // 构造一个公共请求
	request.Method = "POST"                // 设置请求方式
	request.Product = "Ecs"                // 指定产品
	request.Domain = "ecs.aliyuncs.com"    // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2014-05-26"         // 指定产品版本
	request.ApiName = "DescribeSpotAdvice" // 指定接口名
	request.RegionId = regionId
	request.QueryParams["Cores"] = strconv.Itoa(cpuNum)
	request.QueryParams["Memory"] = strconv.Itoa(memorySize)
	request.TransToAcsRequest()
	response, err := client.ProcessCommonRequest(request)
	var ecsPriceArray = make([]aliyun.UcodeEcsPrice, 0)
	if err != nil {
		return ecsPriceArray, errors.New("执行失败")
	}
	var describeSpotAdviceResponse aliyun.DescribeSpotAdviceResponse
	err = json.Unmarshal(response.GetHttpContentBytes(), &describeSpotAdviceResponse)
	if err != nil {
		return ecsPriceArray, errors.New("json解析失败")
	}
	var wg sync.WaitGroup
	var lock sync.Mutex

	for _, spotZone := range describeSpotAdviceResponse.AvailableSpotZones.AvailableSpotZone {
		for _, spotRes := range spotZone.AvailableSpotResources.AvailableSpotResource {
			wg.Add(1)
			go func() {
				priceResult, err := GetSpotPrice(conf, regionId, spotRes.InstanceType, "PayByTraffic", 10, "cloud_efficiency", 20)
				if err == nil {
					lock.Lock()
					ecsPrice := aliyun.UcodeEcsPrice{
						RegionId:            regionId,
						InstanceType:        spotRes.InstanceType,
						TradePrice:          priceResult.PriceInfo.Price.TradePrice,
						AverageSpotDiscount: spotRes.AverageSpotDiscount,
						InterruptionRate:    spotRes.InterruptionRate,
						ZoneId:              spotZone.ZoneID,
					}
					ecsPriceArray = append(ecsPriceArray, ecsPrice)
					defer lock.Unlock()
				}
				wg.Done()
			}()
		}
	}

	wg.Wait()

	return ecsPriceArray, nil
}

/*
*

	@title 获取抢占式报价
	@param InternetChargeType 		PayByBandwidth：按固定带宽计费。 PayByTraffic：按带宽流量计费。
	@param InternetMaxBandwidthOut 		公网出带宽最大值 单位Mbit/s 范围0-100
	@param SystemDisk.Category 		cloud_ssd：SSD 云盘。 cloud_efficiency：高效云盘。
	@param SystemDisk.Size 			系统盘大小，单位为 GiB。取值范围：20~500。

*
*/
func GetSpotPrice(conf *model.CloudConf, regionId string, instanceType string, internetChargeType string, internetMaxBandwidthOut int, systemDiskCategory string, systemDiskSize int) (result aliyun.DescribePriceResponse, _err error) {
	_, client := getClientV2(conf, regionId)
	request := requests.NewCommonRequest() // 构造一个公共请求
	request.Method = "POST"                // 设置请求方式
	request.Product = "Ecs"                // 指定产品
	request.Domain = "ecs.aliyuncs.com"    // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2014-05-26"         // 指定产品版本
	request.ApiName = "DescribePrice"      // 指定接口名
	request.RegionId = regionId
	request.QueryParams["InstanceType"] = instanceType
	request.QueryParams["InternetChargeType"] = internetChargeType
	request.QueryParams["InternetMaxBandwidthOut"] = strconv.Itoa(internetMaxBandwidthOut)
	request.QueryParams["SystemDiskCategory"] = systemDiskCategory
	request.QueryParams["SystemDiskSize"] = strconv.Itoa(systemDiskSize)
	request.QueryParams["Period"] = "1"
	request.QueryParams["PriceUnit"] = "Hour"
	request.QueryParams["SpotStrategy"] = "SpotWithPriceLimit"
	request.QueryParams["SpotDuration"] = "1"
	request.TransToAcsRequest()
	response, err := client.ProcessCommonRequest(request)
	var describePriceResponse aliyun.DescribePriceResponse
	if err != nil {
		return describePriceResponse, errors.New("执行失败")
	}

	err = json.Unmarshal(response.GetHttpContentBytes(), &describePriceResponse)
	if err != nil {
		return describePriceResponse, errors.New("json解析失败")
	}
	return describePriceResponse, nil
}

// 获取自定义镜像列表
func GetDescribeImages(conf *model.CloudConf, regionId string) (result []aliyun.Image, _err error) {
	_, client := getClientV2(conf, regionId)
	request := requests.NewCommonRequest()          // 构造一个公共请求
	request.Method = "POST"                         // 设置请求方式
	request.Product = "Ecs"                         // 指定产品
	request.Domain = "ecs.aliyuncs.com"             // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2014-05-26"                  // 指定产品版本
	request.ApiName = "DescribeImages"              // 指定接口名
	request.QueryParams["ImageOwnerAlias"] = "self" //镜像来源：自定义镜像
	request.RegionId = regionId
	request.TransToAcsRequest()
	response, err := client.ProcessCommonRequest(request)

	var describeImagesResponse aliyun.DescribeImagesResponse
	if err != nil {
		return nil, errors.New("执行失败")
	}

	err = json.Unmarshal(response.GetHttpContentBytes(), &describeImagesResponse)
	if err != nil {
		return nil, errors.New("json解析失败")
	}
	return describeImagesResponse.Images.Image, nil
}

// 创建ECS实例
// 需要自己调用 AllocatePublicIpAddress 分配公网IP
func CreateInstance(conf *model.CloudConf, regionId string, ecsPrice aliyun.UcodeEcsPrice, imageInfo aliyun.Image, netWidth int, diskSize int, password string) (result aliyun.CreateInstanceResponse, _err error) {
	_, client := getClientV2(conf, regionId)
	request := requests.NewCommonRequest() // 构造一个公共请求
	request.Method = "POST"                // 设置请求方式
	request.Product = "Ecs"                // 指定产品
	request.Domain = "ecs.aliyuncs.com"    // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2014-05-26"         // 指定产品版本
	request.ApiName = "CreateInstance"     // 指定接口名
	request.RegionId = regionId
	request.QueryParams["ImageId"] = imageInfo.ImageID
	request.QueryParams["InstanceType"] = ecsPrice.InstanceType
	//request.QueryParams["SecurityGroupId"] = "" //安全组ID
	request.QueryParams["InstanceName"] = "ucode-auto-ecs-plot-" + strconv.FormatInt(time.Now().Unix(), 10)
	request.QueryParams["InternetMaxBandwidthIn"] = strconv.Itoa(netWidth)
	request.QueryParams["InternetMaxBandwidthOut"] = strconv.Itoa(netWidth)
	request.QueryParams["Password"] = password
	request.QueryParams["ZoneId"] = ecsPrice.ZoneId
	request.QueryParams["SystemDisk.Size"] = strconv.Itoa(diskSize)
	request.QueryParams["SystemDisk.Category"] = "cloud_efficiency"
	request.QueryParams["SpotStrategy"] = "SpotWithPriceLimit"
	request.QueryParams["SpotPriceLimit"] = fmt.Sprintf("%f", (ecsPrice.TradePrice * 1.5))
	request.TransToAcsRequest()
	response, err := client.ProcessCommonRequest(request)

	var createInstanceResponse aliyun.CreateInstanceResponse
	if err != nil {
		return createInstanceResponse, errors.New("执行失败")
	}
	err = json.Unmarshal(response.GetHttpContentBytes(), &createInstanceResponse)
	if err != nil {
		return createInstanceResponse, errors.New("json解析失败")
	}
	return createInstanceResponse, nil
}

// 创建公网IP
func AllocatePublicIpAddress(conf *model.CloudConf, regionId string, instanceId string) (result aliyun.AllocatePublicIpAddressResponse, _err error) {

	_, client := getClientV2(conf, regionId)
	request := requests.NewCommonRequest()      // 构造一个公共请求
	request.Method = "POST"                     // 设置请求方式
	request.Product = "Ecs"                     // 指定产品
	request.Domain = "ecs.aliyuncs.com"         // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2014-05-26"              // 指定产品版本
	request.ApiName = "AllocatePublicIpAddress" // 指定接口名
	request.RegionId = regionId
	request.QueryParams["InstanceId"] = instanceId
	request.TransToAcsRequest()
	response, err := client.ProcessCommonRequest(request)

	var allocatePublicIpAddressResponse aliyun.AllocatePublicIpAddressResponse
	if err != nil {
		return allocatePublicIpAddressResponse, errors.New("执行失败")
	}
	err = json.Unmarshal(response.GetHttpContentBytes(), &allocatePublicIpAddressResponse)
	if err != nil {
		return allocatePublicIpAddressResponse, errors.New("json解析失败")
	}
	return allocatePublicIpAddressResponse, nil
}

// 启动实例
func StartInstance(conf *model.CloudConf, regionId string, instanceId string) (result aliyun.StartInstanceResponse, _err error) {
	_, client := getClientV2(conf, regionId)
	request := requests.NewCommonRequest() // 构造一个公共请求
	request.Method = "POST"                // 设置请求方式
	request.Product = "Ecs"                // 指定产品
	request.Domain = "ecs.aliyuncs.com"    // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2014-05-26"         // 指定产品版本
	request.ApiName = "StartInstance"      // 指定接口名
	request.RegionId = regionId
	request.QueryParams["InstanceId"] = instanceId
	request.TransToAcsRequest()
	response, err := client.ProcessCommonRequest(request)

	var startInstanceResponse aliyun.StartInstanceResponse
	if err != nil {
		return startInstanceResponse, errors.New("执行失败")
	}
	err = json.Unmarshal(response.GetHttpContentBytes(), &startInstanceResponse)
	if err != nil {
		return startInstanceResponse, errors.New("json解析失败")
	}
	return startInstanceResponse, nil
}

func DeleteInstance(conf *model.CloudConf, regionId string, instanceId string) (result aliyun.DeleteInstanceResponse, _err error) {
	_, client := getClientV2(conf, regionId)
	request := requests.NewCommonRequest() // 构造一个公共请求
	request.Method = "POST"                // 设置请求方式
	request.Product = "Ecs"                // 指定产品
	request.Domain = "ecs.aliyuncs.com"    // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2014-05-26"         // 指定产品版本
	request.ApiName = "DeleteInstance"     // 指定接口名
	request.RegionId = regionId
	request.QueryParams["InstanceId"] = instanceId
	request.TransToAcsRequest()
	response, err := client.ProcessCommonRequest(request)

	var deleteInstanceResponse aliyun.DeleteInstanceResponse
	if err != nil {
		return deleteInstanceResponse, errors.New("执行失败")
	}
	err = json.Unmarshal(response.GetHttpContentBytes(), &deleteInstanceResponse)
	if err != nil {
		return deleteInstanceResponse, errors.New("json解析失败")
	}
	return deleteInstanceResponse, nil
}

/*
*

	@title 自动创建抢占式实例
	@param conf 		阿里云配置
	@param region 		地区ID
	@param cpuNum 		cpu数量
	@param memorySize 	内存大小
	@param InternetChargeType 		PayByBandwidth：按固定带宽计费。 PayByTraffic：按带宽流量计费。
	@param InternetMaxBandwidthOut 		公网出带宽最大值 单位Mbit/s 范围0-100
	@param SystemDisk.Category 		cloud_ssd：SSD 云盘。 cloud_efficiency：高效云盘。
	@param SystemDisk.Size 			系统盘大小，单位为 GiB。取值范围：20~500。

*
*/
func CreateIntance4AutoSpot(conf *model.CloudConf, regionId string, cpuNum int, memorySize int, internetChargeType string, internetMaxBandwidthOut int, systemDiskCategory string, systemDiskSize int) {
	ecsPriceArray, err := GetSpotAdvice(conf, regionId, 8, 64)
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
	imageArray, err := GetDescribeImages(conf, regionId)
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
		if strings.Compare(image.ImageName, "omind-charge-cloud-2.0.6") == 0 {
			myImage = image
		}
	}
	fmt.Println(myImage.ImageName)
	createInstanceResponse, err := CreateInstance(conf, regionId, ecsPriceArray[0], myImage, internetMaxBandwidthOut, systemDiskSize, "YMKJ-soft-!@#$")
	if err != nil {
		fmt.Println("创建实例失败：" + err.Error())
		return
	}
	_, err = AllocatePublicIpAddress(conf, regionId, createInstanceResponse.InstanceID)
	if err != nil {
		fmt.Println("分配公网IP失败：" + err.Error())
		return
	}
	_, err = StartInstance(conf, regionId, createInstanceResponse.InstanceID)
	if err != nil {
		fmt.Println("启动实例失败：" + err.Error())
		return
	}
}
