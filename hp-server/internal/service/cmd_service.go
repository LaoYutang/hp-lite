package service

import (
	"encoding/json"
	"hp-server/internal/bean"
	"hp-server/internal/db"
	"hp-server/internal/entity"
	"hp-server/internal/message"
	"hp-server/internal/protol"
	"hp-server/pkg/logger"
	"net"
	"sync"
)

// string,conn
var CMD_CACHE_CONN = sync.Map{}

// string memory_info
var CMD_CACHE_MEMORY_INFO = sync.Map{}

type CmdService struct {
}

func (receiver *CmdService) sendMessage(conn net.Conn, message *message.CmdMessage) {
	conn.Write(protol.CmdEncode(message))
}

func (receiver *CmdService) sendTips(conn net.Conn, tips string) {
	c := &message.CmdMessage{
		Data: tips,
		Type: message.CmdMessage_TIPS,
	}
	receiver.sendMessage(conn, c)
}

func SendCloseMsg(deviceKey, msg string) bool {
	value, ok := CMD_CACHE_CONN.Load(deviceKey)
	if ok {
		c := &message.CmdMessage{
			Data: msg,
			Type: message.CmdMessage_DISCONNECT,
		}
		conn := value.(net.Conn)
		conn.Write(protol.CmdEncode(c))
		conn.Close()
		CMD_CACHE_CONN.Delete(deviceKey)
	}
	return ok
}

func NoticeClientUpdateData(deviceKey string) bool {
	var results []entity.UserConfigEntity
	db.DB.Model(&entity.UserConfigEntity{}).Where("device_key = ?", deviceKey).Find(&results)
	if results != nil {
		var results2 []bean.LocalInnerWear
		for _, item := range results {
			wear := bean.LocalInnerWear{
				OutLimit:    -1,
				InLimit:     -1,
				ConnectType: *item.ConnectType,
				ConfigKey:   item.ConfigKey,
				LocalIp:     item.LocalIp,
				LocalPort:   *item.LocalPort,
				ServerIp:    item.ServerIp,
				ServerPort:  *item.ServerPort,
			}
			results2 = append(results2, wear)
		}
		jsonData, err := json.Marshal(results2)
		if err == nil {
			c := &message.CmdMessage{
				Data: string(jsonData),
				Type: message.CmdMessage_LOCAL_INNER_WEAR,
			}
			value, ok := CMD_CACHE_CONN.Load(deviceKey)
			if ok {
				conn := value.(net.Conn)
				conn.Write(protol.CmdEncode(c))
			}
			return ok
		}

	}
	return false
}

func (receiver CmdService) StoreMemInfo(conn net.Conn, message *message.CmdMessage) {
	data := message.GetData()
	if len(data) > 0 {
		info := &bean.MemoryInfo{}
		err := json.Unmarshal([]byte(data), info)
		if err == nil {
			CMD_CACHE_MEMORY_INFO.Delete(message.GetKey())
			CMD_CACHE_MEMORY_INFO.Store(message.GetKey(), info)
			CMD_CACHE_CONN.Delete(message.GetKey())
			CMD_CACHE_CONN.Store(message.GetKey(), conn)
		}
	}
}

func (receiver *CmdService) Connect(conn net.Conn, message *message.CmdMessage) {
	_, ok := CMD_CACHE_CONN.Load(message.GetKey())
	if ok {
		logger.Infof("设备KEY已经在线: %s", message.GetKey())
		receiver.sendTips(conn, "设备KEY已经在线")
		return
	} else {
		receiver.StoreMemInfo(conn, message)
		NoticeClientUpdateData(message.GetKey())
	}
}

func (receiver CmdService) Clear(conn net.Conn) {
	deviceKey := ""
	CMD_CACHE_CONN.Range(func(key, value interface{}) bool {
		deviceKey = key.(string)
		return value != conn
	})
	//清除数据
	if len(deviceKey) > 0 {
		logger.Infof("清除设备key: %s", deviceKey)
		CMD_CACHE_CONN.Delete(deviceKey)
		CMD_CACHE_MEMORY_INFO.Delete(deviceKey)
	}
}
