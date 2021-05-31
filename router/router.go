package router

import (
	"github.com/gin-gonic/gin"
	"pmsGo/app/controller"
	"pmsGo/middleware/auth"
	"pmsGo/middleware/session"
)

func Router(engine *gin.Engine) {
	//支持上传文件的支持
	engine.Static("/public", "./public")
	//注册session组件
	engine.Use(session.Register())
	//账号路由分组
	admin := engine.Group("/admin")
	{
		//账号登录
		admin.POST("/login", auth.Register(), controller.Admin.Login)
		//获取登录用户信息
		admin.GET("/auth", auth.Register(), controller.Admin.Auth)
		//登出
		admin.POST("/logout", auth.Register(), controller.Admin.Logout)
		//获取账号信息
		admin.GET("/profile", auth.Register(), controller.Admin.Profile)
		//保存账号信息
		admin.POST("/profile", auth.Register(), controller.Admin.Profile)
		//获取第三方互联账号
		admin.GET("/connects", auth.Register(), controller.Admin.Connects)
		//获取第三方互联授权地址
		admin.GET("/authorize", auth.Register(), controller.Admin.AuthorizeUrl)
		admin.GET("/authorize-url", auth.Register(), controller.Admin.AuthorizeUrl)
		//获取第三方互联授权账号
		admin.GET("/callback", auth.Register(), controller.Admin.AuthorizeUser)
		admin.GET("/authorize-user", auth.Register(), controller.Admin.AuthorizeUser)
	}

	//设置路由分组
	setting := engine.Group("/setting")
	{
		//前端获取public配置
		setting.GET("", controller.Setting.Index)
		//获取登录互联配置
		setting.GET("/login", controller.Setting.Login)
		//获取登录互联配置
		setting.GET("/connects", auth.Register(), controller.Setting.Connects)
		//编辑页获取配置
		setting.GET("/site", auth.Register(), controller.Setting.Site)
		//编辑页保存配置
		setting.POST("/site", auth.Register(), controller.Setting.Save)
		//获取轮播图配置
		setting.GET("/carousel", auth.Register(), controller.Setting.Carousel)
		//保存轮播图配置
		setting.POST("/carousel", auth.Register(), controller.Setting.Save)
		//获取支付宝互联配置
		setting.GET("/alipay", auth.Register(), controller.Setting.Alipay)
		//保存支付宝互联配置
		setting.POST("/alipay", auth.Register(), controller.Setting.Alipay)
		//获取百度互联配置
		setting.GET("/baidu", auth.Register(), controller.Setting.Baidu)
		//保存百度互联配置
		setting.POST("/baidu", auth.Register(), controller.Setting.Baidu)
		//获取facebook互联配置
		setting.GET("/facebook", auth.Register(), controller.Setting.Facebook)
		//保存facebook互联配置
		setting.POST("/facebook", auth.Register(), controller.Setting.Facebook)
		//获取github互联配置
		setting.GET("/github", auth.Register(), controller.Setting.Github)
		//保存github互联配置
		setting.POST("/github", auth.Register(), controller.Setting.Github)
		//获取码云互联配置
		setting.GET("/gitee", auth.Register(), controller.Setting.Gitee)
		//保存码云互联配置
		setting.POST("/gitee", auth.Register(), controller.Setting.Gitee)
		//获取google互联配置
		setting.GET("/google", auth.Register(), controller.Setting.Google)
		//保存google互联陪孩子
		setting.POST("/google", auth.Register(), controller.Setting.Google)
		//获取line互联配置
		setting.GET("/line", auth.Register(), controller.Setting.Line)
		//保存lime互联配置
		setting.POST("/line", auth.Register(), controller.Setting.Line)
		//获取qq互联配置
		setting.GET("/qq", auth.Register(), controller.Setting.Qq)
		//保存qq互联配置
		setting.POST("/qq", auth.Register(), controller.Setting.Qq)
		//获取twitter互联配置
		setting.GET("/twitter", auth.Register(), controller.Setting.Twitter)
		//保存twitter互联配置
		setting.POST("/twitter", auth.Register(), controller.Setting.Twitter)
		//获取微信互联配置
		setting.GET("/wechat", auth.Register(), controller.Setting.Wechat)
		//保存微信互联配置
		setting.POST("/wechat", auth.Register(), controller.Setting.Wechat)
		//获取微博互联配置
		setting.GET("/weibo", auth.Register(), controller.Setting.Weibo)
		//保存微博互联配置
		setting.POST("/weibo", auth.Register(), controller.Setting.Weibo)
	}

	//首页轮播图路由分组
	carousel := engine.Group("/carousel")
	{
		//前端获取轮播图列表
		carousel.GET("", controller.Carousel.Index)
		//管理页获取轮播图列表
		carousel.GET("/list", auth.Register(), controller.Carousel.List)
		//创建轮播图
		carousel.POST("/create", auth.Register(), controller.Carousel.Create)
		//删除轮播图
		carousel.POST("/delete", auth.Register(), controller.Carousel.Delete)
	}

	//稿件路由分组
	post := engine.Group("/post")
	{
		//前端获取上架稿件列表
		post.GET("", controller.Post.Index)
		//管理页获取稿件列表
		post.GET("/list", auth.Register(), controller.Post.List)
		//保存稿件
		post.POST("/save", auth.Register(), controller.Post.Save)
		//删除稿件
		post.POST("/delete", auth.Register(), controller.Post.Delete)
		//修改状态
		post.POST("/toggle-status", auth.Register(), controller.Post.ToggleStatus)
		//获取详情
		post.GET("/detail", controller.Post.Detail)
	}

	//文件路由分组
	file := engine.Group("/file")
	{
		//管理页获取文件列表
		file.GET("", auth.Register(), controller.File.Index)
		//上传文件
		file.POST("/upload", auth.Register(), controller.File.Upload)
		//删除文件
		file.POST("/delete", auth.Register(), controller.File.Delete)
	}

}
