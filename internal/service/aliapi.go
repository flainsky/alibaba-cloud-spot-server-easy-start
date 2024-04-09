package service

import (
	"CloudBuild/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/duke-git/lancet/v2/strutil"
)

func getClient(conf *model.CloudConf) (_result *ecs20140526.Client, _err error) {
	// 1. 初始化配置
	config := &openapi.Config{
		AccessKeyId:     tea.String(conf.AliYunConfig.AccessKey),
		AccessKeySecret: tea.String(conf.AliYunConfig.AccessSecret),
	}
	// 设置请求地址
	config.Endpoint = tea.String("ecs.aliyuncs.com")
	// 设置连接超时为5000毫秒
	config.ConnectTimeout = tea.Int(5000)
	// 设置读超时为5000毫秒
	config.ReadTimeout = tea.Int(5000)

	_result = &ecs20140526.Client{}
	_result, _err = ecs20140526.NewClient(config)
	return _result, _err
}

// 获取自有实例列表
func GetMyInstances(conf *model.CloudConf, regionCode string) (_err error, ecsArray []*ecs20140526.DescribeInstancesResponseBodyInstancesInstance) {

	if conf == nil {
		return errors.New("config is null"), nil
	}
	if strutil.IsBlank(regionCode) {
		return errors.New("regionCode is empty"), nil
	}
	client, _err := getClient(conf)
	if _err != nil {
		return _err, nil
	}
	describeInstancesRequest := &ecs20140526.DescribeInstancesRequest{
		PageSize: tea.Int32(100),
		RegionId: &regionCode,
	}
	resp, _err := client.DescribeInstances(describeInstancesRequest)
	if _err != nil {
		return _err, nil
	}

	instances := resp.Body.Instances.Instance
	fmt.Println(regionCode + " 下 ECS 实例列表:")
	for _, instance := range instances {
		fmt.Println("  " + tea.StringValue(instance.HostName) + " 实例ID " + tea.StringValue(instance.InstanceId) + " CPU:" + tea.ToString(tea.Int32Value(instance.Cpu)) + "  内存:" + tea.ToString(tea.Int32Value(instance.Memory)) + " MB 规格：" + tea.StringValue(instance.InstanceType) + " 系统:" + tea.StringValue(instance.OSType) + "(" + tea.StringValue(instance.OSName) + ") 状态：" + tea.StringValue(instance.Status))
	}
	return nil, instances
}

// 获取可用的实例类型
// CpuArchitecture：CPU架构 值：X86 ARM
// MinimumCpuCoreCount
// MaximumCpuCoreCount
// MinimumMemorySize
// MaximumMemorySize
func GetIntanceTypes(conf *model.CloudConf, regionCode string, cpuNum int, memorySize int) (_err error) {
	if conf == nil {
		return errors.New("config is null")
	}
	if strutil.IsBlank(regionCode) {
		return errors.New("regionCode is empty")
	}

	client, _err := getClient(conf)
	if _err != nil {
		return _err
	}
	describeInstanceTypesRequest := &ecs20140526.DescribeInstanceTypesRequest{
		MaxResults: tea.Int64(1024),
	}
	resp, _err := client.DescribeInstanceTypes(describeInstanceTypesRequest)
	if _err != nil {
		return _err
	}
	instanceTypeArray := resp.Body.InstanceTypes.InstanceType
	for _, instanceType := range instanceTypeArray {
		b, _ := json.Marshal(instanceType)
		fmt.Println(string(b))
	}
	return nil
}
