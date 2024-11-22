package ws_msg_data

import (
	"encoding/json"
	"fmt"
)

type MessageBody struct {
	Callback string `json:"callback"`
	UUID     string `json:"uuid,omitempty"`
	Aes      bool   `json:"aes"`
	Data     string `json:"data,omitempty"`
}

func NewMessageBody(callback, uuid string, aes bool, data any) MessageBody {
	var result = MessageBody{
		Callback: callback,
		UUID:     uuid,
		Aes:      aes,
		Data:     "",
	}
	_ = result.Marshal(data)
	return result
}

func (m *MessageBody) Marshal(data any) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if m.Aes {
		m.Data = string(AesEncryptECB(marshal, aesKey))
	} else {
		m.Data = string(marshal)
	}
	return nil
}

func (m *MessageBody) Unmarshal(data any) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if m.Aes {
		decrypted := AesDecryptECB([]byte(m.Data), aesKey)
		return json.Unmarshal(decrypted, data)
	}
	return json.Unmarshal([]byte(m.Data), data)
}
