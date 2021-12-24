package model

import (
	"gorm.io/gorm"
	"pmsGo/lib/database"
	"time"
)

const (
	ConnectStatusEnable  = 1
	ConnectStatusDisable = 2
)

const (
	ConnectTypeAlipay   = "alipay"
	ConnectTypeBaidu    = "baidu"
	ConnectTypeFacebook = "facebook"
	ConnectTypeGitee    = "gitee"
	ConnectTypeGithub   = "github"
	ConnectTypeGoogle   = "google"
	ConnectTypeLine     = "line"
	ConnectTypeQq       = "qq"
	ConnectTypeTwitter  = "twitter"
	ConnectTypeWechat   = "wechat"
	ConnectTypeWeibo    = "weibo"
)

type Connect struct {
	ID        int       `gorm:"private_key" json:"id"`
	AdminId   int       `json:"admin_id"`
	Type      string    `json:"type"`
	Avatar    string    `json:"avatar"`
	Account   string    `json:"account"`
	UnionId   string    `json:"_"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (model *Connect) DB() *gorm.DB {
	return database.DB.Model(&model)
}
