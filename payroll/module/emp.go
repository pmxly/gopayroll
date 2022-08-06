package module


type EmpIdRcd struct {
	TenantId int64
	EmpId    string `xorm:"hhr_empid"`
	EmpRcd   int64  `xorm:"hhr_emp_rcd"`
}

