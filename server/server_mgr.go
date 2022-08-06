package server

import (
	"gopayroll/common"
	"gopayroll/db"
	"gopayroll/taskmgr"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"strconv"
	"strings"
	"time"
)

var worker *machinery.Worker

func loadConfig() (*config.Config, error) {
	if common.ConfigPath != "" {
		return config.NewFromYaml(common.ConfigPath, false)
	}

	return config.NewFromEnvironment(true)
}

func StartServer() (*machinery.Server, error) {
	cnf, err := loadConfig()
	if err != nil {
		return nil, err
	}
	// Create server instance
	common.SrvInst, err = machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}
	// Register payroll
	tasksMap := taskmgr.GetTasksMap()
	return common.SrvInst, common.SrvInst.RegisterTasks(tasksMap)
}

func LaunchWorker() error {
	consumerTag := "boogoo_worker"
	server, err := StartServer()
	if err != nil {
		return err
	}
	worker = server.NewWorker(consumerTag, common.Concurrency)
	//worker.SetErrorHandler(errorHandler)
	worker.SetPreTaskHandler(preTaskHandler)
	worker.SetPostTaskHandler(postTaskHandler)
	return worker.Launch()
}

/*func errorHandler(err error) {
	log.ERROR.Println("I am an error handler:", err.Error())
}*/

func preTaskHandler(signature *tasks.Signature) {
	taskState, _ := worker.GetServer().GetBackend().GetState(signature.UUID)
	tenantId := getTenantId(signature.UUID)
	engine := db.OrmEngine("hhr_foundation")
	loc, _ := time.LoadLocation(common.LocalLocation)
	hhrTaskState := taskmgr.HhrTaskState{
		TenantId:     tenantId,
		HhrTaskId:    taskState.TaskUUID,
		HhrTaskName:  taskState.TaskName,
		HhrState:     taskState.State,
		HhrTaskError: taskState.Error,
		HhrCreateAt:  taskState.CreatedAt.In(loc),
	}
	_, err := engine.InsertOne(&hhrTaskState)
	if err != nil{
		common.Logger.Error("[preTaskHandler]", err.Error())
	}
}

func postTaskHandler(signature *tasks.Signature) {
	taskState, _ := worker.GetServer().GetBackend().GetState(signature.UUID)
	tenantId := getTenantId(signature.UUID)
	err := updateTaskState(tenantId, signature.UUID, taskState.State, taskState.Error)
	if err != nil {
		common.Logger.Error("[postTaskHandler]", err.Error())
	}
}

func getTenantId(uuid string) int64 {
	lastIndex := strings.Index(uuid, ":")
	tenantIdStr := uuid[:lastIndex]
	t, err := strconv.Atoi(tenantIdStr)
	if err != nil {
		common.Logger.Error("[getTenantId]", "tenantId conversion error from string to int")
	}
	tenantId := int64(t)
	return tenantId
}

func updateTaskState(tenantId int64, taskUUID string, state string, errText string) error {
	engine := db.OrmEngine("hhr_foundation")
	_, err := engine.Prepare().Exec("UPDATE hhr_foundation.hhr_task_state set hhr_state = ?, hhr_task_error = ?, hhr_end_date = ? " +
		"WHERE tenant_id = ? and hhr_task_id = ?", state, errText, common.CurLocalDate(), tenantId, taskUUID)
	if err != nil {
		common.Logger.Error("[updateTaskState]", err.Error())
		return err
	}
	return nil
}
