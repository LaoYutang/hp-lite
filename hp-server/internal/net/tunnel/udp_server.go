package tunnel

import (
	"hp-server/pkg/logger"
	"net"
	"strconv"
	"sync"

	"github.com/quic-go/quic-go"
)

type UdpServer struct {
	cache   sync.Map
	conn    quic.Connection
	udpConn *net.UDPConn
}

func NewUdpServer(conn quic.Connection) *UdpServer {
	return &UdpServer{
		sync.Map{},
		conn,
		nil,
	}
}

// ConnectLocal 内网服务的TCP链接
func (udpServer *UdpServer) StartServer(port int) bool {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		logger.Errorf("无法解析UDP地址 %d: %v", port, err)
		return false
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Errorf("无法创建UDP服务器 %d: %v", port, err)
		return false
	}
	udpServer.udpConn = conn
	//设置读
	go func() {
		// 创建缓冲区用于接收数据
		buffer := make([]byte, 1450)
		for {
			if udpServer.conn == nil {
				break
			}
			n, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				udpServer.cache.Range(func(key, value any) bool {
					handler := value.(*UdpHandler)
					go handler.ChannelInactive(conn)
					udpServer.cache.Delete(key)
					return true
				})
				break
			}
			value, ok := udpServer.cache.Load(addr.String())
			if !ok {
				handler := NewUdpHandler(udpServer, conn, udpServer.conn, addr)
				go handler.ChannelActive(conn)
				udpServer.cache.Store(addr.String(), handler)
			} else {
				handler := value.(*UdpHandler)
				go handler.ChannelRead(conn, buffer[:n])
			}
		}

		udpServer.cache.Range(func(key, value any) bool {
			handler := value.(*UdpHandler)
			go handler.ChannelInactive(conn)
			udpServer.cache.Delete(key)
			return true
		})

	}()
	return true
}

func (udpServer *UdpServer) CLose() {
	if udpServer.udpConn != nil {
		udpServer.udpConn.Close()
		udpServer.udpConn = nil
	}
}
