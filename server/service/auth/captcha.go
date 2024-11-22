package auth

import (
	"encoding/base64"
	"errors"
	"go.uber.org/zap"
	"net-share/pkg/captcha"
	"net-share/pkg/utils"
	"net-share/server/constant"
	"net-share/server/global"
	"time"
)

func (*service) Captcha() (key, bs64 string, err error) {
	code := utils.RandStr(4, utils.NumDict)
	bytes, err := captcha.Generate(120, 40, code)
	if err != nil {
		global.Logger.Error("生成图片验证码失败", zap.Error(err))
		return "", "", errors.New("获取失败")
	}
	bs64 = base64.StdEncoding.EncodeToString(bytes)
	key = utils.RandStr(16, utils.AllDict)
	global.Cache.Set(constant.CacheLoginCaptChaKey+key, code, time.Minute*5)
	return
}
