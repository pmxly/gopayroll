package routes

import (
	"gopayroll/middleware"
	"gopayroll/payroll/handler"
	"github.com/gin-gonic/gin"
)

func RouteManager(engine *gin.Engine) {
	rg := engine.Group("/v1")
	rg.Use(middleware.AuthMiddleware())
	{
		//薪资计算
		rg.POST("/payroll/calc", handler.PayrollCalc)
	}
}

