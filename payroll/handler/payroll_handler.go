package handler

import (
	"gopayroll/common"
	pyCommon "gopayroll/payroll/common"
	"gopayroll/payroll/module"
	"gopayroll/payroll/pyexecute"
	"gopayroll/taskmgr"
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"net/http"
)

func PayrollCalc1(ctx *gin.Context) {
	var pyCalParam common.PyCalParam
	if err := ctx.ShouldBindBodyWith(&pyCalParam, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenantId := pyCalParam.TenantId
	var (
		identifyFlag, reIdentFlag, pyCalcFlag, cancelIdentFlag,
		reCalcFlag, reIdentCalcFlag, finishFlag, cancelCalcFlag, cancelFlag, closeFlag string
	)
	pyAction := pyCalParam.Action
	switch pyAction {
	case "IDENTIFY":
		identifyFlag = "Y"
	case "RE_IDEN":
		reIdentFlag = "Y"
	case "PY_CALC":
		pyCalcFlag = "Y"
	case "CANCEL_IDEN":
		cancelIdentFlag = "Y"
	case "RE_CALC":
		reCalcFlag = "Y"
	case "RE_IDEN_CALC":
		reIdentCalcFlag = "Y"
	case "FINISH":
		finishFlag = "Y"
	case "CANCEL_CALC":
		cancelCalcFlag = "Y"
	case "CANCEL_BOTH":
		cancelFlag = "Y"
	case "CLOSE":
		closeFlag = "Y"
	default:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Action is not supported"})
	}

	runParam := module.PayCalcRunParam{
		UserId:          pyCalParam.UserId,
		TenantId:        tenantId,
		CalGrpID:        pyCalParam.CalId,
		CalID:           pyCalParam.CalId,
		IdentityFlag:    identifyFlag,
		ReIdentityFlag:  reIdentFlag,
		PyCalcFlag:      pyCalcFlag,
		ReCalcFlag:      reCalcFlag,
		ReIdentCalcFlag: reIdentCalcFlag,
		CancelIdentFlag: cancelIdentFlag,
		CancelCalcFlag:  cancelCalcFlag,
		CancelFlag:      cancelFlag,
		FinishFlag:      finishFlag,
		CloseFlag:       closeFlag,
		LogFlag:         pyCalParam.LogFlag,
		LogEmpRecStr:    pyCalParam.LogEmpRecStr}

	//重新标记&计算标志
	if reIdentCalcFlag == "Y" {
		runParam.ReIdentityFlag = "Y"
		runParam.ReCalcFlag = "Y"
	}
	//标记人员 或 重新标记
	if identifyFlag == "Y" || reIdentFlag == "Y" {
		runParam.TaskName = pyCommon.PyIdentify
		jsonRunParam, _ := json.Marshal(runParam)
		pyIdentTask := &tasks.Signature{
			Name: "pyIdentity",
			Args: []tasks.Arg{
				{
					Type:  "[]byte",
					Value: jsonRunParam,
				},
			},
		}

		taskID0 := uuid.New().String()
		pyIdentTask.UUID = fmt.Sprintf("%v:task_%v", tenantId, taskID0)
		err := taskmgr.SendTask(pyIdentTask)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "send task error"})
		}

	}

	/*	taskID1 := uuid.New().String()
		taskspool.TaskPool.AddTask1.UUID = fmt.Sprintf("task_%v", taskID1)
		fmt.Println("taskID1:", taskID1)

		chain, _ := payroll.NewChain(&taskspool.TaskPool.AddTask0, &taskspool.TaskPool.AddTask1)
		taskmgr.SendChainTasks(chain)*/
	ctx.String(http.StatusOK, "你好")
}

func PayrollCalc(ctx *gin.Context) {
	var pyCalParam common.PyCalParam
	if err := ctx.ShouldBindBodyWith(&pyCalParam, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var runParam = &module.PayCalcRunParam{
		UserId:   "u101",
		TenantId: 0,
		CalGrpID: "H01-2020-01",
		TaskName: pyCommon.PyIdentify,
	}
	pyexecute.RunDistPyEngine(runParam)

	ctx.String(http.StatusOK, "success")
}
