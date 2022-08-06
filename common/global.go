package common

import (
	"github.com/RichardKnop/machinery/v1"
	"time"
)
var (
	ConfigPath string = "config.yml"
	SrvInst *machinery.Server
	CSTZone = time.FixedZone("CST", 8*3600)
	ZeroTime = time.Time{}
)

type BasicParam struct {
	TenantId    	int    `json:"tenant_id"`
	UserId      	string `json:"user_id"`
	SecretToken 	string `json:"secret_token"`
}

type PyCalParam struct {
	TenantId    	int64  `json:"tenant_id"`
	UserId      	string `json:"user_id"`
	CalId 			string `json:"cal_id"`
	Action			string `json:"action"`
	LogFlag			string `json:"log_flag"`
	LogEmpRecStr	string `json:"log_emp_rec_str"`
}