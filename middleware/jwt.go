package middleware

import (
	"errors"
	"fmt"
	"second_hand_mall/internal/global"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT 用户基础结构
type JwtBase struct {
	UserId   int64  // 公有系统用户ID -- user表对用户ID
	UserName string // 用户名字
	Nick     string // 用户昵称
	Phone    string // 手机号
	Email    string // 邮箱
}

// 已登录的token信息
type Logined struct {
	NowToken string //存储当前登录的token
	ExpTime  int64  //存储当前登录的token过期时间
}

// 存储已登录
var tokenMap = sync.Map{}

func storeToken(key string, value Logined) {
	tokenMap.Store(key, value)
}

func loadToken(key string) (Logined, bool) {
	value, ok := tokenMap.Load(key)
	if ok {
		return value.(Logined), ok
	}
	return Logined{}, ok
}

func deleteToken(key string) {
	tokenMap.Delete(key)
}

// 登录
type SaleLogin struct {
	JwtBase       // 用户基本信息
	CId     int64 // 公司ID
	jwt.StandardClaims
}

// CreateToken 生成一个token
func (this *SaleLogin) CreateToken() *jwt.Token {
	this.StandardClaims.NotBefore = time.Now().Unix() - 30                                                                 // 签名生效时间
	this.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(global.GVAL_CONFIG.JWT.ExpiresTime) * time.Second).Unix() // 过期时间
	this.StandardClaims.Issuer = global.GVAL_CONFIG.JWT.Issuer                                                             // 签名的发行者
	return jwt.NewWithClaims(jwt.SigningMethodHS256, *this)
}

// 解析Token
func (this *SaleLogin) ParseToken(c *gin.Context, s_token string) (string, error) {
	Token, err := jwt.ParseWithClaims(s_token, &SaleLogin{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.GVAL_CONFIG.JWT.SigningKey), nil
	})
	if err != nil {
		fmt.Println("秘钥", global.GVAL_CONFIG.JWT.SigningKey)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return "", errors.New("非法Token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", errors.New("登录过期,请重新登录")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return "", errors.New("令牌尚未激活")
			} else {
				return "", errors.New("无法处理此令牌")
			}
		}
	}
	claims, ok := Token.Claims.(*SaleLogin)
	if ok && Token.Valid {
		// 单会话登录 检查
		if t, ok := loadToken(claims.UserName); ok && t.NowToken != s_token && t.ExpTime > claims.StandardClaims.ExpiresAt {
			return "", errors.New("此账号已在其他地点登录,请重新登录!")
		}
		storeToken(claims.UserName, Logined{s_token, claims.StandardClaims.ExpiresAt})

		this.UserId = claims.UserId
		this.UserName = claims.UserName
		this.Nick = claims.Nick
		this.Phone = claims.Phone
		this.Email = claims.Email
		this.CId = claims.CId
		this.StandardClaims.NotBefore = claims.StandardClaims.NotBefore // 签名生效时间
		this.StandardClaims.ExpiresAt = claims.StandardClaims.ExpiresAt // 过期时间 一天
		this.StandardClaims.Issuer = claims.StandardClaims.Issuer       // 签发的发行者

		return "", nil
	}
	return "", errors.New("无法处理此令牌")
}

// 登录
func JwtSaleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		str_token := c.Request.Header.Get("Authorization")
		if str_token == "" {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}
		var t SaleLogin
		token, err := t.ParseToken(c, str_token)
		if nil != err {
			c.JSON(200, gin.H{
				"code": -2,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}
		fmt.Println("=================================token.CId", t.CId)
		c.Set("c_id", t.CId)
		c.Set("x-token", token)
	}
}

// 登录鉴权
func JwtSaleAdminLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		str_token := c.Request.Header.Get("Authorization")
		if str_token == "" {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}
		var t SaleLogin
		_, err := t.ParseToken(c, str_token)
		if nil != err {
			c.JSON(200, gin.H{
				"code": -2,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}
		if 0 != t.CId {
			c.JSON(200, gin.H{
				"code": -3,
				"msg":  "此账号无权限",
			})
			c.Abort()
			return
		}
		c.Set("c_id", t.CId)
	}
}

// 解析token
func ParseSaleLogin(c *gin.Context) SaleLogin {
	str_token := c.Request.Header.Get("Authorization")
	fmt.Println("Login TOKEN：", str_token)
	var t SaleLogin
	t.ParseToken(c, str_token)
	return t
}
