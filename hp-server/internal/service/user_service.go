package service

import (
	"hp-server/internal/bean"
	"hp-server/internal/config"
	"hp-server/internal/db"
	"hp-server/internal/entity"
	"strings"
)

type UserService struct {
}

func (receiver *UserService) Login(login bean.ReqLogin) *bean.ResLoginUser {
	if strings.Compare(config.ConfigData.Admin.Username, login.Email) == 0 && strings.Compare(config.ConfigData.Admin.Password, login.Password) == 0 {
		return bean.NewAdminUser(login)
	} else {
		userQuery := entity.UserCustomEntity{}
		db.DB.Where("username = ?  and password = ?", login.Email, login.Password).First(&userQuery)
		if userQuery.Id != nil {
			return bean.NewClientUser(*userQuery.Id, userQuery.Username)
		}
	}
	return nil
}
