package model

type Client struct {
	Base
	Name string `json:"name"`
	Key  string `json:"key"`
}
