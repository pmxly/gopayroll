package common

//薪资计算过程中的变量对象
type RunVariableObject struct {
	VarID string
	Value interface{}
}

type VariableObject struct {
	TenantId int64
	VarId string `xorm:"hhr_variable_id"`
	DataType string `xorm:"hhr_data_type"`
	VarType string `xorm:"hhr_var_type"`
	Country string `xorm:"hhr_country"`
	CoverEnable string `xorm:"hhr_cover_enable"`
	HasCovered string `xorm:"-"`
	Value interface{} `xorm:"-"`
}

type RunVarMap map[string]*RunVariableObject
type VarMap map[string]*VariableObject

var GlobalRunVarCache = make(map[string]RunVarMap)
var GlobalVarCache = make(map[string]VarMap)

func InitGlobalRunVarCache(globalKey string) RunVarMap {
	GlobalRunVarCache[globalKey] = make(RunVarMap)
	return GlobalRunVarCache[globalKey]
}

func InitGlobalVarCache(globalKey string) VarMap {
	GlobalVarCache[globalKey] = make(VarMap)
	return GlobalVarCache[globalKey]
}

func RunVarCache(globalKey string) RunVarMap {
	return GlobalRunVarCache[globalKey]
}

func VarCache(globalKey string) VarMap {
	return GlobalVarCache[globalKey]
}

func (runVarMap RunVarMap) SetRunVarObj(varID string, value interface{}) {
	runVarMap[varID] = &RunVariableObject{VarID: varID, Value: value}
}

func (runVarMap RunVarMap) GetRunVarValue(varId string) interface{} {
	runVarObj, ok := runVarMap[varId]
	if !ok {
		return nil
	}
	return runVarObj.Value
}

func (varMap VarMap) SetVarObj(key string, value *VariableObject) {
	varMap[key] = value
}

func (varMap VarMap) GetVarValue(varId string) interface{} {
	varObj, ok := varMap[varId]
	if !ok {
		return nil
	}
	return varObj.Value
}
