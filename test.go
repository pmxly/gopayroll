package main

import (
	"gopayroll/common"
	"gopayroll/db"
	"gopayroll/payroll/entity"
	"gopayroll/payroll/module"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type PayCalcRunParam struct {
	UserId          string `json:"user_id"`
	TenantId        int    `json:"tenant_id"`
	CalGrpID        string `json:"cal_grp_id"`
	CalID           string `json:"cal_id"`
	IdentityFlag    string `json:"identity_flag"`
	ReIdentityFlag  string `json:"re_identity_flag"`
	PyCalcFlag      string `json:"py_calc_flag"`
	ReCalcFlag      string `json:"re_calc_flag"`
	ReIdentCalcFlag string `json:"re_ident_calc_flag"`
	CancelIdentFlag string `json:"cancel_ident_flag"`
	CancelCalcFlag  string `json:"cancel_calc_flag"`
	CancelFlag      string `json:"cancel_flag"`
	FinishFlag      string `json:"finish_flag"`
	CloseFlag       string `json:"close_flag"`
	LogFlag         string `json:"log_flag"`
	LogEmpRecStr    string `json:"log_emp_rec_str"`
	TaskName        string `json:"task_name"`
}

type test01 struct {
	a string
	t time.Time
}

func main() {
	//var s1 []*test01
	//s1 = append(s1, &test01{a: "xyz"},  &test01{a: "uyt"})
	//fmt.Println("s1", s1)
	//for i, v := range s1{
	//	fmt.Printf("%p:\n", v)
	//	fmt.Println("i:", i)
	//	v.a = "yyy"
	//}
	//fmt.Println("s1:", s1[1].a)
	//fmt.Println("t:", s1[1].t.IsZero())

	var a interface{}
	fmt.Println("a:", a == nil)

	var s1 = []string{"A", "B", "C"}
	var s2 = &s1
	*s2 = append(*s2, "X")
	fmt.Println("s1:", s1)
	fmt.Println("s2:", s2)

	/*t := time.Now()
	num:=fibonacci(45)
	fmt.Println(num)
	latency := time.Since(t)
	log.Print("耗时：", latency)*/

	db.InitDBEngine()
	type Payee struct {
		TenantId    int64
		CalId       string `xorm:"hhr_py_cal_id"`
		EmpId       string `xorm:"hhr_empid"`
		EmpRcd      int64  `xorm:"hhr_emp_rcd"`
		PyGroupId   string `xorm:"hhr_pygroup_id"`
		UpdateCalId string `xorm:"hhr_upt_calid"`
	}
	type Test struct {
		TenantId int64
		EmpId    string
		EmpRcd   int64
	}

	//engine := db.OrmEngine("hhr_payroll")
	//var valueSlc = make([]map[string]interface{}, 0)
	//err := engine.Table("hhr_py_cal_emp_add").Where("tenant_id=? and hhr_py_cal_id = 'T'", 0).Find(&valueSlc)
	//fmt.Println("valueSlc:", valueSlc)
	//fmt.Println("valueSlc[0][hhr_empid]:", string(valueSlc[0]["hhr_empid"].([]uint8)))

	/*test := make([]*Test, 0)
	err := engine.SQL("select hhr_empid emp_id, hhr_emp_rcd from hhr_py_cal_emp_add where tenant_id=? and hhr_py_cal_id = 'T'", 0).Find(&test)
	if err != nil {
		common.Logger.Error("[test]", err.Error())
	}
	fmt.Println("test:", test[0].EmpId)*/

	//calGrp := entity.NewCalendarGroup(0, "H01-2020-01")
	//fmt.Println("calGrp:", calGrp)
	/*type TempCatalog struct {
		TenantId   int64
		HhrEmpid   string
		HhrEmpRcd  int64
		HhrSeqNum  int64
		HhrHistSeq int64
		HhrUpdSeq  int64
		CalGrpId   string `xorm:"hhr_pycalgrp_id"`
		HhrFPrdNum int64
	}
	temp := &TempCatalog{}
	engine := db.OrmEngine("hhr_payroll")
	_, err := engine.Table("hhr_py_cal_catalog").Where("tenant_id=?", 0).Get(temp)
	if err != nil {
		common.Logger.Error("[test]", err.Error())
	}
	common.Logger.Debug("------test--------", temp)
	if temp.HhrFPrdNum == 0{
		fmt.Println("0000000000000000000000000")
	}*/
	test400()
}

func fibonacci(num int) int {
	if num < 2 {
		return 1
	}
	return fibonacci(num-1) + fibonacci(num-2)
}

func CurLocalDate() time.Time {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	t := time.Now()
	t = t.In(loc)
	return t
}

func test100() {
	engine := db.OrmEngine("hhr_payroll")
	stmt2 := "select A.hhr_empid from hhr_corehr.hhr_org_per_jobdata A where A.tenant_id = ? and A.hhr_efft_date > ? " +
		"and A.hhr_efft_date <= ? and A.hhr_status = 'Y' and A.hhr_empid = ? and A.hhr_emp_rcd = ? " +
		"union " +
		"select a.hhr_empid from hhr_corehr.hhr_org_per_jobdata a where a.tenant_id = ? and a.hhr_status = 'Y' " +
		"and ? between a.hhr_efft_date and a.hhr_efft_end_date " +
		"and a.hhr_efft_seq = (select max(a1.hhr_efft_seq) from hhr_corehr.hhr_org_per_jobdata a1 where a1.tenant_id = a.tenant_id " +
		"and a1.hhr_empid=a.hhr_empid and a1.hhr_emp_rcd = a.hhr_emp_rcd and a1.hhr_efft_date = a.hhr_efft_date) " +
		"and a.hhr_empid = ? and a.hhr_emp_rcd = ? "
	empIdRcdSlc := make([]*module.EmpIdRcd, 0)
	err := engine.SQL(stmt2, 0, "2010-10-01", "2020-05-17", "E00000004", 1, 0, "2010-10-01", "E00000004", 1).Find(&empIdRcdSlc)
	if err != nil {
		common.Logger.Error("[test100]", err.Error())
	}
	fmt.Println("empIdRcdSlc[0]:", empIdRcdSlc[0])

}

func test200() {
	engine := db.OrmEngine("hhr_payroll")
	furtherPayees := make([]*entity.Payee, 0)
	err := engine.Table("hhr_py_cal_catalog").Cols("hhr_empid", "hhr_emp_rcd").Distinct("hhr_empid", "hhr_emp_rcd").
		Where("tenant_id=? and hhr_f_prd_end_dt=?", 0, "2020-02-23").Find(&furtherPayees)
	if err != nil {
		common.Logger.Error("[beginIdentify-3]", err.Error())
	}
	fmt.Println("furtherPayees[0]:", furtherPayees[0])
}

func test300() {
	x := time.Time{}
	fmt.Println("x:", x)
	fmt.Println("x.IsZero():", x.IsZero())
}

func test400() {
	engine := db.OrmEngine("hhr_payroll")
	var ExceptPersList []*entity.Payee
	err := engine.Table("hhr_py_cal_emp_expt").Cols("hhr_empid", "hhr_emp_rcd").
		Where("tenant_id=? and hhr_py_cal_id=?", 0, "a").Find(&ExceptPersList)
	if err != nil {
		common.Logger.Error("[test400-3]", err.Error())
	}
	fmt.Println("ExceptPersList:", ExceptPersList)
	p := &ExceptPersList
	for _, i := range *p {
		fmt.Println("i:", i)
	}
}
