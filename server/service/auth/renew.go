package auth

import (
	"github.com/gin-gonic/gin"
	"net-share/server/global"
	"net-share/server/router/middleware"
	"time"
)

type RenewResponse struct {
	Token      string `json:"token"`
	TokenExpAt string `json:"tokenExpAt"`
}

func (*service) Renew(c *gin.Context) (result RenewResponse, err error) {
	claims := middleware.GetClaims(c)
	if time.Now().Unix()+int64(time.Hour)*24 > claims.ExpiresAt {
		tokenExpTime := time.Hour * 24 * 7
		token, err := global.JwtTool.GenerateToken(global.JwtTool.NewClaims(claims.Id, map[string]string{}, tokenExpTime))
		if err != nil {
			return RenewResponse{}, err
		}
		result = RenewResponse{
			Token:      token,
			TokenExpAt: time.Now().Add(tokenExpTime).Format(time.DateTime),
		}
	} else {
		result = RenewResponse{
			Token:      middleware.GetToken(c),
			TokenExpAt: time.Unix(claims.ExpiresAt, 0).Format(time.DateTime),
		}
	}
	return result, nil
}
