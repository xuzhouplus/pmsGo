package router

import (
	"pmsGo/controller"
	"pmsGo/lib/core"
)

func Router(engine *core.Engine) {
	//支持上传文件的支持
	engine.Static("/public", "./public")
	engine.Router(controller.Admin)
	engine.Router(controller.Setting)
	engine.Router(controller.Post)
	engine.Router(controller.File)
	engine.Router(controller.Carousel)
	/**
	//账号路由分组
	admin := engine.Group("/admin")
	{
		//账号登录
		admin.POST("/login", controller.Admin.Login)
		//获取登录用户信息
		admin.GET("/auth", controller.Admin.Auth)
		//登出
		admin.POST("/logout", controller.Admin.Logout)
		//获取账号信息
		admin.GET("/profile", controller.Admin.Profile)
		//保存账号信息
		admin.POST("/profile", controller.Admin.Profile)
		//获取第三方互联账号
		admin.GET("/connects", controller.Admin.Connects)
		//获取第三方互联授权地址
		admin.GET("/authorize", controller.Admin.AuthorizeUrl)
		admin.GET("/authorize-url", controller.Admin.AuthorizeUrl)
		//获取第三方互联授权账号
		admin.GET("/callback", controller.Admin.AuthorizeUser)
		admin.GET("/authorize-user", controller.Admin.AuthorizeUser)
	}

	//设置路由分组
	setting := engine.Group("/setting")
	{
		//前端获取public配置
		setting.GET("",  controller.Setting.Index)
		//获取登录互联配置
		setting.GET("/login", controller.Setting.Login)
		//获取登录互联配置
		setting.GET("/connects", controller.Setting.Connects)
		//编辑页获取配置
		setting.GET("/site", controller.Setting.Site)
		//编辑页保存配置
		setting.POST("/site", controller.Setting.Save)
		//获取轮播图配置
		setting.GET("/carousel", controller.Setting.Carousel)
		//保存轮播图配置
		setting.POST("/carousel", controller.Setting.Save)
		//获取支付宝互联配置
		setting.GET("/alipay", controller.Setting.Alipay)
		//保存支付宝互联配置
		setting.POST("/alipay", controller.Setting.Alipay)
		//获取百度互联配置
		setting.GET("/baidu", controller.Setting.Baidu)
		//保存百度互联配置
		setting.POST("/baidu", controller.Setting.Baidu)
		//获取facebook互联配置
		setting.GET("/facebook", controller.Setting.Facebook)
		//保存facebook互联配置
		setting.POST("/facebook", controller.Setting.Facebook)
		//获取github互联配置
		setting.GET("/github", controller.Setting.Github)
		//保存github互联配置
		setting.POST("/github", controller.Setting.Github)
		//获取码云互联配置
		setting.GET("/gitee", controller.Setting.Gitee)
		//保存码云互联配置
		setting.POST("/gitee", controller.Setting.Gitee)
		//获取google互联配置
		setting.GET("/google", controller.Setting.Google)
		//保存google互联陪孩子
		setting.POST("/google", controller.Setting.Google)
		//获取line互联配置
		setting.GET("/line", controller.Setting.Line)
		//保存lime互联配置
		setting.POST("/line", controller.Setting.Line)
		//获取qq互联配置
		setting.GET("/qq", controller.Setting.Qq)
		//保存qq互联配置
		setting.POST("/qq", controller.Setting.Qq)
		//获取twitter互联配置
		setting.GET("/twitter", controller.Setting.Twitter)
		//保存twitter互联配置
		setting.POST("/twitter", controller.Setting.Twitter)
		//获取微信互联配置
		setting.GET("/wechat", controller.Setting.Wechat)
		//保存微信互联配置
		setting.POST("/wechat", controller.Setting.Wechat)
		//获取微博互联配置
		setting.GET("/weibo", controller.Setting.Weibo)
		//保存微博互联配置
		setting.POST("/weibo", controller.Setting.Weibo)
	}

	//首页轮播图路由分组
	carousel := engine.Group("/carousel")
	{
		//前端获取轮播图列表
		carousel.GET("", controller.Carousel.Index)
		//管理页获取轮播图列表
		carousel.GET("/list", controller.Carousel.List)
		//创建轮播图
		carousel.POST("/create", controller.Carousel.Create)
		//删除轮播图
		carousel.POST("/delete", controller.Carousel.Delete)
	}

	//稿件路由分组
	post := engine.Group("/post")
	{
		//前端获取上架稿件列表
		post.GET("", controller.Post.Index)
		//管理页获取稿件列表
		post.GET("/list", controller.Post.List)
		//保存稿件
		post.POST("/save", controller.Post.Save)
		//删除稿件
		post.POST("/delete", controller.Post.Delete)
		//修改状态
		post.POST("/toggle-status", controller.Post.ToggleStatus)
		//获取详情
		post.GET("/detail", controller.Post.Detail)
	}

	//文件路由分组
	file := engine.Group("/file")
	{
		//管理页获取文件列表
		file.GET("", controller.File.Index)
		//上传文件
		file.POST("/upload", controller.File.Upload)
		//删除文件
		file.POST("/delete", controller.File.Delete)
	}
	 */
}
