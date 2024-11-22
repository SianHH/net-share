# 介绍
基于GOST开发的内网穿透系统，包含最基础的http、tcp、udp转发，不内置https转发，推荐使用Nginx代理添加多域名SSL，灵活性更高

[GOST开源地址](https://github.com/go-gost/gost)

# 优点
- 在网页新增修改隧道配置
- 速率限制
- 传输效率高

# 快速开始

1. 下载二进制文件 
> 略过

2. 运行服务端
> 第一次运行会以默认配置文件运行，默认配置文件生成在运行目录的config.yaml

配置文件说明：
```yaml
# 基础系统配置
# Web服务的地址
addr: 0.0.0.0:8080
# 运行模式，一般不需要修改
mode: prod
# Web登录的账号
account: admin
# Web登录的密码
password: "123456"
# jwt随机密钥
jwt-key: 1hvsestbk5v38nvq0yqcgi7ev6t0hif6
# 日志，一般默认即可
logger:
    file: application.log
    level: 0
    console: true

# 服务器IP，用于告知客户端连接的IP
ip: 127.0.0.1

# 域名解析配置，缺少任意一项，域名解析功能都无法使用
# 使用http转发时，使用的基础域名
domain: example.com
# http转发时，流量进入的端口，可反向代理此端口，给http添加上ssl证书
entrypoint: "18080"
# 转发http隧道时，客户端会和此端口交换数据
host-port: "2096"

# 端口转发配置，缺少任意一项，端口转发功能都无法使用
# 转发tcp/udp时，客户端会和此端口交换数据
forward-port: "2097"
# 端口转发时，可以分配的端口
ports:
    - 10001-11000
    - "20000"
    - "30000"
```

3. 运行客户端

> 需要预先在后台管理添加一个客户端，用于获取客户端的连接密钥

客户端启动命令：
```shell
./client -s ws://127.0.0.1:8080 -key xxxxxx 
```
参数说明：
- `-s`配置服务端地址，服务端向客户端分发配置使用的是websocket连接，如果web管理服务暴露的是https服务，需要将`ws://`修改为`wss://`
- `-key`客户端的连接密钥，这里`xxxxx`是示例

# 演示网页
> 此项目并非演示网站源码，但是由演示网站功能精简而来，去除了多用户扩展出来的相关功能，两则的使用方式差别不大

[https://gost.sian.one](https://gost.sian.one)

账号：guest@guest.com

密码：guest
