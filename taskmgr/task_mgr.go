package taskmgr

import (
	"gopayroll/common"
	"gopayroll/taskspool"
	"gopayroll/taskspool/payroll"
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func SendTask(task *tasks.Signature) error {
	common.Logger.Debug("[SendTask] -----Send task---->>")
	_, err := common.SrvInst.SendTask(task)
	if err != nil {
		common.Logger.Error("[SendTask]", err.Error())
		return fmt.Errorf("could not send task: %s", err.Error())
	}
	return nil
}

func SendChainTasks(chain *tasks.Chain) error {
	_, err := common.SrvInst.SendChain(chain)
	if err != nil {
		return err
	}
	return nil
}

func GetTasksMap() map[string]interface{} {
	tasksMap := map[string]interface{}{
		"pyIdentity": payroll.PyIdentity,
		"pyCalc":     payroll.PyCalc,
		"add":        taskspool.Add,
	}
	return tasksMap
}