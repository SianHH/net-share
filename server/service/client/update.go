package client

import (
	"errors"
	"go.uber.org/zap"
	"net-share/server/global"
	"time"
)

type UpdateRequest struct {
	Code string `binding:"required" json:"code"`
	Name string `binding:"required" json:"name"`
}

type UpdateResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (*service) Update(params UpdateRequest) (UpdateResponse, error) {
	client, err := global.ClientFs.Query(params.Code)
	if err != nil {
		return UpdateResponse{}, err
	}
	client.Name = params.Name
	client.UpdatedAt = time.Now()
	if err := global.ClientFs.Update(client); err != nil {
		global.Logger.Error("修改client失败", zap.Error(err))
		return UpdateResponse{}, errors.New("修改失败")
	}
	return UpdateResponse{
		Code: client.Code,
		Name: client.Name,
		Key:  client.Key,
	}, nil
}
