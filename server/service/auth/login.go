package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net-share/server/constant"
	"net-share/server/global"
	"time"
)

type LoginRequest struct {
	Account  string `binding:"required" json:"account"`
	Password string `binding:"required" json:"password"`
	Key      string `json:"key"`
	Captcha  string `json:"captcha"`
}

type LoginResponse struct {
	Token      string `json:"token,omitempty"`
	TokenExpAt string `json:"tokenExpAt,omitempty"`
}

func (*service) Login(c *gin.Context, param LoginRequest) (result LoginResponse, err error) {
	// 获取安全IP缓存
	security := global.Cache.GetString(constant.CacheSecurityIPKey+c.ClientIP()) == ""
	if !security {
		// 判断验证码
		code := global.Cache.GetString(constant.CacheLoginCaptChaKey + param.Key)
		if code == "" {
			return result, errors.New("验证码错误或已失效")
		}
		if code != param.Captcha {
			return result, errors.New("验证码错误或已失效")
		}
		global.Cache.Delete(constant.CacheLoginCaptChaKey + param.Key)
	}

	// 拉黑IP鉴权10分钟
	global.Cache.SetString(constant.CacheSecurityIPKey+c.ClientIP(), "1", time.Minute*10)

	if global.App.Account != param.Account || global.App.Password != param.Password {
		return result, errors.New("账号或密码错误")
	}
	tokenExpTime := time.Hour * 24 * 7
	token, err := global.JwtTool.GenerateToken(global.JwtTool.NewClaims(global.App.Account, map[string]string{}, tokenExpTime))
	if err != nil {
		return LoginResponse{}, err
	}
	result = LoginResponse{
		Token:      token,
		TokenExpAt: time.Now().Add(tokenExpTime).Format(time.DateTime),
	}
	// 登录成功，取消黑IP
	global.Cache.Delete(constant.CacheSecurityIPKey + c.ClientIP())
	return result, nil
}
