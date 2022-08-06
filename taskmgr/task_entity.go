package taskmgr

import "time"

type HhrTaskState struct {
	ID int64 `xorm:"pk autoincr 'id'"`
	TenantId int64 `xorm:"tenant_id"`
	HhrTaskId string `xorm:"varchar(40) 'hhr_task_id'"`
	HhrTaskName string `xorm:"varchar(100) 'hhr_task_name'"`
	HhrState string `xorm:"varchar(20) 'hhr_state'"`
	HhrTaskError string `xorm:"hhr_task_error"`
	HhrCreateAt time.Time `xorm:"hhr_create_at"`
	HhrEndDate time.Time `xorm:"hhr_end_date"`
}