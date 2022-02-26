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

#### 数据库

新建数据库并导入`sql`目录下的`pms.sql`和`patch.sql`文件

#### 配置文件
`config/app/yaml`

- 站点配置

对应配置文件中的`site`节点。

| 参数名      | 类型     | 说明        |
|:---------|:-------|:----------|
| name     | string | 网站名称      |
| maintain | bool   | 是否开启维护    |
| debug    | bool   | 是否开启debug |
| listen   | int    | 服务监听端口    |
| cors     | bool   | 是否允许跨域    |

- 数据库配置

对应配置文件中的`database`节点。

| 参数名      | 类型     | 说明   |
|:---------|:-------|:-----|
| host     | string | 服务地址 |
| username | string | 账号   |
| password | string | 密码   |
| database | string | 数据库  |
| charset  | string | 编码   |
| prefix   | string | 表前缀  |

- Redis配置

对应配置文件中的`redis`节点。

| 参数名  | 类型     | 说明   |
|:-----|:-------|:-----|
| host | string | 服务地址 |
| port | int    | 端口   |
| auth | string | 密码   |

- 缓存配置

对应配置文件中的`cache`节点。

| 参数名    | 类型     | 说明   |
|:-------|:-------|:-----|
| prefix | string | 键名前缀 |
| expire | int    | 缓存时长 |

- 会话配置

对应配置文件中的`session`节点。

| 参数名    | 类型     | 说明       |
|:-------|:-------|:---------|
| name   | string | cookie键名 |
| prefix | string | 缓存键名前缀   |
| secret | string | 加密密钥     |
| idle   | int    | 会话空闲时长   |


- 日志配置

日志使用`uber`的`zap`，对应配置文件中的`log`节点。

| 参数名   | 类型     | 说明        |
|:------|:-------|:----------|
| file  | string | 文件名       |
| level | string | 日志等级      |
| json  | string | 是否保存为json |

- 网页配置

对应配置文件中的`web`节点。

| 参数名                             | 类型     | 说明                                                                                                     |
|:--------------------------------|:-------|:-------------------------------------------------------------------------------------------------------|
| security                        | object | 安全配置                                                                                                   |
| security.primaryKey             | string | 加密私钥地址，与前端的`publicKey`配置对应                                                                             |
| security.salt                   | string | 数据加密加盐                                                                                                 |
| host                            | string | 服务访问域名                                                                                                 |
| connects                        | array  | 第三方授权支持类型，`alipay`、`github`、`baidu`、`gitee`、`facebook`、`google`、`line`、`qq`、`twitter`、`wechat`、`weibo` |
| upload                          | object | 文件上传配置                                                                                                 |
| upload.path                     | string | 文件存储目录                                                                                                 |
| upload.url                      | string | 文件存储目录访问地址                                                                                             |
| upload.video                    | object | 视频文件上传配置                                                                                               |
| upload.video.ffmpeg             | string | ffmpeg bin文件地址                                                                                         |
| upload.video.ffprobe            | string | ffprobe bin文件地址                                                                                        |
| upload.video.extensions         | array  | 支持文件后缀                                                                                                 |
| upload.video.maxSize            | string | 最大单个文件大小                                                                                               |
| upload.video.maxFiles           | string | 最大文件数量                                                                                                 |
| upload.video.mimeTypes          | array  | 支持文件类型                                                                                                 |
| upload.image                    | object | 图片文件上传配置                                                                                               |
| upload.image.extensions         | array  | 支持文件后缀                                                                                                 |
| upload.image.maxSize            | string | 最大单个文件大小                                                                                               |
| upload.image.maxFiles           | string | 最大文件数量                                                                                                 |
| upload.image.mimeTypes          | array  | 支持文件类型                                                                                                 |

### 构建
构建为Linux 64位可执行程序，生成文件在工程根目录

```shell
GOOS=linux GOARCH=amd64 go build .
```



