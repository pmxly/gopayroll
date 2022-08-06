/*
Desc: 标记受款人
Author: 潘承勋
Date: 2020-05-03
*/

package identify

import (
	"gopayroll/common"
	"gopayroll/db"
	pyCommon "gopayroll/payroll/common"
	"gopayroll/payroll/entity"
	"gopayroll/payroll/module"
	"gopayroll/payroll/pysysutils"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
)

type PayeesIdentify struct {
	CalGrpObj    *entity.CalendarGroup
	CalendarDic  map[string]*entity.Calendar
	RunParam     *module.PayCalcRunParam
	LockedPayees []*entity.PayeeCalcStat
	Engine       *xorm.Engine
}

func NewPayeesIdentify(globalKey string, runParam *module.PayCalcRunParam) *PayeesIdentify {
	runVarCache := pyCommon.RunVarCache(globalKey)
	pi := &PayeesIdentify{}
	pi.CalGrpObj = runVarCache["CAL_GRP_OBJ"].Value.(*entity.CalendarGroup)
	pi.CalendarDic = runVarCache["CAL_OBJ_DIC"].Value.(map[string]*entity.Calendar)
	pi.RunParam = runParam
	pi.LockedPayees = make([]*entity.PayeeCalcStat, 0)
	pi.Engine = db.OrmEngine("hhr_payroll")
	return pi
}

func (pi *PayeesIdentify) IdentifyPayees() error {
	//1.数据一致性校验
	if err := pi.checkConsistency(); err != nil {
		return err
	}
	//2.更新/删除原操作记录和结果
	if err := pi.clrPayeeCalcStat(); err != nil {
		return err
	}
	if err := pi.clrPayCalcResult(); err != nil {
		return err
	}
	if err := pi.clrPayCalcLog(); err != nil {
		return err
	}
	if err := pi.clrIdentifiedPayees(); err != nil {
		return err
	}
	if err := pi.beginIdentify(); err != nil {
		return err
	}
	return nil
}

//数据一致性校验
func (pi *PayeesIdentify) checkConsistency() error {
	if len(pi.CalendarDic) == 0 {
		common.Logger.Error("[checkConsistency] no calendar exists")
		return errors.New("[checkConsistency] no calendar exists")
	}
	var errText string
	var country = pi.CalGrpObj.Country
	var calLst []*entity.Calendar
	var payGrplst []*entity.PayGroup
	var runTypeLst []*entity.RunType
	for _, cal := range pi.CalendarDic {
		calLst = append(calLst, cal)
		payGrplst = append(payGrplst, cal.PyGroupEntity)
		runTypeLst = append(runTypeLst, cal.RunTypeEntity)
	}

	//1.日历组的国家/地区、日历的薪资组的国家/地区、日历的运行类型的国家/地区必须一致
	for _, pg := range payGrplst {
		if pg.Country != country {
			errText = fmt.Sprintf("日历组%s的国家/地区与薪资组%v的国家/地区不一致", pi.CalGrpObj.CalGrpId, pg.PyGroupId)
			return errors.New(errText)
		}
	}
	for _, rt := range runTypeLst {
		if rt.Country != country {
			errText = fmt.Sprintf("日历组%s的国家/地区与运行类型%v的国家/地区不一致", pi.CalGrpObj.CalGrpId, rt.RunTypeId)
			return errors.New(errText)
		}
	}

	//2.日历的期间编码、日历的薪资组的期间编码必须一致
	for _, cal := range calLst {
		if cal.PeriodId != cal.PyGroupEntity.PeriodCd {
			errText = fmt.Sprintf("日历%s的期间编码与薪资组%v的期间编码不一致", cal.CalId, cal.PyGroupEntity.PyGroupId)
			return errors.New(errText)
		}
	}
	//3.日历组下所有日历的期间年度、期间序号必须一致
	//4.日历组下所有日历的计算类型必须一致
	//5.日历组下所有日历的运行类型的性质必须一致
	//6.当日历组勾选了追溯时，校验日历的运行类型的性质，必须为“C-周期”。为“O-非周期”时，日历组不能勾选追溯
	if len(calLst) > 0 {
		var year = calLst[0].PeriodYear
		var seq = calLst[0].PeriodNum
		var calType = calLst[0].PayCalType
		var cycle = calLst[0].RunTypeEntity.Cycle
		for _, cal := range calLst {
			if len(calLst) > 1 {
				if year != cal.PeriodYear {
					errText = fmt.Sprintf("日历%s的期间年度%v与日历%s的期间年度%v不一致", calLst[0].CalId, year, cal.CalId, cal.PeriodYear)
					return errors.New(errText)
				}
				if seq != cal.PeriodNum {
					errText = fmt.Sprintf("日历%s的期间序号%v与日历%s的期间序号%v不一致", calLst[0].CalId, seq, cal.CalId, cal.PeriodNum)
					return errors.New(errText)
				}
				if calType != cal.PayCalType {
					errText = fmt.Sprintf("日历%s的计算类型%v与日历%s的计算类型%v不一致", calLst[0].CalId, calType, cal.CalId, cal.PayCalType)
					return errors.New(errText)
				}
				if cycle != cal.RunTypeEntity.Cycle {
					errText = fmt.Sprintf("日历%s的运行类型的性质%v与日历%s的运行类型的性质%v不一致", calLst[0].CalId, cycle, cal.CalId, cal.RunTypeEntity.Cycle)
					return errors.New(errText)
				}
			}
			if cal.RunTypeEntity.Cycle != "C" {
				errText = fmt.Sprintf("日历组%s启用了追溯，但是日历%s的运行类型的性质并不是C-周期", pi.CalGrpObj.CalGrpId, cal.CalId)
				return errors.New(errText)
			}
		}
	}
	return nil
}

func (pi *PayeesIdentify) clrPayeeCalcStat() error {
	session := pi.Engine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Logger.Error("[clrPayeeCalcStat]", err.Error())
	}
	var tenantId = pi.RunParam.TenantId
	var calGrpId = pi.CalGrpObj.CalGrpId
	var userId = pi.RunParam.UserId
	err = pysysutils.DelPayeeCalcStat(session, tenantId, calGrpId, userId, "ALL")
	if err != nil {
		return err
	}
	return session.Commit()
}

/**
清除已保存的薪资计算结果
*/
func (pi *PayeesIdentify) clrPayCalcResult() error {
	session := pi.Engine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Logger.Error("[clrPayCalcResult]", err.Error())
	}
	var tenantId = pi.RunParam.TenantId
	var calGrpId = pi.CalGrpObj.CalGrpId
	var userId = pi.RunParam.UserId
	err = pysysutils.DelPayCalcResult(session, tenantId, calGrpId, userId)
	if err != nil {
		return err
	}
	return session.Commit()
}

/**
清除已保存的过程日志
*/
func (pi *PayeesIdentify) clrPayCalcLog() error {
	session := pi.Engine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Logger.Error("[clrPayCalcLog]", err.Error())
	}
	var tenantId = pi.RunParam.TenantId
	var calGrpId = pi.CalGrpObj.CalGrpId
	err = pysysutils.DelPayCalcLog(session, tenantId, calGrpId)
	if err != nil {
		return err
	}
	return session.Commit()
}

/**
清除已标记的人员
*/
func (pi *PayeesIdentify) clrIdentifiedPayees() error {
	session := pi.Engine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Logger.Error("[clrIdentifiedPayees]", err.Error())
	}
	var tenantId = pi.RunParam.TenantId
	var calGrpId = pi.CalGrpObj.CalGrpId
	err = pysysutils.DelIdentifiedPayees(session, tenantId, calGrpId)
	if err != nil {
		return err
	}
	return session.Commit()
}

/**
开始标记受款人
*/
func (pi *PayeesIdentify) beginIdentify() error {
	session := pi.Engine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Logger.Error("[beginIdentify-1]", err.Error())
		return err
	}
	var tenantId = pi.RunParam.TenantId
	//一次性获取所有被其他日历锁定的人员任职记录
	err = session.Table("hhr_py_payee_calc_stat").Cols("hhr_empid", "hhr_emp_rcd", "hhr_py_cal_id").
		Where("tenant_id=? and hhr_lock_flg='1'", tenantId).Find(&pi.LockedPayees)
	if err != nil {
		common.Logger.Error("[beginIdentify-2]", err.Error())
		return err
	}
	//处理日历组的每个日历
	for _, cal := range pi.CalendarDic {
		prdBgnDt := cal.PeriodCalEntity.StartDate
		prdEndDt := cal.PeriodCalEntity.EndDate
		pyGroupId := cal.PyGroupId
		calType := cal.PayCalType
		//针对每个日历，一次性获取所有后续期间已存在计算结果的人员任职记录
		furtherPayees := make([]*entity.Payee, 0)
		err = session.Table("hhr_py_cal_catalog").Cols("hhr_empid", "hhr_emp_rcd").Distinct("hhr_empid", "hhr_emp_rcd").
			Where("tenant_id=? and hhr_f_prd_end_dt=?", tenantId, prdEndDt).Find(&furtherPayees)
		if err != nil {
			common.Logger.Error("[beginIdentify-3]", err.Error())
			return err
		}

		//针对每个日历，一次性获取当前期间已存在周期计算结果的人员
		cyclePayees := make([]*entity.Payee, 0)
		if cal.RunTypeEntity.Cycle == "C" && (calType == "A" || calType == "B") {
			err = session.Table("hhr_py_cal_catalog").Cols("hhr_empid", "hhr_emp_rcd").Distinct("hhr_empid", "hhr_emp_rcd").
				Where("tenant_id=? and hhr_f_rt_cycle='C' and (hhr_f_prd_bgn_dt <= ? and hhr_f_prd_end_dt >= ?) ",
					tenantId, prdEndDt, prdBgnDt).Find(&cyclePayees)
			if err != nil {
				common.Logger.Error("[beginIdentify-4]", err.Error())
				return err
			}
		}

		//常规
		if calType == "A" {
			//指定薪资组在薪资计算的期间内在职的人员（至少1天在职）
			//根据薪资组、期间结束日期从基础薪酬表中获取员工任职（在期间结束日期属于此薪资组）
			stmt1 := "select a.hhr_empid, a.hhr_emp_rcd from hhr_payroll.hhr_py_assign_pg a where a.tenant_id = ? " +
				"and ? between a.hhr_efft_date and a.hhr_efft_end_date " +
				"and a.hhr_efft_seq = (select max(a2.hhr_efft_seq) from hhr_payroll.hhr_py_assign_pg a2 where a2.tenant_id = a.tenant_id " +
				"and a2.hhr_empid = a.hhr_empid and a2.hhr_emp_rcd = a.hhr_emp_rcd and a2.hhr_efft_date = a.hhr_efft_date) " +
				"and a.hhr_pygroup_id = ? "
			stmt2 := "select A.hhr_empid from hhr_corehr.hhr_org_per_jobdata A where A.tenant_id = ? and A.hhr_efft_date > ? " +
				"and A.hhr_efft_date <= ? and A.hhr_status = 'Y' and A.hhr_empid = ? and A.hhr_emp_rcd = ? " +
				"union " +
				"select a.hhr_empid from hhr_corehr.hhr_org_per_jobdata a where a.tenant_id = ? and a.hhr_status = 'Y' " +
				"and ? between a.hhr_efft_date and a.hhr_efft_end_date and a.hhr_empid = ? and a.hhr_emp_rcd = ? "
			empIdRcdSlc := make([]*module.EmpIdRcd, 0)
			err = session.SQL(stmt1, tenantId, prdEndDt, pyGroupId).Find(&empIdRcdSlc)
			if err != nil {
				common.Logger.Error("[beginIdentify-5]", err.Error())
				return err
			}
			for _, empIdRcd := range empIdRcdSlc {
				empId := empIdRcd.EmpId
				empRcd := empIdRcd.EmpRcd
				if cal.ShouldBeExcluded(empId, empRcd) {
					continue
				}

				exists, err := session.SQL(stmt2, tenantId, prdBgnDt, prdEndDt, empId, empRcd, tenantId, prdBgnDt, empId, empRcd).Exist()
				if err != nil {
					common.Logger.Error("[beginIdentify-6]", err.Error())
					return err
				}
				if exists {
					payee := &entity.Payee{
						TenantId: tenantId,
						CalId:    cal.CalId,
						EmpId:    empId,
						EmpRcd:   empRcd,
					}
					err = pi.payCheck(session, cal, payee, furtherPayees, cyclePayees)
					if err != nil {
						common.Logger.Error("[beginIdentify-7]", err.Error())
						return err
					}
				}
			}
		}

	}
	return session.Commit()
}

/**
校验受款人并标记
*/
func (pi *PayeesIdentify) payCheck(session *xorm.Session, cal *entity.Calendar, payee *entity.Payee, furPayees []*entity.Payee, cyclePayees []*entity.Payee) error {
	var err error
	calType := cal.PayCalType
	userId := pi.RunParam.UserId
	tenantId := payee.TenantId
	calGrpId := pi.CalGrpObj.CalGrpId
	calId := cal.CalId
	empId := payee.EmpId
	empRcd := payee.EmpRcd

	//1.校验人员是否已被其他日历锁定，若是则此人员标记失败
	for _, lockPayee := range pi.LockedPayees {
		if lockPayee.EmpId == empId && lockPayee.EmpRcd == empRcd {
			lockCalId := lockPayee.CalId
			msgTxt := fmt.Sprintf("该员工已被日历%s锁定", lockCalId)
			err = entity.InsPyCalcMsg(session, calGrpId, payee, cal, "A", "E", msgTxt, userId)
			if err != nil {
				common.Logger.Error("[payCheck-1]", err.Error())
				return err
			}
		}
	}
	//2.校验人员后续期间是否已存在计算结果，若已存在则此人员标记失败（存在已有薪资计算结果的结束日期>当前结束日期，标记失败）
	for _, furPayee := range furPayees {
		if furPayee.EmpId == empId && furPayee.EmpRcd == empRcd {
			msgTxt := "该员工后续期间已存在计算结果"
			err = entity.InsPyCalcMsg(session, calGrpId, payee, cal, "A", "E", msgTxt, userId)
			if err != nil {
				common.Logger.Error("[payCheck-2]", err.Error())
				return err
			}
		}
	}
	//3.当运行类型的性质为“C-周期”时，对常规计算或单独计算的人员，校验人员当前期间是否已存在历经期运行类型为周期的计算结果，
	//若已存在则此人员标记失败（存在已有周期薪资计算结果的期间与当前期间交叉，标记失败）
	for _, cyclePayee := range cyclePayees {
		if cyclePayee.EmpId == empId && cyclePayee.EmpRcd == empRcd {
			msgTxt := "该员工已存在周期的计算结果"
			err = entity.InsPyCalcMsg(session, calGrpId, payee, cal, "A", "E", msgTxt, userId)
			if err != nil {
				common.Logger.Error("[payCheck-3]", err.Error())
				return err
			}
		}
	}
	//4.对常规计算的增补人员，校验期间结束日期的薪资组是否一致，若不一致则此人员标记失败
	//5.对单独计算的人员，校验期间结束日期的薪资组是否一致，若不一致则此人员标记失败
	//6.对更正计算的人员，校验期间结束日期的薪资组是否一致，若不一致则此人员标记失败
	checkPass, err := entity.CheckPayeePayGroup(session, payee, cal, common.ZeroTime, common.ZeroTime)
	if err != nil {
		common.Logger.Error("[payCheck-4]", err.Error())
		return err
	}
	if checkPass == false {
		msgTxt := "员工期间结束日期内的薪资组与日历的薪资组不一致"
		err = entity.InsPyCalcMsg(session, calGrpId, payee, cal, "A", "E", msgTxt, userId)
		if err != nil {
			common.Logger.Error("[payCheck-5]", err.Error())
			return err
		}
	}
	//7.对更正计算的人员，校验更正日历的期间编码、年度、序号、运行类型是否相同，并且是此人员序号最大的所在期日历，若不一致则此人员标记失败
	if calType == "C" {
		uCalId := payee.UpdCalId
		uPrdId := payee.UpdCalEntity.PeriodId
		uPrdYear := payee.UpdCalEntity.PeriodYear
		uPrdNum := payee.UpdCalEntity.PeriodNum
		uRunType := payee.UpdCalEntity.RunTypeId

		calPrdId := cal.PeriodId
		calPrdYear := cal.PeriodYear
		calPrdNum := cal.PeriodNum
		calRunType := cal.RunTypeId
		if (uPrdId != calPrdId) || (uPrdYear != calPrdYear) || (uPrdNum != calPrdNum) || (uRunType != calRunType) {
			msgTxt := fmt.Sprintf("更正日历%s的期间编码、年度、序号、运行类型与运行日历%s不一致", uCalId, calId)
			err = entity.InsPyCalcMsg(session, calGrpId, payee, cal, "A", "E", msgTxt, userId)
			if err != nil {
				common.Logger.Error("[payCheck-6]", err.Error())
				return err
			}
		} else {
			//获取人员序号最大的所在期日历
			t1 := "select hhr_py_cal_id from hhr_payroll.hhr_py_cal_catalog a where a.tenant_id = ? and a.hhr_empid = ? " +
				"and hhr_emp_rcd = ? order by a.hhr_seq_num desc"
			var preCalId string
			_, err = session.SQL(t1, tenantId, empId, empRcd).Get(&preCalId)
			if err != nil {
				common.Logger.Error("[payCheck-7]", err.Error())
				return err
			}
			if preCalId != uCalId {
				common.Logger.Error("[payCheck-8]",  errors.New("更正日历必须是此人员序号最大的所在期日历"))
				return errors.New("更正日历必须是此人员序号最大的所在期日历")
			}
		}
	}
	return nil
}
