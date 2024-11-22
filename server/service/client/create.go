package client

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net-share/server/global"
	"net-share/server/model"
	"time"
)

type CreateRequest struct {
	Name string `binding:"required" json:"name"`
}

type CreateResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (*service) Create(params CreateRequest) (CreateResponse, error) {
	var client = model.Client{
		Base: model.Base{
			Code:      uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: params.Name,
		Key:  uuid.NewString(),
	}

	if err := global.ClientFs.Create(client); err != nil {
		global.Logger.Error("新增client失败", zap.Error(err))
		return CreateResponse{}, errors.New("新增失败")
	}
	return CreateResponse{
		Code: client.Code,
		Name: client.Name,
		Key:  client.Key,
	}, nil
}
