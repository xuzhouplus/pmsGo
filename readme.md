## 基于Gin的pms服务端

### 目录结构

```
controller #控制器
model      #数据模型
service    #服务
config         #yaml配置文件
lib
	app        #程序初始化
	cache      #缓存初始化
	controller #基础控制器
	database   #数据库初始化
	file       #文件处理
	helper     #杂项
	log        #zap-uber初始化
	middleware #中间件
	model      #基础模型
	oauth      #账号登录认证
	security   #随机数、加解密等
public         #静态资源
router         #路由
sql            #初始sql
```
### 环境搭建

#### 开发工具
 - GoLand 2021.3.3
 - go1.17.7
 - MySQL5.7

#### GoLand配置
在设置中`Go Modules`项配置`envrionment`值`GOPROXY=https://mirrors.aliyun.com/goproxy/`

### 构建
构建为Linux 64位可执行程序，生成文件在工程根目录
```
GOOS=linux GOARCH=amd64 go build .
```



