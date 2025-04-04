package server

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"hp-server/internal/net/base"
	"hp-server/internal/protol"
	"hp-server/pkg/logger"
	"math/big"
	"strconv"
	"time"

	"github.com/quic-go/quic-go"
)

type HpServer struct {
	base.QuicHandler
	listener *quic.Listener
}

func NewHPServer(handler base.QuicHandler) *HpServer {
	return &HpServer{
		handler,
		nil,
	}
}

func (quicServer *HpServer) generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"HP_LITE"},
	}
}

// ConnectLocal 内网服务的TCP链接
func (quicServer *HpServer) StartServer(port int) {
	q := &quic.Config{
		//最大空闲时间，超过就重连
		MaxIdleTimeout:        time.Duration(20) * time.Second,
		MaxIncomingStreams:    1000000,
		MaxIncomingUniStreams: 1000000,
		//空闲时，应该发送心跳包
		KeepAlivePeriod: time.Duration(5) * time.Second,
	}
	listener, err := quic.ListenAddr(":"+strconv.Itoa(port), quicServer.generateTLSConfig(), q)
	if err != nil {
		logger.Errorf("无法创建QUIC服务器 %d: %v", port, err)
	}
	quicServer.listener = listener
	//设置读
	go func() {
		for {
			conn, err := listener.Accept(context.Background())
			if err != nil {
				logger.Errorf("QUIC获取连接错误：%v", err)
			}
			go func() {
				for {
					stream, err := conn.AcceptStream(context.Background())
					if err != nil {
						go quicServer.ChannelInactive(stream, conn)
						logger.Errorf("接收流错误：全部关闭: %v", err)
						return
					}
					// 为每个连接启动一个新的处理 goroutine
					quicServer.handler(stream, conn)
				}
			}()
		}
	}()
	logger.Infof("数据传输服务启动成功UDP: %d", port)
}

func (quicServer *HpServer) handler(stream quic.Stream, conn quic.Connection) {
	go func() {
		defer stream.Close()
		quicServer.ChannelActive(stream, conn)
		reader := bufio.NewReader(stream)
		//避坑点：多包问题，需要重复读取解包
		for {
			decode, e := protol.Decode(reader)
			if e != nil {
				quicServer.ChannelInactive(stream, conn)
				return
			}
			if decode != nil {
				e := quicServer.ChannelRead(stream, decode, conn)
				if e != nil {
					return
				}
			}
		}
	}()
}

func (quicServer *HpServer) CLose() {
	if quicServer.listener != nil {
		quicServer.listener.Close()
		quicServer.listener = nil
	}
}
