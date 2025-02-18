package controller

import (
	"encoding/json"
	"hp-server/internal/bean"
	"hp-server/internal/service"
	"hp-server/pkg/util"
	"net/http"
	"strconv"
)

type MonitorController struct {
	service.MonitorService
}

func (receiver MonitorController) getUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	token := r.Header.Get("token")
	userId, _, _, err := util.DecodeToken(token)
	if err != nil {
		json.NewEncoder(w).Encode(bean.ResErrorCode(-2, "用户权限校验失败"))
		return 0, err
	}
	return userId, nil
}

func (receiver MonitorController) List(w http.ResponseWriter, r *http.Request) {
	id, err := receiver.getUserId(w, r)
	if err == nil {
		data := receiver.ListData(id)
		if data != nil {
			json.NewEncoder(w).Encode(bean.ResOk(data))
			return
		} else {
			json.NewEncoder(w).Encode(bean.ResError("登陆失败"))
		}
	}
}

func (receiver MonitorController) GetMonitorData(w http.ResponseWriter, r *http.Request) {
	configId, err := strconv.Atoi(r.URL.Query().Get("configId"))
	if err != nil {
		json.NewEncoder(w).Encode(bean.ResErrorCode(-2, "配置id错误"))
		return
	}
	data := receiver.GetMonitorDataByConfigId(configId)
	json.NewEncoder(w).Encode(bean.ResOk(data))
}
