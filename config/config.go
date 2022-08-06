package config

import (
	"gopayroll/common"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

var reloadDelay = time.Second * 10

type Config struct {
	DataSource *DataSource `yaml:"data_source"`
}

type DataSource struct {
	DriverName  string   `yaml:"driver_name"`
	DBUserName  string   `yaml:"db_username"`
	DBPassword  string   `yaml:"db_password"`
	DBHost      string   `yaml:"db_host"`
	DBPort      int      `yaml:"db_port"`
	DBSchemas   []string `yaml:"db_schemas"`
	MaxOpenConn int      `yaml:"max_open_conn"`
	MaxIdleConn int      `yaml:"max_idle_conn"`
	ShowSql     bool     `yaml:"show_sql"`
}

func LoadConfig() (*Config, error) {
	return NewFromYaml(common.ConfigPath, false)
}

func NewFromYaml(cnfPath string, keepReloading bool) (*Config, error) {
	cnf, err := fromFile(cnfPath)
	if err != nil {
		return nil, err
	}
	common.Logger.WithFields(logrus.Fields{"filepath": cnfPath}).Debug("[NewFromYaml] Successfully loaded config from file")
	if keepReloading {
		// Open a goroutine to watch remote changes forever
		go func() {
			for {
				// Delay after each request
				time.Sleep(reloadDelay)

				// Attempt to reload the config
				newCnf, newErr := fromFile(cnfPath)
				if newErr != nil {
					common.Logger.WithFields(logrus.Fields{"filepath": cnfPath}).Warn("[NewFromYaml] Failed to reload config from file ", newErr.Error())
					continue
				}

				*cnf = *newCnf
				common.Logger.WithFields(logrus.Fields{"filepath": cnfPath}).Info("[NewFromYaml] Successfully reloaded config from file")
			}
		}()
	}

	return cnf, nil
}

// ReadFromFile reads data from a file
func ReadFromFile(cnfPath string) ([]byte, error) {
	file, err := os.Open(cnfPath)

	// Config file not found
	if err != nil {
		return nil, fmt.Errorf("Open file error: %s", err)
	}

	// Config file found, try to read it
	data := make([]byte, 1000)
	count, err := file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("Read from file error: %s", err)
	}

	return data[:count], nil
}

func fromFile(cnfPath string) (*Config, error) {
	cnf := new(Config)
	data, err := ReadFromFile(cnfPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cnf); err != nil {
		return nil, fmt.Errorf("Unmarshal YAML error: %s", err)
	}
	return cnf, nil
}
