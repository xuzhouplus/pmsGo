package nas

import (
	"fmt"
	"pmsGo/lib/nas/gateway"
	"pmsGo/lib/nas/gateway/baidu"
)

type nas struct {
	Type    string
	Gateway gateway.Gateway
}

func NewNas(gatewayType string) (*nas, error) {
	var gatewayInstance gateway.Gateway
	var err error
	switch gatewayType {
	default:
		return nil, fmt.Errorf("不支持的类型：%v", gatewayType)
	case baidu.GatewayType:
		gatewayInstance, err = baidu.NewBaidu()
	}
	if err != nil {
		return nil, err
	}
	service := &nas{}
	service.Type = gatewayType
	service.Gateway = gatewayInstance
	return service, nil
}
