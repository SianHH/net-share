package main

import (
	"crypto/tls"
	"net-share/client/ws"

	// Register connectors
	_ "github.com/go-gost/x/connector/direct"
	_ "github.com/go-gost/x/connector/forward"
	_ "github.com/go-gost/x/connector/http"
	_ "github.com/go-gost/x/connector/http2"
	_ "github.com/go-gost/x/connector/relay"
	_ "github.com/go-gost/x/connector/serial"
	_ "github.com/go-gost/x/connector/sni"
	_ "github.com/go-gost/x/connector/socks/v4"
	_ "github.com/go-gost/x/connector/socks/v5"
	_ "github.com/go-gost/x/connector/ss"
	_ "github.com/go-gost/x/connector/ss/udp"
	_ "github.com/go-gost/x/connector/sshd"
	_ "github.com/go-gost/x/connector/tcp"
	_ "github.com/go-gost/x/connector/tunnel"
	_ "github.com/go-gost/x/connector/unix"
	"github.com/lxzan/gws"
	"net/http"

	// Register dialers
	_ "github.com/go-gost/x/dialer/direct"
	_ "github.com/go-gost/x/dialer/dtls"
	_ "github.com/go-gost/x/dialer/ftcp"
	_ "github.com/go-gost/x/dialer/grpc"
	_ "github.com/go-gost/x/dialer/http2"
	_ "github.com/go-gost/x/dialer/http2/h2"
	_ "github.com/go-gost/x/dialer/http3"
	_ "github.com/go-gost/x/dialer/http3/wt"
	_ "github.com/go-gost/x/dialer/icmp"
	_ "github.com/go-gost/x/dialer/kcp"
	_ "github.com/go-gost/x/dialer/mtcp"
	_ "github.com/go-gost/x/dialer/mtls"
	_ "github.com/go-gost/x/dialer/mws"
	_ "github.com/go-gost/x/dialer/obfs/http"
	_ "github.com/go-gost/x/dialer/obfs/tls"
	_ "github.com/go-gost/x/dialer/pht"
	_ "github.com/go-gost/x/dialer/quic"
	_ "github.com/go-gost/x/dialer/serial"
	_ "github.com/go-gost/x/dialer/ssh"
	_ "github.com/go-gost/x/dialer/sshd"
	_ "github.com/go-gost/x/dialer/tcp"
	_ "github.com/go-gost/x/dialer/tls"
	_ "github.com/go-gost/x/dialer/udp"
	_ "github.com/go-gost/x/dialer/unix"
	_ "github.com/go-gost/x/dialer/ws"

	// Register handlers
	_ "github.com/go-gost/x/handler/auto"
	_ "github.com/go-gost/x/handler/dns"
	_ "github.com/go-gost/x/handler/file"
	_ "github.com/go-gost/x/handler/forward/local"
	_ "github.com/go-gost/x/handler/forward/remote"
	_ "github.com/go-gost/x/handler/http"
	_ "github.com/go-gost/x/handler/http2"
	_ "github.com/go-gost/x/handler/http3"
	_ "github.com/go-gost/x/handler/metrics"
	_ "github.com/go-gost/x/handler/redirect/tcp"
	_ "github.com/go-gost/x/handler/redirect/udp"
	_ "github.com/go-gost/x/handler/relay"
	_ "github.com/go-gost/x/handler/serial"
	_ "github.com/go-gost/x/handler/sni"
	_ "github.com/go-gost/x/handler/socks/v4"
	_ "github.com/go-gost/x/handler/socks/v5"
	_ "github.com/go-gost/x/handler/ss"
	_ "github.com/go-gost/x/handler/ss/udp"
	_ "github.com/go-gost/x/handler/sshd"
	_ "github.com/go-gost/x/handler/tap"
	_ "github.com/go-gost/x/handler/tun"
	_ "github.com/go-gost/x/handler/tunnel"
	_ "github.com/go-gost/x/handler/unix"

	// Register listeners
	_ "github.com/go-gost/x/listener/dns"
	_ "github.com/go-gost/x/listener/dtls"
	_ "github.com/go-gost/x/listener/ftcp"
	_ "github.com/go-gost/x/listener/grpc"
	_ "github.com/go-gost/x/listener/http2"
	_ "github.com/go-gost/x/listener/http2/h2"
	_ "github.com/go-gost/x/listener/http3"
	_ "github.com/go-gost/x/listener/http3/h3"
	_ "github.com/go-gost/x/listener/http3/wt"
	_ "github.com/go-gost/x/listener/icmp"
	_ "github.com/go-gost/x/listener/kcp"
	_ "github.com/go-gost/x/listener/mtcp"
	_ "github.com/go-gost/x/listener/mtls"
	_ "github.com/go-gost/x/listener/mws"
	_ "github.com/go-gost/x/listener/obfs/http"
	_ "github.com/go-gost/x/listener/obfs/tls"
	_ "github.com/go-gost/x/listener/pht"
	_ "github.com/go-gost/x/listener/quic"
	_ "github.com/go-gost/x/listener/redirect/tcp"
	_ "github.com/go-gost/x/listener/redirect/udp"
	_ "github.com/go-gost/x/listener/rtcp"
	_ "github.com/go-gost/x/listener/rudp"
	_ "github.com/go-gost/x/listener/serial"
	_ "github.com/go-gost/x/listener/ssh"
	_ "github.com/go-gost/x/listener/sshd"
	_ "github.com/go-gost/x/listener/tap"
	_ "github.com/go-gost/x/listener/tcp"
	_ "github.com/go-gost/x/listener/tls"
	_ "github.com/go-gost/x/listener/tun"
	_ "github.com/go-gost/x/listener/udp"
	_ "github.com/go-gost/x/listener/unix"
	_ "github.com/go-gost/x/listener/ws"
)

import (
	"flag"
	"fmt"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/x/config/parsing"
	"log"
	"os"
	"time"

	xlogger "github.com/go-gost/x/logger"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	logger.SetDefault(xlogger.NewLogger(xlogger.LevelOption(logger.FatalLevel)))
	parsing.BuildDefaultTLSConfig(nil)
}

func main() {
	var server string
	var key string
	flag.StringVar(&server, "s", "", "websocket url,example: wss://gost.sian.one")
	flag.StringVar(&key, "key", "", "client key")
	flag.Parse()
	if server == "" {
		fmt.Println("pleas enter the websocket url")
		os.Exit(1)
	}
	if key == "" {
		fmt.Println("pleas enter the client key")
		os.Exit(1)
	}
	fmt.Println("WebSocketUrl:", server+"/api/v1/client/ws")
	fmt.Println("ClientKey:", key)
	for {
		socket, _, err := gws.NewClient(ws.NewService("ws", key), &gws.ClientOption{
			Addr:          server + "/api/v1/client/ws",
			TlsConfig:     &tls.Config{InsecureSkipVerify: true},
			RequestHeader: http.Header{"key": []string{key}},
			PermessageDeflate: gws.PermessageDeflate{
				Enabled:               true,
				ServerContextTakeover: true,
				ClientContextTakeover: true,
			},
		})
		if err != nil {
			fmt.Println("conn fail,please wait 5 second,retry conn", err)
			time.Sleep(time.Second * 5)
			continue
		}
		_ = socket.WritePing(nil)
		socket.ReadLoop()
	}
}
