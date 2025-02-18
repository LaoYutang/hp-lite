package controller

import (
	"encoding/json"
	"hp-server/internal/bean"
	"hp-server/internal/service"
	"hp-server/pkg/logger"
	"net/http"
)

type LoginController struct {
	service.UserService
}

func (receiver LoginController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var msg bean.ReqLogin
	// 解析请求体中的JSON数据
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		logger.Errorf("解析请求体中的JSON数据失败: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	login := receiver.Login(msg)
	if login != nil {
		json.NewEncoder(w).Encode(bean.ResOk(login))
		return
	} else {
		json.NewEncoder(w).Encode(bean.ResError("登陆失败"))
	}
}
