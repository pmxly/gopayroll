package entity

import (
	"gopayroll/common"
	"gopayroll/db"
	"time"
)

type SegRule struct {
	TenantId    int64
	SegRuleCd   string    `xorm:"hhr_seg_rule_cd"`
	Country     string    `xorm:"hhr_country"`
	EffDate     time.Time `xorm:"hhr_efft_date"`
	Numerator   string    `xorm:"hhr_numerator"`
	Denominator string    `xorm:"hhr_denominator"`
	FixedDays   float64   `xorm:"decimal hhr_fixed_days"`
	Priority    string    `xorm:"hhr_priority"`
}

func NewSegRule(tenantId int64, segRuleCd string, toDate time.Time) *SegRule {
	segRule := &SegRule{TenantId: tenantId, SegRuleCd: segRuleCd}
	engine := db.OrmEngine("hhr_payroll")
	_, err := engine.Table("hhr_py_seg_rule").
		Where("tenant_id=? and ? between hhr_efft_date and hhr_efft_end_date and hhr_status = 'Y' ", tenantId, toDate).Get(segRule)
	if err != nil {
		common.Logger.Error("[NewSegRule]", err.Error())
	}
	common.Logger.Debug("------NewSegRule--------", segRule)
	return segRule
}
