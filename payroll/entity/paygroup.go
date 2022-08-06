/*
Desc: 薪资组实体
Author: 潘承勋
Date: 2020-05-03
*/

package entity

import (
	"gopayroll/common"
	"gopayroll/db"
	"time"
)

type PayGroup struct {
	TenantId       int64
	//薪资组编码
	PyGroupId      string    `xorm:"hhr_pygroup_id"`
	Country        string    `xorm:"hhr_country"`
	//生效日期
	EffDate        time.Time `xorm:"hhr_efft_date"`
	//状态Y-有效，N-无效
	Status         string    `xorm:"hhr_status"`
	//期间ID
	PeriodCd       string    `xorm:"hhr_period_code"`
	//追溯标志,Y-追溯，N-不追溯
	RetroFlag      string    `xorm:"hhr_retro_flag"`
	//分段规则
	SegRuleCd      string    `xorm:"hhr_seg_rule_cd"`
	//最大回溯期间数
	MaxRetroPrdNum int32     `xorm:"hhr_max_retro_prd_num"`
	//取整规则
	RoundRule      string    `xorm:"hhr_round_rule"`
	//工作计划
	WorkPlan       string    `xorm:"hhr_work_plan"`
	//货币
	Currency       string    `xorm:"hhr_currency"`
	//支付银行
	PayBank        string    `xorm:"hhr_payment_bank"`
	//考勤评估规则
	AbsEvalRule    string    `xorm:"hhr_abs_eval_rule"`
	//考勤期间差异
	AbsPrdDiffSw   string    `xorm:"hhr_abs_prd_diff_sw"`
	//考勤期间
	AbsPeriod      string    `xorm:"hhr_abs_period"`
	//期间序号差异
	PrdSeqDiffSw   string    `xorm:"hhr_prd_seq_diff_sw"`
	//差异数
	DiffNum        int32     `xorm:"hhr_diff_num"`
	//采用考勤周期薪资标准
	AbsCycleStdSw  string    `xorm:"hhr_abs_cycl_std_sw"`
	//分段规则实体
	SegRuleEntity  *SegRule   `xorm:"-"`
}

func NewPayGroup(tenantId int64, payGroupId string, toDate time.Time) *PayGroup {

	payGroup := &PayGroup{TenantId: tenantId, PyGroupId: payGroupId}
	engine := db.OrmEngine("hhr_payroll")
	_, err := engine.Table("hhr_py_paygroup").
		Where("tenant_id=? and ? between hhr_efft_date and hhr_efft_end_date and hhr_status = 'Y' ", tenantId, toDate).Get(payGroup)
	if err != nil {
		common.Logger.Error("[NewPayGroup-1]", err.Error())
	}
	payGroup.SegRuleEntity = NewSegRule(tenantId, payGroupId, payGroup.EffDate)
	common.Logger.Debug("------NewPayGroup--------", payGroup)
	return payGroup
}
