package utils

import (
	"gopayroll/common"
	"gopayroll/db"
)

type Switch struct {
	HhrCfgSwCode string
	HhrCfgSwVal  string
}

//获取开关值
func GetSwitchMap(tenantId int64, values ...interface{}) (map[string]string, error) {
	result := make(map[string]string)
	engine := db.OrmEngine("hhr_foundation")
	rsSlice := make([]*Switch, 0)
	err := engine.Table("hhr_switch_cfg").Where("tenant_id=?", tenantId).In("hhr_cfg_sw_code", values).
		Cols("hhr_cfg_sw_code", "hhr_cfg_sw_val").Find(&rsSlice)
	if err != nil {
		common.Logger.Error("[GetSwitchMap]", err.Error())
	}
	for i := 0; i < len(rsSlice); i++ {
		result[rsSlice[i].HhrCfgSwCode] = rsSlice[i].HhrCfgSwVal
	}
	return result, err
}
