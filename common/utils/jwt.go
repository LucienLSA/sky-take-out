package utils

import (
	"errors"
	"fmt"
	"skytakeout/global"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomPayload 自定义载荷继承原有接口并附带自己的字段
type CustomPayload struct {
	UserId   uint64
	UserName string
	jwt.RegisteredClaims
}

// GenerateToken 生成Token uid 用户id subject 签发对象  secret 加盐
func GenerateTokenV0(uid uint64, uname string, secret string) (string, error) {
	claim := CustomPayload{
		UserId:   uid,
		UserName: uname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   //签发者
			Subject:   uname,                                           //签发对象
			Audience:  jwt.ClaimStrings{"PC", "Wechat_Program"},        //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
	return token, err
}

func GenerateTokenV1(uid uint64, uname string, secret string) (accessToken, refreshToken string, err error) {
	nowTime := time.Now()
	ttl := global.Config.Jwt.Admin.TTL
	expireTime := nowTime.Add(time.Duration(ttl) * time.Minute)
	rtExpireTime := nowTime.Add(time.Duration(ttl) * time.Hour)
	claim := CustomPayload{
		UserId:   uid,
		UserName: uname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   //签发者
			Subject:   uname,                                           //签发对象
			Audience:  jwt.ClaimStrings{"PC", "Wechat_Program"},        //签发受众
			ExpiresAt: jwt.NewNumericDate(expireTime),                  //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	claim2 := CustomPayload{
		UserId:   uid,
		UserName: uname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   //签发者
			Subject:   uname,                                           //签发对象
			Audience:  jwt.ClaimStrings{"PC", "Wechat_Program"},        //签发受众
			ExpiresAt: jwt.NewNumericDate(rtExpireTime),                //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim2).SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// 解析token
func ParseToken(token string, secret string) (myclaims *CustomPayload, err error) {
	myclaims = &CustomPayload{}
	tokenClaims, err := jwt.ParseWithClaims(token, myclaims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secret), nil
	})
	fmt.Println(tokenClaims)
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*CustomPayload); ok && tokenClaims.Valid { // 校验token
			return claims, nil
		}
	}
	return nil, errors.New("invalid token")
}

// GenerateAccessToken 只生成 access_token
func GenerateAccessToken(uid uint64, uname string, secret string) (accessToken string, err error) {
	nowTime := time.Now()
	ttl := global.Config.Jwt.Admin.TTL
	expireTime := nowTime.Add(time.Duration(ttl) * time.Minute)
	claim := CustomPayload{
		UserId:   uid,
		UserName: uname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   //签发者
			Subject:   uname,                                           //签发对象
			Audience:  jwt.ClaimStrings{"PC", "Wechat_Program"},        //签发受众
			ExpiresAt: jwt.NewNumericDate(expireTime),                  //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
	return
}

// ParseRefreshToken 验证用户token
func ParseRefreshToken(aToken, rToken, secret string) (newAToken, newRToken string, err error) {
	accessClaim, err := ParseToken(aToken, secret)
	if err != nil {
		return
	}

	refreshClaim, err := ParseToken(rToken, secret)
	if err != nil {
		return
	}

	// 阈值：5分钟
	const refreshThreshold = 5 * time.Minute

	// access_token 没过期
	if accessClaim.ExpiresAt.After(time.Now()) {
		remaining := time.Until(accessClaim.ExpiresAt.Time)
		if remaining < refreshThreshold {
			// 只刷新 access_token
			newAToken, err = GenerateAccessToken(accessClaim.UserId, accessClaim.UserName, secret)
			if err != nil {
				return
			}
			return newAToken, rToken, nil
		}
		return aToken, rToken, nil
	}

	// access_token 过期，但 refresh_token 没过期，可以刷新一对新的 token
	if refreshClaim.ExpiresAt.After(time.Now()) {
		return GenerateTokenV1(accessClaim.UserId, accessClaim.UserName, secret)
	}

	// 两个都过期，强制重新登录
	return "", "", errors.New("身份过期，重新登陆")
}
