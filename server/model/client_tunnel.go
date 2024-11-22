package model

import (
	"fmt"
	"strings"
)

type ClientTunnel struct {
	Base
	Name         string `json:"name"`
	Target       string `json:"target"`
	DomainPrefix string `json:"domainPrefix"`
	ClientCode   string `json:"clientCode"`
	RateLimiter  int    `json:"rateLimiter"`
	Key          string `json:"key"`
	AuthUser     string `json:"authUser"`
	AuthPwd      string `json:"authPwd"`
}

func (ClientTunnel) TableName() string {
	return "gost_client_tunnel"
}

func (l ClientTunnel) GetLimits() []string {
	unit := "KB"
	return []string{
		fmt.Sprintf("$ %d%s %d%s", l.RateLimiter, unit, l.RateLimiter, unit),
	}
}

func (l ClientTunnel) GetTargetIpAndPort() (ip, port string) {
	split := strings.Split(l.Target, ":")
	if len(split) == 2 {
		return split[0], split[1]
	}
	return "", ""
}
