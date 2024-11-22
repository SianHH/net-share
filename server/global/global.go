package global

import (
	"net-share/pkg/cache"
	"net-share/pkg/file_storage"
	"net-share/pkg/jwt"
	"net-share/pkg/logger"
	"net-share/server/config"
	"net-share/server/model"
	"sync"
)

var App config.App

var Logger *logger.Logger

var ClientFs *file_storage.FileStorage[model.Client]
var ClientForwardFs *file_storage.FileStorage[model.ClientForward]
var ClientHostFs *file_storage.FileStorage[model.ClientHost]
var ClientTunnelFs *file_storage.FileStorage[model.ClientTunnel]
var ClientHostDomainFs *file_storage.FileStorage[model.Base]

var Cache *cache.Cache

var JwtTool *jwt.Tool

var WsMap = &sync.Map{}
