package server

import (
	"errors"
	hpMessage "hp-server/internal/message"
	"hp-server/internal/service"
	"hp-server/pkg/logger"

	"github.com/quic-go/quic-go"
)

type HPClientHandler struct {
	*service.HpService
}

func NewHPHandler() *HPClientHandler {
	return &HPClientHandler{
		&service.HpService{},
	}
}

// ChannelActive 连接激活时，发送注册信息给云端
func (h *HPClientHandler) ChannelActive(stream quic.Stream, conn quic.Connection) {
	logger.Infof("HP传输打开流：%d", stream.StreamID())
}

func (h *HPClientHandler) ChannelRead(stream quic.Stream, data interface{}, conn quic.Connection) error {
	message := data.(*hpMessage.HpMessage)
	if message == nil {
		logger.Errorf("消息类型:解码异常|ip:%s", conn.RemoteAddr().String())
		stream.Close()
		return errors.New("HP消息类型:解码异常")
	}
	logger.Infof("流ID:%d|HP消息类型:%s|IP:%s", stream.StreamID(), message.Type.String(), conn.RemoteAddr())
	switch message.Type {
	case hpMessage.HpMessage_REGISTER:
		{
			h.Register(message, conn)
		}
	}
	return nil
}

func (h *HPClientHandler) ChannelInactive(stream quic.Stream, conn quic.Connection) {
	if stream != nil {
		logger.Infof("HP传输断开流：%d", stream.StreamID())
	}
}
