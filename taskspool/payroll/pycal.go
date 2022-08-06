package payroll

import (
	"gopayroll/common"
	"gopayroll/payroll/module"
	"gopayroll/payroll/pyexecute"
	"encoding/json"
)

//标记受款人
func PyIdentity(jsonRunParam []byte) error {
	var runParam = &module.PayCalcRunParam{}
	err := json.Unmarshal(jsonRunParam, runParam)
	common.Logger.Info("======runParam=====", runParam)
	if err != nil {
		common.Logger.Error("[PyIdentity]", err.Error())
		return err
	}
	return pyexecute.RunDistPyEngine(runParam)
}

//薪资计算
func PyCalc(runParam []byte) error {

	return nil
}