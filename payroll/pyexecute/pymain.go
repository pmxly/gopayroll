/*
Desc: 薪资计算主入口
Author: 潘承勋
Date: 2020-05-03
*/

package pyexecute

import (
	"gopayroll/common"
	"gopayroll/db"
	pyCommon "gopayroll/payroll/common"
	"gopayroll/payroll/entity"
	"gopayroll/payroll/module"
	"gopayroll/payroll/pyexecute/identify"
	"gopayroll/payroll/pysysutils"
	"gopayroll/utils"
)

//分布式薪资计算入口函数
func RunDistPyEngine(runParam *module.PayCalcRunParam) error {
	//初始化日历组对象
	tenantId, calGrpId := runParam.TenantId, runParam.CalGrpID

	//初始化(租户+日历组)对应的缓存
	globalKey := pysysutils.GetGlobalKey(tenantId, calGrpId)
	pyCommon.InitGlobalRunVarCache(globalKey)
	pyCommon.InitGlobalVarCache(globalKey)

	calGrp := entity.NewCalendarGroup(tenantId, calGrpId)
	runVarCache := pyCommon.RunVarCache(globalKey)
	runVarCache.SetRunVarObj("CAL_GRP_OBJ", calGrp)

	entity.InitVariables(tenantId, calGrp)

	runVarCache.SetRunVarObj("TENANT_ID", tenantId)

	//获取控制开关值
	switchMap, err := utils.GetSwitchMap(tenantId, "PY_AUTOCAL_ENTRY", "PY_AUTOCAL_PERIOD")
	if err != nil {
		return err
	}
	//common.Logger.WithFields(logrus.Fields{"switchmap": switchMap,}).Info()
	pyAutoCalEntry, _ := switchMap["PY_AUTOCAL_ENTRY"]
	pyAutoCalPeriod, _ := switchMap["PY_AUTOCAL_PERIOD"]
	runVarCache.SetRunVarObj("VR_AUTOCAL_ENTRY", pyAutoCalEntry)
	runVarCache.SetRunVarObj("VR_AUTOCAL_PERIOD", pyAutoCalPeriod)

	runVarCache.SetRunVarObj("COMM_LOG_FLAG", runParam.LogFlag)
	//将日历组下的所有日历放进run变量
	runVarCache.SetRunVarObj("CAL_OBJ_DIC", calGrp.GetCalendarMap(tenantId, calGrpId))

	switch runParam.TaskName {
	case pyCommon.PyIdentify:
		if err = IdentifyPayees(globalKey, runParam); err != nil {
			return err
		}
	case pyCommon.PyCalc:
	}
	return nil
}

//标记受款人
func IdentifyPayees(globalKey string, runParam *module.PayCalcRunParam) error {
	engine := db.OrmEngine("hhr_payroll")
	session := engine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Logger.Error("[IdentifyPayees]", err.Error())
	}
	pi := identify.NewPayeesIdentify(globalKey, runParam)
	err = pi.IdentifyPayees()
	if err != nil {
		common.Logger.Error("[IdentifyPayees]error occurred when identifying payees:", err.Error())
		_ = session.Rollback()
		return err
	}
	return session.Commit()
}