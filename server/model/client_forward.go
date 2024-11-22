package model

import (
	"fmt"
	"strings"
)

type ClientForward struct {
	Base
	Name        string `json:"name"`
	Target      string `json:"target"`
	Port        string `json:"port"`
	Nodelay     int    `json:"nodelay"`
	ClientCode  string `json:"clientCode"`
	RateLimiter int    `json:"rateLimiter"`
}

func (ClientForward) TableName() string {
	return "gost_client_forward"
}

func (l ClientForward) GetLimits() []string {
	unit := "KB"
	return []string{
		fmt.Sprintf("$ %d%s %d%s", l.RateLimiter, unit, l.RateLimiter, unit),
	}
}

func (l ClientForward) GetTargetIpAndPort() (ip, port string) {
	split := strings.Split(l.Target, ":")
	if len(split) == 2 {
		return split[0], split[1]
	}
	return "", ""
}
