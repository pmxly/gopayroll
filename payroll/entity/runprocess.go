/*
Desc: 运行过程实体
Author: 潘承勋
Date: 2020-05-03
*/

package entity

import (
	"gopayroll/common"
	"gopayroll/db"
	pyCommon "gopayroll/payroll/common"
	"gopayroll/payroll/pyformula"
	"gopayroll/payroll/pysysutils"
	"time"
)

type RunProcess struct {
	TenantId     int64
	RunProcessId string          `xorm:"hhr_runprocess_id"`
	Country      string          `xorm:"hhr_country"`
	LineData     []*ProcessChild `xorm:"-"`
	Calendar     *Calendar       `xorm:"-"`
}

type ProcessChild struct {
	TenantId      int64
	RunProcessId  string      `xorm:"hhr_runprocess_id"`
	SeqNum        uint32      `xorm:"hhr_seqnum"`
	ElementType   string      `xorm:"hhr_element_type"`
	ElementId     string      `xorm:"hhr_element_id"`
	Country       string      `xorm:"-"`
	StartDate     time.Time   `xorm:"hhr_start_dt"`
	EndDate       time.Time   `xorm:"hhr_end_dt"`
	ElementEntity interface{} `xorm:"-"`
}

func NewRunProcess(tenantId int64, runProcessId string, calendar *Calendar) *RunProcess {
	runP := &RunProcess{TenantId: tenantId, RunProcessId: runProcessId, Calendar: calendar}
	engine := db.OrmEngine("hhr_payroll")
	_, err := engine.Table("hhr_py_runprocess").Where("tenant_id=?", tenantId).Get(runP)
	if err != nil {
		common.Logger.Error("[NewRunProcess-1]", err.Error())
	}

	runP.LineData = make([]*ProcessChild, 0)
	err = engine.Table("hhr_py_runprocess_child").Where("tenant_id=? and hhr_runprocess_id=?", tenantId, runProcessId).
		Asc("hhr_seqnum").Find(&runP.LineData)
	if err != nil {
		common.Logger.Error("[NewRunProcess-2]", err.Error())
	}
	for _, processChild := range runP.LineData {
		if processChild.EndDate.IsZero() {
			processChild.EndDate, _ = time.Parse(common.TimeLayout, "9999/12/31")
		}
		if processChild.Country == "ALL" {
			globalKey := pysysutils.GetGlobalKey(tenantId, calendar.CalGrpId)
			varCache := pyCommon.VarCache(globalKey)
			processChild.Country = varCache.GetVarValue("VR_F_COUNTRY").(string)
		}
	}
	common.Logger.Debug("------NewRunProcess--------", runP)
	return runP
}

//todo
func (runChild *ProcessChild) ElementExec() {
	switch runChild.ElementType {
	case "FM":
		runChild.ElementEntity.(*pyformula.Formula).Exec()
	case "FC":
	case "WT":
	case "WC":
	}
}
