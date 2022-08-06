package entity

import "time"

type PayeeCalcMsg struct {
	TenantId   int64
	CalGrpId   string    `xorm:"hhr_pycalgrp_id"`
	CalId      string    `xorm:"hhr_py_cal_id"`
	EmpId      string    `xorm:"hhr_empid"`
	EmpRcd     int64     `xorm:"hhr_emp_rcd"`
	MsgClass   string    `xorm:"hhr_py_msg_class"`
	SeqNum     int64     `xorm:"hhr_seq_num"`
	FCalId     string    `xorm:"hhr_f_cal_id"`
	MsgType    string    `xorm:"hhr_py_msg_type"`
	MsgTxt     string    `xorm:"hhr_msg_txt"`
	CreateDtm  time.Time `xorm:"hhr_create_dttm"`
	CreateUser string    `xorm:"hhr_create_user"`
	ModifyDtm  time.Time `xorm:"hhr_modify_dttm"`
	ModifyUser string    `xorm:"hhr_modify_user"`
}
