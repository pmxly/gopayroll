package entity

type Catalog struct {
	TenantId  int64
	Empid     string `xorm:"hhr_empid"`
	EmpRcd    int64  `xorm:"hhr_emp_rcd"`
	SeqNum    int64  `xorm:"hhr_seq_num"`
	FCalGrpId string `xorm:"hhr_f_calgrp_id"`
	FCalId    string `xorm:"hhr_f_cal_id"`
}
