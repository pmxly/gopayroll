/*
Desc: 期间、期间日历实体
Author: 潘承勋
Date: 2020-05-03
*/

package entity

import (
	"gopayroll/common"
	"gopayroll/db"
	"time"
)

//期间
type Period struct {
	TenantId int64
	PeriodId string `xorm:"hhr_period_code"`
	/*频率
	  10:每月
	  20:每半月
	  30:每周
	  40:每两周
	  50:每四周
	  60:每年
	  70:每季度
	  80:每半年*/
	Frequency string `xorm:"hhr_period_unit"`
}

//期间日历
type PeriodCalender struct {
	TenantId       int64
	PeriodId       string    `xorm:"hhr_period_code"`
	Year           int32     `xorm:"hhr_period_year"`
	PeriodNum      int32     `xorm:"hhr_prd_num"`
	StartDate      time.Time `xorm:"hhr_period_start_date"`
	EndDate        time.Time `xorm:"hhr_period_end_date"`
	LastPeriodYear int32     `xorm:"hhr_period_last_year"`
	LastPeriodNum  int32     `xorm:"hhr_last_prd_num"`
}

func NewPeriod(tenantId int64, periodId string) *Period {
	period := &Period{TenantId: tenantId, PeriodId: periodId}
	engine := db.OrmEngine("hhr_payroll")
	_, err := engine.Table("hhr_py_period").Where("tenant_id=?", tenantId).Cols("hhr_period_unit").Get(period)
	if err != nil {
		common.Logger.Error("[NewPeriod]", err.Error())
	}
	common.Logger.Debug("------NewPeriod--------", period)
	return period
}

func NewPeriodCalendar(tenantId int64, periodId string, periodYear int32, periodNum int32) *PeriodCalender {

	periodCal := &PeriodCalender{TenantId: tenantId, PeriodId: periodId, Year: periodYear, PeriodNum: periodNum}
	engine := db.OrmEngine("hhr_payroll")
	_, err := engine.Table("hhr_py_period_calendar_line").Where("tenant_id=?", tenantId).Get(periodCal)
	if err != nil {
		common.Logger.Error("[NewPeriodCalendar]", err.Error())
	}
	common.Logger.Debug("------NewPeriodCalendar--------", periodCal)
	return periodCal
}