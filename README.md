# HTTP应用-博客后端

此项目源自《Go语言编程之旅》第二章内容，作者源码：[go-programming-tour-book/blog-service: 《Go 语言编程之旅：一起用 Go 做项目》第二章：博客程序（HTTP Server） (github.com)](https://github.com/go-programming-tour-book/blog-service)

## 一、使用 gin web框架

- gin 初始化

```go
func main()  {
   gin.SetMode(global.ServerSetting.RunMode)
   // 下面三行打印用来校验配置是否真正地映射到了配置结构体上
   //log.Println(global.AppSetting)
   //log.Println(global.ServerSetting)
   //log.Println(global.DatabaseSetting)
   router := routers.NewRouter()
   s := &http.Server{
      Addr:        ":"+global.ServerSetting.HttpPort,
      Handler:      router,
      ReadTimeout:    10 * time.Second,
      WriteTimeout:   10 * time.Second,
      MaxHeaderBytes: 1<<20,
   }
   global.Logger.Infof("%s:go-programming-tour-book/%s","eddycjy","blog-service")
   s.ListenAndServe()
}
```

- 配置中间件及路由

```go
func NewRouter() *gin.Engine {
   r := gin.New()
   r.Use(gin.Logger())
   r.Use(gin.Recovery())
   r.Use(middleware.Translations())
   r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
   article := v1.NewArticle()
   tag := v1.NewTag()

   // 第七节添加代码，上传图片路由
   upload := api.NewUpload()
   r.POST("/upload/file",upload.UploadFile)
   r.StaticFS("/static",http.Dir(global.AppSetting.UploadSavePath))
   // 第八节新增下一行代码
   r.POST("/auth",api.GetAuth)
   apiv1 := r.Group("/api/v1")
   apiv1.Use(middleware.JWT())
   {
      apiv1.POST("/tags",tag.Create)
      apiv1.DELETE("/tags/:id",tag.Delete)
      apiv1.PUT("/tags/:id",tag.Update)
      apiv1.PATCH("/tags/:id/state",tag.Update)
      apiv1.GET("/tags/",tag.List)

      apiv1.POST("/articles",article.Create)
      apiv1.DELETE("/articles/:id",article.Delete)
      apiv1.PUT("/articles/:id",article.Update)
      apiv1.PATCH("/articles/:id/state",article.Update)
      apiv1.GET("/articles/:id",article.Get)
      apiv1.GET("/articles",article.List)
   }

   // part8.接入JWT中间件，利用gin中分组路由的概念
   // 只对apiv1的路由分组进行JWT中间件的引入，也就是说只有apiv1路由分组里的路由方法
   // 会受到此中间件的约束

   return r
}
```

中间件写法：

```go
// CustomRecoveryWithWriter returns a middleware for a given writer that recovers from any panics and calls the provided handle func to handle it.
func CustomRecoveryWithWriter(out io.Writer, handle RecoveryFunc) HandlerFunc {
   var logger *log.Logger
   if out != nil {
      logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
   }
   return func(c *Context) {
      defer func() {
          ...
```

与之前设计的 Gee 框架一样，有 Next 函数

## 二、项目设计

无设计不开发，先对本次需求的内容进行设计和评审。本节对项目的目录结构、接口方案、路由注册等内容进行设计和开发。

### 目录结构

- configs 配置文件，内有文件 config.yaml，使用 viper 读取配置
- docs 文档集合，swagger 生成文档
- global 存放全局变量，eg：DBEngine、设置中的全局变量 ServerSetting、AppSetting、DatabaseSetting、Logger、JWTSetting，主要就是对应配置文件，这些变量主体是 setting 中的结构体，结构体需要与配置文件对应
- internal 内部模块

  - dao 数据访问层（Database Access Object），所有与数据相关的操作都会在 dao 层进行，如 MySQL
  	- auth.go
  	- dao.go 处理标签模块的 dao 操作
  		dao 层进行了数据访问对象的封装，并对业务所需的字段进行了处理
  	- tag.go
  - middleware：HTTP 中间件
  	- jwt.go
  	- translations.go 辅助jwt
  - model：模型层，用于存放 model 对象
  	- article.go 创建文章 Model，与数据库字段一一对应
  	- article_tag.go 创建文章标签 Model
  	- auth.go
  	- model.go 编写公共字段结构体 Model，根据 DSN，生成 DBEngine
  	- tag.go 创建 Tag Model，也就是 Tag 结构体
  - routers：路由相关的逻辑
  	- router.go 注册路由，路由管理
  	- api 文件夹，是一个分组，内部新建文件夹 v1，v1 中新建 tag.go 和 article.go，这里编写对应路由的处理方法，写好处理方法后，将其注册到对应的路由即可
  		- article.go
  		- tag.go
  	
  - service：项目核心业务逻辑，专注业务逻辑，对于其中需要的数据库操作，都用 dao 实现
  	- auth.go 认证信息，AuthRequest 结构体用于接口入参的校验，AppKey 和 AppSecret 都为必选项，CheckAuth 验证客户端传入的认证信息是否存在，若不存在返回错误信息
  	- service.go 声明 Service 结构体，包含 context 和 dao.Dao，还有一个 New 函数，之后 tag 对标签进行增删查改都需要用到 Service 结构体
  	- tag.go  包含对标签进行 增删查改的方法，都有对应的结构体，标签中 form 和 binding 分别代表着表单的映射字段名和入参校验的规则内容，CountTagRequest 我不记得是做什么的了
  		在接口校验一节中，添加 binding 标签
  	- upload.go
- pkg：项目相关的模块包

  - app
  	- form.go 对入参校验的方法进行二次封装，在 BindAndValid 方法中，通过 ShouldBind 进行参数绑定和入参校验。发生错误后，通过中间件 Translations 中设置的 Translator 对错误消息体进行具体的翻译
  - convert
  - errcode
    - common_code.go 预定义项目中的一些公共错误码
    - errcode.go 编写常用的一些错误处理公共方法，标准化错误输出
    - module_code.go 针对标签模块的错误码
  - logger 日志标准化处理
    - logger.go 对日志分等级，然后编写具体的方法，对日志的实例初始化和标准化参数进行绑定
  - setting：对读取配置的行为进行封装
    - section.go 声明配置属性的结构体，编写读取区段配置的配置方法
    - setting.go NewSetting函数用于初始化本项目配置的基础属性，即设定配置文件的名称为 config、配置类型为 yaml，并且设置其配置路径为相对路径 configs/，以确保在该项目目录下能够成功启动编写组件
  - upload
    - file.go
  - util
    - md5.go
- storage：项目生成的临时文件
	- logs
		- app.log

	- uploads

- scripts：各类构建、安装、分析等操作的脚本
- third_party：第三方的资源工具，如 Swagger UI
- main.go 

## 二、viper 读取配置文件

本项目的配置管理使用最常见的文件配置作为我们的选型

configs 目录新建 config.yaml，写入配置

- Server：服务配置，设置 gin 的运行模式、默认的 HTTP 监听端口、允许读取和写入的最大持续时间
- APP：应用配置，设置默认每页数量、所允许的最大每页数量，以及默认的应用日志存储路径
- Database：数据库配置，主要是连接实例所必需的基础参数

pkg/setting/section.go 中 存放对应的结构体，ReadSection 是读取区段配置的方法

初始化的代码片段NewSetting 函数用于读取配置

```go
func setupSetting() error {
   setting,err := setting.NewSetting()
   if err != nil{
      return err
   }
   err = setting.ReadSection("Server",&global.ServerSetting)
   if err != nil{
      return err
   }
   err = setting.ReadSection("App",&global.AppSetting)
   if err != nil{
      return err
   }
   err = setting.ReadSection("Database",&global.DatabaseSetting)
   if err != nil{
      return err
   }
   // 8section,设置JWT的一些相关配置
   err = setting.ReadSection("JWT",&global.JWTSetting)
   if err != nil{
      return err
   }

   global.ServerSetting.ReadTimeout *= time.Second
   global.ServerSetting.WriteTimeout *= time.Second
   global.JWTSetting.Expire *= time.Second
   return nil
}
```

调用方式，main.go 中的 init() 中，setupSetting（）

## 三、gorm

引入 gorm 开源库的同时需要引入 MySQL驱动库，若不引入驱动，会有报错

使用 MySQL，初始化代码

后面主要就是 回调处理公共字段 和 CRUD

使用ORM 的目的在于完全屏蔽数据库细节，解放生产力。

初始化在 main.go 中的 setupDBEngine

```go
func setupDBEngine() error {
	var err error
	// 这里一定不能用 :=
	// 由于 := 会重新声明并创建左侧的新局部变量，会导致等式右方变量没有赋值给全局变量 global.DBEngine
	// 其他包调用 global.DBEngine 时，它仍然是 nil
	global.DBEngine,err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil{
		return err
	}

	return nil
}
```

## 四、lumberjack 写日志

该开源库的核心功能是将日志写入滚动文件中，允许我们设置单日志文件的最大占用空间、最大生命周期等

pkg 目录下新建 logger 目录，创建 logger.go，写入日志分级相关代码

logger.go：写入日志分级相关代码，预定义了应用日志的 Level 和 Fields 的具体类型，把日志分为 debug、info、warn、error、fatal 和 panic 六个等级，以便在不同的使用场景中记录不同级别的日志。

完成日志分级后，编写具体的方法，用来对日志的实例初始化和标准化参数进行绑定

- WithLevel：设置日至等级
- WithFields：设置日志公共字段
- WithContext：设置日志上下文属性
- WithCaller：设置当前某一层调用栈的信息（程序计数器、文件信息和行号）
- WithCallersFrames：设置当前的整个调用栈信息

后面是编写日志内容的格式化和日志输出方法

最后是日志输出等级

总结：

- 日志分级
- 日志标准化
- 日志格式化和输出
- 日志分级输出

**设置全局变量**，在 globa 文件夹下的 setting.go 中写入 Logger 全局变量 

**在 mian.go 中进行初始化**

```go
func setupLogger() error {
   filename := global.AppSetting.LogSavePath+"/"+global.AppSetting.LogFileName+global.AppSetting.LogFileExt
   global.Logger = logger.NewLogger(&lumberjack.Logger{
      Filename: filename,
      MaxSize:  600,
      MaxAge:   10,
      LocalTime: true,
   },"",log.LstdFlags).WithCaller(2)

   return nil
}
```

init 方法中新增了日志组件的流程，并在 setupLogger 方法内部对 global 的包全局变量 Logger 进行了初始化。

## 五、swagger 生成文档

描述一个 API 的基本信息：

- 有关该 API 的描述
- 可用路径（或资源）
- 在每个路径上的可用操作（获取和提交等）
- 每个操作的输入和输出格式

写入注解：

- @Summary：摘要
- @Produce：API可以产生的MIME类型的列表。我们可以把MIME类型简单地理解为相应类型，如 JSON、XML、HTML 等
- @Param：参数格式，从左到右分别为：参数名、入参类型、数据类型、是否必填和注释
- @Success：响应成功，从左到右分别为：状态码、参数类型、数据类型和注释
- @Failure：响应失败，从左到右分别为：状态码、参数类型、数据类型和注释
- @Router：路由，从左到右分别为：路由地址和 HTTP 方法

写完注释后，swag init 就能生成文档

在 routers 中进行默认初始化和注册对应路由，就能在对应路由访问该文档

## 六、validator 接口入参校验

gin 使用该中间件做参数验证，现在对其作一定的定制，并添加到中间件中

validator 可以对字段要求 必填、大于、大于等于、小于、小于等于、最大值、最小值、其中之一、长度要求与len一致。标签 form 和 binding 分别代表着表单的映射字段名 和 入参校验的规则内容。

绑定

在 pkg/app 目录中新建 form.go，用 gin.ShouldBind 进行参数绑定和入参校验

中间件 Translations，主要做国际化和验证器注册

具体使用：api/v1/tag.o

```go
func (t Tag) Update(c *gin.Context) {
   param := service.UpdateTagRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
   response := app.NewResponse(c)
   valid, errs := app.BindAndValid(c, &param)
   if !valid {
      global.Logger.Errorf("app.BindAndValid errs: %v", errs)
      response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
      return
   }

   svc := service.New(c.Request.Context())
   err := svc.UpdateTag(&param)
   if err != nil {
      global.Logger.Errorf("svc.UpdateTag err: %v", err)
      response.ToErrorResponse(errcode.ErrorUpdateTagFail)
      return
   }

   response.ToResponse(gin.H{})
   return
}
```

去煎鱼源码看看，t 是否真的是空结构体，我这是，

## 七、jwt - API访问控制

在 pkg/app/jwt.go：

一个函数用于生成 jwt

```go
//这一部分承担了整个流程中比较重要的职责，也就是生成JWT的行为
//函数流程逻辑是根据客户端传入的AppKey和AppSecret以及在项目配置中所设置的
//Issuer和ExpiresAt，根据指定的算法生成签名后的Token
func GenerateToken(appKey, appSecret string) (string, error) {
   nowTime := time.Now()
   expireTime := nowTime.Add(global.JWTSetting.Expire)
   claims := Claims{
      AppKey:    util.EncodeMD5(appKey),
      AppSecret: util.EncodeMD5(appSecret),
      StandardClaims: jwt.StandardClaims{
         ExpiresAt: expireTime.Unix(),
         Issuer:    global.JWTSetting.Issuer,
      },
   }
   // 根据Claims结构体创建Token实例
   tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
   // 根据传入的秘钥生成签名并返回标准的 Token
   token, err := tokenClaims.SignedString(GetJWTSecret())
   return token, err
}
```

一个函数用于解析 jwt

```go
// 解析和校验 Token，承担着与 Generate Token 相对的功能，
// 函数流程主要是解析传入的 Token，然后根据 Claims 的相关属性要求进行校验
func ParseToken(token string) (*Claims, error) {
   tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
      return GetJWTSecret(), nil
   })
   if err != nil {
      return nil, err
   }
   if tokenClaims != nil {
      // Valid：验证基于时间的声明，例如：ExpiresAt、Issuer、Not Before
      // 需要注意的是，在没有任何声明的令牌中，仍然认为是有效的
      if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
         return claims, nil
      }
   }

   return nil, err
}
```

应用 jwt 中间件

在 middleware/jwt.go 中写入相应代码

接入 jwt 中间件：apiv1 的路由分组进行 jwt 中间件的引入

```go
apiv1 := r.Group("/api/v1")
apiv1.Use(middleware.JWT())
```

一些别的处理

相应处理

## bug 修复

tag 的启用状态为未启用，也就是 0 时，无法更新成功。

global.DBEngine,err := model.NewDBEngine(global.DatabaseSetting) ，是错误的，由于 := 会重新声明并创建左侧的新局部变量，因此在其他包中调用 global.DBEngine 变量时，其值仍然是 nil

正确的做法是 var err error，然后 global.DBEngine,err = model.NewDBEngine(global.DatabaseSetting) 