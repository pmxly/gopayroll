/*
Desc: 运行类型实体
Author: 潘承勋
Date: 2020-05-03
*/

package entity

import (
	"gopayroll/common"
	"gopayroll/db"
)

type RunType struct {
	TenantId  int64
	RunTypeId string          `xorm:"hhr_runtype_id"`
	Country   string          `xorm:"hhr_country"`
	Cycle     string          `xorm:"hhr_cycle"`
	LineData  []*RunTypeChild `xorm:"-"`
	Calendar  *Calendar       `xorm:"-"`
}

type RunTypeChild struct {
	TenantId         int64
	RunTypeId        string      `xorm:"hhr_runtype_id"`
	SeqNum           uint32      `xorm:"hhr_seqnum"`
	RunProcessId     string      `xorm:"hhr_runprocess_id"`
	Active           string      `xorm:"hhr_active"`
	RunProcessEntity *RunProcess `xorm:"-"`
}

func NewRunType(tenantId int64, runTypeId string, calendar *Calendar) *RunType {
	runType := &RunType{TenantId: tenantId, RunTypeId: runTypeId, Calendar: calendar}
	engine := db.OrmEngine("hhr_payroll")
	_, err := engine.Table("hhr_py_runtype").Where("tenant_id=?", tenantId).Get(runType)
	if err != nil {
		common.Logger.Error("[NewRunType-1]", err.Error())
	}

	runType.LineData = make([]*RunTypeChild, 0)
	err = engine.Table("hhr_py_runtype_child").Where("tenant_id=? and hhr_runtype_id=? and hhr_active='Y' ", tenantId, runTypeId).
		Asc("hhr_seqnum").Find(&runType.LineData)
	if err != nil {
		common.Logger.Error("[NewRunType-2]", err.Error())
	}
	for _, runTypeChild := range runType.LineData {
		runTypeChild.RunProcessEntity = NewRunProcess(tenantId, runTypeChild.RunProcessId, calendar)
	}
	common.Logger.Debug("------NewRunType--------", runType)
	return runType
}
