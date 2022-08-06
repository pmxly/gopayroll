package entity

//受款人
type Payee struct {
	TenantId     int64
	CalId        string    `xorm:"hhr_py_cal_id"`
	EmpId        string    `xorm:"hhr_empid"`
	EmpRcd       int64     `xorm:"hhr_emp_rcd"`
	UpdCalId     string    `xorm:"hhr_upt_calid"`
	UpdCalEntity *Calendar `xorm:"-"`
}