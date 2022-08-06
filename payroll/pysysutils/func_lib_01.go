package pysysutils

import (
	"gopayroll/common"
	"fmt"
	"github.com/go-xorm/xorm"
)

type TempCatalog struct {
	TenantId   int64
	HhrEmpid   string
	HhrEmpRcd  int64
	HhrSeqNum  int64
	HhrHistSeq int64
	HhrUpdSeq  int64
	CalGrpId   string `xorm:"hhr_pycalgrp_id"`
}

func GetGlobalKey(tenantId int64, calGrpId string) string {
	globalKey := fmt.Sprintf("%v:%s", tenantId, calGrpId)
	return globalKey
}

/**
删除标记/计算结果表、标记/计算消息表中的对应记录
@param tenant_id: 租户ID
@param cal_grp_id: 日历组ID
@param action_user: 操作用户
@param option: PAY-只处理薪资计算记录； ALL-删除标记和薪资计算记录
*/
func DelPayeeCalcStat(session *xorm.Session, tenantId int64, calGrpId string, actionUser string, option string) error {
	var err error
	if option == "ALL" {
		_, err = session.Exec("delete from hhr_payroll.hhr_py_payee_calc_stat where tenant_id = ? and hhr_pycalgrp_id = ?", tenantId, calGrpId)
		_, err = session.Exec("delete from hhr_payroll.hhr_py_payee_calc_msg where tenant_id = ? and hhr_pycalgrp_id = ?", tenantId, calGrpId)
	} else if option == "PAY" {
		_, err = session.Exec("update hhr_payroll.hhr_py_payee_calc_stat set hhr_py_calc_stat = '', hhr_modify_dttm = ?, hhr_modify_user = ? "+
			"where tenant_id = ? and hhr_pycalgrp_id = ?", common.CurLocalDate(), actionUser, tenantId, calGrpId)
		_, err = session.Exec("delete from hhr_payroll.hhr_py_payee_calc_msg where tenant_id = ? and hhr_pycalgrp_id = ? "+
			"and hhr_py_msg_class = 'B'", tenantId, calGrpId)
	}
	if err != nil {
		return err
	}
	return nil
}

/**
删除人员薪资计算结果
@param tenant_id: 租户ID
@param cal_grp_id: 日历组ID
@param action_user: 操作用户
*/
func DelPayCalcResult(session *xorm.Session, tenantId int64, calGrpId string, actionUser string) error {
	var err error
	var exists = false
	//清除/更新人员薪资计算目录数据
	catalogSlice := make([]*TempCatalog, 0)
	if err = session.Table("hhr_py_cal_catalog").Where("tenant_id=? and hhr_pycalgrp_id=?", tenantId, calGrpId).Find(&catalogSlice); err != nil {
		return err
	}
	//需要更新的记录：所在期日历组为当前日历组的记录中，历史序号列记录的序号对应的记录。从P或者O(历史)还原为A(活动)
	t1 := "update hhr_payroll.hhr_py_cal_catalog A set A.hhr_py_rec_stat = 'A', A.hhr_modify_dttm = ?, A.hhr_modify_user = ? " +
		"where (A.hhr_py_rec_stat = 'P' or A.hhr_py_rec_stat = 'O') and A.tenant_id = ? and A.hhr_empid = ? and A.hhr_emp_rcd = ? and A.hhr_seq_num = ?"
	//清除薪资结果-变量
	t2 := "delete from hhr_payroll.hhr_py_rslt_var where tenant_id = ? and hhr_empid = ? and hhr_emp_rcd = ? and hhr_seq_num = ?"
	//清除薪资结果-期间分段
	t3 := "delete from hhr_payroll.hhr_py_rslt_seg where tenant_id = ? and hhr_empid = ? and hhr_emp_rcd = ? and hhr_seq_num = ?"
	//清除薪资结果-任职
	t4 := "delete from hhr_payroll.hhr_py_rslt_job where tenant_id = ? and hhr_empid = ? and hhr_emp_rcd = ? and hhr_seq_num = ?"
	//清除薪资结果-工作日历
	t5 := "delete from hhr_payroll.hhr_py_rslt_wkcal where tenant_id = ? and hhr_empid = ? and hhr_emp_rcd = ? and hhr_seq_num = ?"
	//清除薪资结果-分段因子hhr_sal_plan_cd
	t6 := "delete from hhr_payroll.hhr_py_rslt_segfact where tenant_id = ? and hhr_empid = ? and hhr_emp_rcd = ? and hhr_seq_num = ?"
	//清除薪资结果-考勤记录
	t7 := "delete from hhr_payroll.hhr_py_rslt_abs where tenant_id = ? and hhr_empid = ? and hhr_emp_rcd = ? and hhr_seq_num = ?"
	//清除薪资结果-薪资项目
	t8 := "delete from hhr_payroll.hhr_py_rslt_pin where tenant_id = ? and hhr_empid = ? and hhr_emp_rcd = ? and hhr_seq_num = ?"
	//清除薪资结果-薪资项目累计
	t9 := "delete from hhr_payroll.hhr_py_rslt_accm where tenant_id = ? and hhr_empid = ? and hhr_emp_rcd = ? and hhr_seq_num = ?"
	//清除可写数组
	t10 := "delete from hhr_payroll.hhr_py_rslt_array_01 where tenant_id = ? and hhr_empid = ? and hhr_emp_rcd = ? and hhr_seq_num = ?"

	for i := 0; i < len(catalogSlice); i++ {
		exists = true
		empId := catalogSlice[i].HhrEmpid
		empRcd := catalogSlice[i].HhrEmpRcd
		seqNum := catalogSlice[i].HhrSeqNum
		histSeq := catalogSlice[i].HhrHistSeq
		updSeq := catalogSlice[i].HhrUpdSeq
		if updSeq != 0 {
			histSeq = updSeq
		}
		if histSeq != 0 {
			_, err = session.Exec(t1, common.CurLocalDate(), actionUser, tenantId, empId, empRcd, seqNum)
			if err != nil {
				return err
			}
		}

		_, err = session.Exec(t2, tenantId, empId, empRcd, seqNum)
		if err != nil {
			return err
		}
		_, err = session.Exec(t3, tenantId, empId, empRcd, seqNum)
		if err != nil {
			return err
		}
		_, err = session.Exec(t4, tenantId, empId, empRcd, seqNum)
		if err != nil {
			return err
		}
		_, err = session.Exec(t5, tenantId, empId, empRcd, seqNum)
		if err != nil {
			return err
		}
		_, err = session.Exec(t6, tenantId, empId, empRcd, seqNum)
		if err != nil {
			return err
		}
		_, err = session.Exec(t7, tenantId, empId, empRcd, seqNum)
		if err != nil {
			return err
		}
		_, err = session.Exec(t8, tenantId, empId, empRcd, seqNum)
		if err != nil {
			return err
		}
		_, err = session.Exec(t9, tenantId, empId, empRcd, seqNum)
		if err != nil {
			return err
		}
		_, err = session.Exec(t10, tenantId, empId, empRcd, seqNum)
		if err != nil {
			return err
		}
	}

	if exists == true {
		t11 := "delete from hhr_payroll.hhr_py_cal_catalog where tenant_id =? and hhr_pycalgrp_id =?"
		_, err = session.Exec(t11, tenantId, calGrpId)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
删除人员薪资计算过程日志
@param tenant_id: 租户ID
@param cal_grp_id: 日历组ID
*/
func DelPayCalcLog(session *xorm.Session, tenantId int64, calGrpId string) error {
	logTreeSql := "delete from hhr_payroll.hhr_py_log_tree where tenant_id = ? and hhr_py_cal_id = ?"
	_, err := session.Exec(logTreeSql, tenantId, calGrpId)
	if err != nil {
		return err
	}
	wtLogSql := "delete from hhr_payroll.hhr_py_wt_log where tenant_id = ? and hhr_py_cal_id = ?"
	_, err = session.Exec(wtLogSql, tenantId, calGrpId)
	if err != nil {
		return err
	}
	wcLogSql := "delete from hhr_payroll.hhr_py_wc_log where tenant_id = ? and hhr_py_cal_id = ?"
	_, err = session.Exec(wcLogSql, tenantId, calGrpId)
	if err != nil {
		return err
	}
	vrLogSql := "delete from hhr_payroll.hhr_py_vr_log where tenant_id = ? and hhr_py_cal_id = ?"
	_, err = session.Exec(vrLogSql, tenantId, calGrpId)
	if err != nil {
		return err
	}
	fcParamSql := "delete from hhr_payroll.hhr_py_fc_param_log where tenant_id = ? and hhr_py_cal_id = ?"
	_, err = session.Exec(fcParamSql, tenantId, calGrpId)
	if err != nil {
		return err
	}
	return nil
}

/**
删除已标记的人员
@param tenant_id: 租户ID
@param cal_grp_id: 日历组ID
*/
func DelIdentifiedPayees(session *xorm.Session, tenantId int64, calGrpId string) error {
	statSql := "delete from hhr_payroll.hhr_py_payee_calc_stat where tenant_id = ? and hhr_py_cal_id = ?"
	_, err := session.Exec(statSql, tenantId, calGrpId)
	if err != nil {
		return err
	}
	msgSql := "delete from hhr_payroll.hhr_py_payee_calc_msg where tenant_id = ? and hhr_py_cal_id = ?"
	_, err = session.Exec(msgSql, tenantId, calGrpId)
	if err != nil {
		return err
	}
	return nil
}
