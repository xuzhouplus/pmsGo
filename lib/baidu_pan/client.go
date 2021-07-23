package baidu_pan

import "pmsGo/lib/oauth/gateway"

type baiduPan struct {
	oauth *gateway.Baidu
}

func NewClient() {
	var BaiduPanClient = &baiduPan{}
	BaiduPanClient.oauth, _ = gateway.NewBaidu()
}

func (client baiduPan) AuthorizeUrl() {

}
