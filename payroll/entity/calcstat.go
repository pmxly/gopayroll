package entity

type PayeeCalcStat struct {
	TenantId int64
	CalGrpId string `xorm:"hhr_pycalgrp_id"`
	CalId    string `xorm:"hhr_py_cal_id"`
	EmpId    string `xorm:"hhr_empid"`
	EmpRcd   int64  `xorm:"hhr_emp_rcd"`
	UpdCalId string `xorm:"hhr_upt_calid"`
	LockFlag string `xorm:"hhr_lock_flg"`
	IndStat  string `xorm:"hhr_py_ind_stat"`
	CalcStat string `xorm:"hhr_py_calc_stat"`
}
