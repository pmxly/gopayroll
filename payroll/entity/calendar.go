/**
Desc: 日历、日历组、受款人实体
Author: 潘承勋
Date: 2020-05-03
*/

package entity

import (
	"gopayroll/common"
	"gopayroll/db"
	"time"
)

//日历
type Calendar struct {
	TenantId   int64
	CalId      string    `xorm:"hhr_py_cal_id"`
	PyGroupId  string    `xorm:"hhr_pygroup_id"`
	RunTypeId  string    `xorm:"hhr_runtype_id"`
	PeriodId   string    `xorm:"hhr_period_code"`
	PeriodYear int32     `xorm:"hhr_period_year"`
	PeriodNum  int32     `xorm:"hhr_prd_num"`
	PayDate    time.Time `xorm:"hhr_pay_date"`
	//计算类型，A-常规，B-单独，C-更正
	PayCalType string `xorm:"hhr_pycalc_type"'`
	//薪资组实体
	PyGroupEntity *PayGroup `xorm:"-"`
	//运行类型实体
	RunTypeEntity *RunType `xorm:"-"`
	//期间实体
	PeriodEntity *Period `xorm:"-"`
	//期间日历实体
	PeriodCalEntity *PeriodCalender `xorm:"-"`
	//排除人员list
	ExceptPersList []*Payee `xorm:"-"`
	//手动添加人员list
	AddPersList []*Payee `xorm:"-"`
	//更正人员list
	UpdatePerList []*Payee `xorm:"-"`
	//实际计算人员list
	PersList []*Payee `xorm:"-"`
	//当前日历的最小追溯期间日历
	MinRtoPrdCal *PeriodCalender `xorm:"-"`
	//当前日历的最小追溯期间日历(新入职员工)
	MinRtoPrdCalForNewEntry *PeriodCalender `xorm:"-"`
	CalGrpId                string          `xorm:"-"`
}

//日历组
type CalendarGroup struct {
	TenantId int64
	CalGrpId string `xorm:"hhr_pycalgrp_id"`
	Country  string `xorm:"hhr_country"`
	Retro    string `xorm:"hhr_retro_flag"`
}

func NewCalendarGroup(tenantId int64, calGrpId string) *CalendarGroup {
	calGrp := &CalendarGroup{TenantId: tenantId, CalGrpId: calGrpId}
	engine := db.OrmEngine("hhr_payroll")
	_, err := engine.Table("hhr_py_cal_grp").Where("tenant_id=?", tenantId).Cols("hhr_country", "hhr_retro_flag").Get(calGrp)
	if err != nil {
		common.Logger.Error("[NewCalendarGroup]", err.Error())
	}
	common.Logger.Debug("------NewCalendarGroup--------", calGrp)
	return calGrp
}

func NewCalendar(tenantId int64, calId string) *Calendar {
	cal := &Calendar{TenantId: tenantId, CalId: calId}
	engine := db.OrmEngine("hhr_payroll")
	_, err := engine.Table("hhr_py_calendar").Where("tenant_id=?", tenantId).Get(cal)
	if err != nil {
		common.Logger.Error("[NewCalendar]", err.Error())
	}
	cal.PeriodEntity = NewPeriod(tenantId, cal.PeriodId)
	cal.PeriodCalEntity = NewPeriodCalendar(tenantId, cal.PeriodId, cal.PeriodYear, cal.PeriodNum)
	cal.PyGroupEntity = NewPayGroup(tenantId, cal.PyGroupId, cal.PeriodCalEntity.EndDate)
	cal.RunTypeEntity = NewRunType(tenantId, cal.RunTypeId, cal)
	cal.InitExceptAddList("EXPT")
	cal.InitExceptAddList("ADD")
	cal.InitExceptAddList("UPT")
	common.Logger.Debug("------NewCalendar--------", cal)
	return cal
}

func NewCalendarOnce(tenantId int64, calId string) *Calendar {
	cal := &Calendar{TenantId: tenantId, CalId: calId}
	engine := db.OrmEngine("hhr_payroll")
	_, err := engine.Table("hhr_py_calendar").Where("tenant_id=?", tenantId).Get(cal)
	if err != nil {
		common.Logger.Error("[NewCalendar]", err.Error())
	}
	cal.PeriodEntity = NewPeriod(tenantId, cal.PeriodId)
	cal.PeriodCalEntity = NewPeriodCalendar(tenantId, cal.PeriodId, cal.PeriodYear, cal.PeriodNum)
	cal.PyGroupEntity = NewPayGroup(tenantId, cal.PyGroupId, cal.PeriodCalEntity.EndDate)
	cal.RunTypeEntity = NewRunType(tenantId, cal.RunTypeId, cal)
	cal.InitExceptAddList("EXPT")
	cal.InitExceptAddList("ADD")
	cal.InitExceptAddList("UPT")
	common.Logger.Debug("------NewCalendar--------", cal)
	return cal
}

func (calGrp *CalendarGroup) GetCalendarMap(tenantId int64, calGrpId string) map[string]*Calendar {
	result := make(map[string]*Calendar)
	engine := db.OrmEngine("hhr_payroll")
	pyCalIdSlice := make([]string, 0)
	err := engine.Table("hhr_py_calgrp_cal").Cols("hhr_py_cal_id").
		Where("tenant_id=? and hhr_pycalgrp_id=?", tenantId, calGrpId).Find(&pyCalIdSlice)
	if err != nil {
		common.Logger.Error("[GetCalendarMap]", err.Error())
	}
	for i := 0; i < len(pyCalIdSlice); i++ {
		cal := NewCalendar(tenantId, pyCalIdSlice[i])
		result[pyCalIdSlice[i]] = cal
	}
	//common.Logger.Info("[GetCalendarMap]", result)
	return result
}

/**
初始化人员列表
*/
func (cal *Calendar) InitExceptAddList(flag string) {
	engine := db.OrmEngine("hhr_payroll")
	var payeesPtr *[]*Payee
	var tableName = ""

	if flag == "EXPT" {
		tableName = "hhr_py_cal_emp_expt"
		payeesPtr = &cal.ExceptPersList
	} else if flag == "ADD" {
		tableName = "hhr_py_cal_emp_add"
		payeesPtr = &cal.AddPersList
	} else if flag == "UPT" {
		tableName = "hhr_py_cal_emp_upt"
		payeesPtr = &cal.UpdatePerList
	}

	if tableName != "" {
		if flag == "UPT" {
			err := engine.Table(tableName).Cols("hhr_empid", "hhr_emp_rcd", "hhr_upt_calid").
				Where("tenant_id=? and hhr_py_cal_id=?", cal.TenantId, cal.CalId).Find(payeesPtr)
			if err != nil {
				common.Logger.Error("[InitExceptAddList-UPT]", err.Error())
			}
			for _, uptPayee := range *payeesPtr {
				uptPayee.UpdCalEntity = NewCalendar(cal.TenantId, uptPayee.UpdCalId)
			}
		} else {
			err := engine.Table(tableName).Cols("hhr_empid", "hhr_emp_rcd").
				Where("tenant_id=? and hhr_py_cal_id=?", cal.TenantId, cal.CalId).Find(payeesPtr)
			if err != nil {
				common.Logger.Error("[InitExceptAddList-OTHERS]", err.Error())
			}
		}

	}
}

func (cal *Calendar) ShouldBeExcluded(empId string, empRcd int64) bool {
	for _, payee := range cal.ExceptPersList {
		if payee.EmpId == empId && payee.EmpRcd == empRcd {
			return true
		}
	}
	return false
}
