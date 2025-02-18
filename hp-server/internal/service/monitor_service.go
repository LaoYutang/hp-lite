package service

import (
	"hp-server/internal/db"
	"hp-server/internal/entity"
)

type MonitorService struct {
}

func (receiver MonitorService) ListData(userId int) map[int][]entity.UserStatisticsEntity {
	var results2 []entity.UserStatisticsEntity
	if userId > 0 {
		var configIds []int
		var results []entity.UserConfigEntity
		db.DB.Model(entity.UserConfigEntity{}).Where("user_id = ?", userId).Find(&results)
		for _, result := range results {
			configIds = append(configIds, *result.Id)
		}
		db.DB.Model(&entity.UserStatisticsEntity{}).Where("config_id in ?", configIds).Find(&results2)
	} else {
		db.DB.Model(&entity.UserStatisticsEntity{}).Find(&results2)
	}

	result := make(map[int][]entity.UserStatisticsEntity)

	for _, user := range results2 {
		// 使用 configId 作为 key，将 UserStatisticsEntity 添加到对应的列表中
		result[user.ConfigId] = append(result[user.ConfigId], user)
	}
	return result
}

// 通过配置id获取监控数据
func (receiver MonitorService) GetMonitorDataByConfigId(configId int) []entity.UserStatisticsEntity {
	var results []entity.UserStatisticsEntity
	tx := db.DB.Model(&entity.UserStatisticsEntity{}).Where("config_id = ?", configId).Order("time").Find(&results)
	if tx.Error != nil {
		return []entity.UserStatisticsEntity{}
	}

	return results
}
