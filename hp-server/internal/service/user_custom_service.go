package service

import (
	"hp-server/internal/bean"
	"hp-server/internal/db"
	"hp-server/internal/entity"
	"time"
)

type UserCustomService struct {
}

func (receiver *UserCustomService) AddData(custom entity.UserCustomEntity) {
	if custom.Id == nil {
		custom.CreateTime = time.Now()
	}
	db.DB.Save(&custom)
}

func (receiver *UserCustomService) ListData(page int, pageSize int) *bean.ResPage {
	var results []entity.UserCustomEntity
	var total int64
	// 计算总记录数并执行分页查询
	db.DB.Model(&entity.UserCustomEntity{}).Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Find(&results)
	return bean.PageOk(total, results)
}

func (receiver *UserCustomService) RemoveData(id int) {
	db.DB.Delete(&entity.UserCustomEntity{Id: &id})
}
