package gojwt

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Conf struct {
	Secret string `yaml:"secret"` // 签名密钥
	Header string `yaml:"header"` // 传token时http head的名称
	TTL    uint   `yaml:"ttl"`    // 有效时间，单位：秒
	Issuer string `yaml:"issuer"` // 颁发者，一般指系统名
}

// JwtToken jwt的token
type JwtToken struct {
	UID           string `json:"uid"`           // 所有者ID
	Role          string `json:"role"`          // 角色
	Token         string `json:"token"`         // token
	EffectiveTime uint   `json:"effectiveTime"` // 有效时间，单位：秒
}

func (t JwtToken) JSON() string {
	str, _ := json.Marshal(t)
	return string(str)
}

// RoleClaims 带角色的claims
type RoleClaims struct {
	Role string
	jwt.StandardClaims
}

// CreateToken 创建jwt token
func CreateToken(uid, role string, cfg Conf) (*JwtToken, error) {
	now := time.Now()
	expireTime := now.Add(time.Duration(cfg.TTL) * time.Second)

	claims := RoleClaims{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    cfg.Issuer,
			IssuedAt:  now.Unix(),
			Subject:   uid,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(cfg.Secret))
	if err != nil {
		return nil, err
	}
	return &JwtToken{
		UID:           uid,
		Role:          role,
		Token:         token,
		EffectiveTime: cfg.TTL,
	}, nil
}

// ParseToken 解析token
func ParseToken(token, secret string) (*RoleClaims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &RoleClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*RoleClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
