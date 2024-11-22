package model

import (
	"fmt"
	"strings"
)

type ClientHost struct {
	Base
	Name         string `json:"name"`
	Target       string `json:"target"`
	DomainPrefix string `json:"domainPrefix"`
	ClientCode   string `json:"clientCode"`
	RateLimiter  int    `json:"rateLimiter"`
}

func (ClientHost) TableName() string {
	return "gost_client_host"
}

func (l ClientHost) GetLimits() []string {
	unit := "KB"
	return []string{
		fmt.Sprintf("$ %d%s %d%s", l.RateLimiter, unit, l.RateLimiter, unit),
	}
}

func (l ClientHost) GetTargetIpAndPort() (ip, port string) {
	split := strings.Split(l.Target, ":")
	if len(split) == 2 {
		return split[0], split[1]
	}
	return "", ""
}
