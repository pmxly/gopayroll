package db

import (
	"gopayroll/common"
	"gopayroll/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
)

var (
	allOrmEngines = make(map[string]*xorm.Engine)
)

func InitDBEngine() {
	cnf, _ := config.LoadConfig()
	for _, schema := range cnf.DataSource.DBSchemas {
		dataSrcName := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=true&loc=%s", cnf.DataSource.DBUserName, cnf.DataSource.DBPassword,
			cnf.DataSource.DBHost, cnf.DataSource.DBPort, schema, common.LocalEscLoc)
		engine, err := xorm.NewEngine(cnf.DataSource.DriverName, dataSrcName)
		if err != nil {
			common.Logger.Error("[InitDBEngine]", err.Error())
		}
		engine.ShowSQL(cnf.DataSource.ShowSql)
		engine.SetMaxOpenConns(cnf.DataSource.MaxOpenConn)
		engine.SetMaxIdleConns(cnf.DataSource.MaxIdleConn)
		setEngine(schema, engine)
	}
}

func setEngine(key string, e *xorm.Engine) {
	allOrmEngines[key] = e
}

func OrmEngine(key string) *xorm.Engine {
	if key != "" {
		value, ok := allOrmEngines[key]
		if !ok {
			common.Logger.WithFields(logrus.Fields{"key" : key}).Error("[OrmEngine] table schema does not exists in engine map")
		}
		return value
	}else{
		common.Logger.Error("[OrmEngine] map key should not be empty")
		return nil
	}
}