package module

type PayCalcRunParam struct {
	UserId          string `json:"user_id"`
	TenantId        int64  `json:"tenant_id"`
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
