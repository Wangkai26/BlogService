package app

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/pkg/util"
	"time"
)
// 处理JWT令牌
type Claims struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	// StandarClaims中的字段都是非强制性的但官方建议使用的预定义权利要求
	// 能够提供一组有用的，可互操作的约定
	jwt.StandardClaims
}
// 从全局变量获取秘钥
func GetJWTSecret() []byte {
	return []byte(global.JWTSetting.Secret)
}


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

// 到此，完成了 JWT 令牌的生成、解析 和 校验 方法的编写
// 现在我们需要在后续的应用中间件中对其进行调用，使其能够在应用程序中将一整套的动作给串起来