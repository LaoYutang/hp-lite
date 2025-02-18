package server

import (
	"bufio"
	"hp-server/pkg/logger"
	"net"
	"strconv"
)

type CmdServer struct {
	listener net.Listener
}

func NewCmdServer() *CmdServer {
	return &CmdServer{
		nil,
	}
}

// ConnectLocal 内网服务的TCP链接
func (tcpServer *CmdServer) StartServer(port int) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		logger.Errorf("不能创建TCP服务器：" + ":" + strconv.Itoa(port) + " 原因：" + err.Error() + " 提示：" + err.Error())
		return
	}
	tcpServer.listener = listener
	//设置读
	go func() {
		for {
			if tcpServer.listener == nil {
				return
			}
			conn, err := listener.Accept()
			if err == nil {
				tcpServer.handler(conn)
			} else {
				logger.Errorf("TCP错误连接: ", err)
			}
		}
	}()
	logger.Infof("指令传输服务启动成功TCP: %d", port)

}

func (tcpServer *CmdServer) handler(conn net.Conn) {
	go func() {
		defer conn.Close()
		handler := NewCmdHandler()
		handler.ChannelActive(conn)
		reader := bufio.NewReader(conn)
		for {
			if tcpServer.listener == nil {
				return
			}
			//尝试读检查连接激活
			_, err := reader.Peek(1)
			if err != nil {
				handler.ChannelInactive(conn)
				return
			}

			decode, e := handler.Decode(reader)
			if e != nil {
				logger.Errorf("解码异常:%v", e)
				handler.ChannelInactive(conn)
				return
			}
			if conn != nil {
				err := handler.ChannelRead(conn, decode)
				if err != nil {
					return
				}
			} else {
				return
			}
		}
	}()
}

func (tcpServer *CmdServer) CLose() {
	if tcpServer.listener != nil {
		tcpServer.listener.Close()
		tcpServer.listener = nil
	}
}
