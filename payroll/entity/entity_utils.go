package entity

import (
	"gopayroll/common"
	"gopayroll/db"
	pyCommon "gopayroll/payroll/common"
	"gopayroll/payroll/pysysutils"
	"github.com/go-xorm/xorm"
	"time"
)

func InitVariables(tenantId int64, calGrp *CalendarGroup) {
	country := calGrp.Country
	varObjLst := make([]*pyCommon.VariableObject, 0)
	engine := db.OrmEngine("hhr_payroll")
	err := engine.Table("hhr_py_variable").
		Where("tenant_id=? and hhr_status='Y' and (hhr_country='ALL' or hhr_country=?)", tenantId, country).Find(&varObjLst)
	if err != nil {
		common.Logger.Error("[InitVariables]", err.Error())
	}
	globalKey := pysysutils.GetGlobalKey(tenantId, calGrp.CalGrpId)
	varCache := pyCommon.VarCache(globalKey)
	for _, varObj := range varObjLst {
		varObj.Country = "ALL"
		varObj.HasCovered = "N"
		varCache.SetVarObj(varObj.VarId, varObj)
	}
}

/**
插入薪资标记&计算消息
@msgClass: A-标记，B-计算
@msgType: S-成功，E-失败，W-警告，I-通知
@msgText: 消息文本
*/
func InsPyCalcMsg(session *xorm.Session, calGrpId string, payee *Payee, cal *Calendar, msgClass string, msgType string, msgTxt string, userId string) error {
	tenantId := cal.TenantId
	calId := cal.CalId
	empId := payee.EmpId
	empRcd := payee.EmpRcd
	maxSeqSql := "select max(hhr_seq_num) from hhr_payroll.hhr_py_payee_calc_msg where tenant_id = ? and hhr_pycalgrp_id = ? " +
		"and hhr_py_cal_id = ? and hhr_empid = ? and hhr_emp_rcd = ? and hhr_py_msg_class = ? "
	var seqNum int64
	var newSeq int64
	_, err := session.SQL(maxSeqSql, tenantId, calGrpId, calId, empId, empRcd, msgClass).Get(&seqNum)
	if err != nil {
		common.Logger.Error("[beginIdentify-5]", err.Error())
		return err
	}
	if seqNum != 0 {
		newSeq = seqNum + 1
	} else {
		newSeq = 1
	}
	payeeCalcMsg := &PayeeCalcMsg{
		TenantId:   tenantId,
		CalGrpId:   calGrpId,
		CalId:      calId,
		EmpId:      empId,
		EmpRcd:     empRcd,
		MsgClass:   msgClass,
		SeqNum:     newSeq,
		FCalId:     "",
		MsgType:    msgType,
		MsgTxt:     msgTxt,
		CreateDtm:  common.CurLocalDate(),
		CreateUser: userId,
		ModifyDtm:  common.CurLocalDate(),
		ModifyUser: userId,
	}
	_, err = session.Table("hhr_py_payee_calc_msg").InsertOne(payeeCalcMsg)
	if err != nil {
		common.Logger.Error("[beginIdentify-6]", err.Error())
	}
	return nil
}

func InsPyCalStat() {

}

/**
校验受款人在期间结束日期内的薪资组是否与所在日历薪资组一致
*/
func CheckPayeePayGroup(session *xorm.Session, payee *Payee, cal *Calendar, bgnDate time.Time, endDate time.Time) (bool, error) {
	tenantId := cal.TenantId
	empId := payee.EmpId
	empRcd := payee.EmpRcd
	var prdEndDt time.Time
	if bgnDate.IsZero() && endDate.IsZero() {
		prdEndDt = cal.PeriodCalEntity.EndDate
	} else {
		prdEndDt = endDate
	}
	pyGrpId := cal.PyGroupId
	t1 := "select a.hhr_pygroup_id from hhr_payroll.hhr_py_assign_pg a where a.tenant_id = ? and a.hhr_efft_date <= ? " +
		"and a.hhr_empid = ? and hhr_emp_rcd = ? order by hhr_efft_date desc, hhr_efft_seq desc"
	var empPyGrpId string
	_, err := session.SQL(t1, tenantId, prdEndDt, empId, empRcd).Get(&empPyGrpId)
	if err != nil {
		common.Logger.Error("[CheckPayeePayGroup]", err.Error())
		return false, err
	}
	if empPyGrpId != pyGrpId {
		return false, nil
	} else {
		return true, nil
	}
}
