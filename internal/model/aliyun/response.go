package aliyun

import "time"

// ---------查询指定地域下，抢占式实例近30天的实例平均释放率、平均折扣率等信息
type DescribeSpotAdviceResponse struct {
	RequestID          string             `json:"RequestId"`
	AvailableSpotZones AvailableSpotZones `json:"AvailableSpotZones"`
	RegionID           string             `json:"RegionId"`
}
type AvailableSpotResource struct {
	InterruptRateDesc   string  `json:"InterruptRateDesc"`
	InstanceType        string  `json:"InstanceType"`
	AverageSpotDiscount int     `json:"AverageSpotDiscount"`
	InterruptionRate    float64 `json:"InterruptionRate"`
}
type AvailableSpotResources struct {
	AvailableSpotResource []AvailableSpotResource `json:"AvailableSpotResource"`
}
type AvailableSpotZone struct {
	ZoneID                 string                 `json:"ZoneId"`
	AvailableSpotResources AvailableSpotResources `json:"AvailableSpotResources"`
}
type AvailableSpotZones struct {
	AvailableSpotZone []AvailableSpotZone `json:"AvailableSpotZone"`
}

// ------查询云服务器ECS资源的最新价格
type DescribePriceResponse struct {
	RequestID string    `json:"RequestId"`
	PriceInfo PriceInfo `json:"PriceInfo"`
}
type Price struct {
	OriginalPrice             float64 `json:"OriginalPrice"`
	ReservedInstanceHourPrice float64 `json:"ReservedInstanceHourPrice"`
	DiscountPrice             float64 `json:"DiscountPrice"`
	Currency                  string  `json:"Currency"`
	TradePrice                float64 `json:"TradePrice"`
}
type Rule struct {
	Description string `json:"Description"`
	RuleID      int    `json:"RuleId"`
}
type Rules struct {
	Rule []Rule `json:"Rule"`
}
type PriceInfo struct {
	Price Price `json:"Price"`
	Rules Rules `json:"Rules"`
}

// --------- 查询自定义镜像返回
type DescribeImagesResponse struct {
	TotalCount int    `json:"TotalCount"`
	PageSize   int    `json:"PageSize"`
	RequestID  string `json:"RequestId"`
	PageNumber int    `json:"PageNumber"`
	Images     Images `json:"Images"`
	RegionID   string `json:"RegionId"`
}

type Features struct {
	NvmeSupport string `json:"NvmeSupport"`
}
type Tag struct {
	TagKey   string `json:"TagKey"`
	TagValue string `json:"TagValue"`
}
type Tags struct {
	Tag []Tag `json:"Tag"`
}
type DiskDeviceMapping struct {
	SnapshotID      string `json:"SnapshotId"`
	Type            string `json:"Type"`
	Progress        string `json:"Progress"`
	Format          string `json:"Format"`
	Device          string `json:"Device"`
	Size            string `json:"Size"`
	ImportOSSBucket string `json:"ImportOSSBucket"`
	ImportOSSObject string `json:"ImportOSSObject"`
}
type DiskDeviceMappings struct {
	DiskDeviceMapping []DiskDeviceMapping `json:"DiskDeviceMapping"`
}
type Item struct {
	RiskCode  string `json:"RiskCode"`
	Value     string `json:"Value"`
	RiskLevel string `json:"RiskLevel"`
	Name      string `json:"Name"`
}
type Items struct {
	Item []Item `json:"Item"`
}
type DetectionOptions struct {
	Status string `json:"Status"`
	Items  Items  `json:"Items"`
}
type Image struct {
	ImageOwnerAlias         string             `json:"ImageOwnerAlias"`
	IsSelfShared            string             `json:"IsSelfShared"`
	Description             string             `json:"Description"`
	Platform                string             `json:"Platform"`
	ResourceGroupID         string             `json:"ResourceGroupId"`
	Size                    int                `json:"Size"`
	IsSubscribed            bool               `json:"IsSubscribed"`
	BootMode                string             `json:"BootMode"`
	OSName                  string             `json:"OSName"`
	IsPublic                bool               `json:"IsPublic"`
	ImageID                 string             `json:"ImageId"`
	DetectionOptions        DetectionOptions   `json:"DetectionOptions,omitempty"`
	Features                Features           `json:"Features"`
	OSNameEn                string             `json:"OSNameEn"`
	Tags                    Tags               `json:"Tags"`
	LoginAsNonRootSupported bool               `json:"LoginAsNonRootSupported"`
	Status                  string             `json:"Status"`
	Progress                string             `json:"Progress"`
	Usage                   string             `json:"Usage"`
	Architecture            string             `json:"Architecture"`
	ProductCode             string             `json:"ProductCode"`
	IsCopied                bool               `json:"IsCopied"`
	ImageFamily             string             `json:"ImageFamily"`
	IsSupportIoOptimized    bool               `json:"IsSupportIoOptimized"`
	IsSupportCloudinit      bool               `json:"IsSupportCloudinit"`
	ImageName               string             `json:"ImageName"`
	DiskDeviceMappings      DiskDeviceMappings `json:"DiskDeviceMappings"`
	ImageVersion            string             `json:"ImageVersion"`
	OSType                  string             `json:"OSType"`
	CreationTime            time.Time          `json:"CreationTime"`
}
type Images struct {
	Image []Image `json:"Image"`
}

// CreateInstance
type CreateInstanceResponse struct {
	RequestID  string `json:"RequestId"`
	InstanceID string `json:"InstanceId"`
}

// AllocatePublicIpAddressResponse
type AllocatePublicIpAddressResponse struct {
	RequestID string `json:"RequestId"`
	IpAddress string `json:"IpAddress"`
}

type StartInstanceResponse struct {
	RequestID string `json:"RequestId"`
}

type DeleteInstanceResponse struct {
	RequestID string `json:"RequestId"`
}
