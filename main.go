package main

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
)

type Demo struct {
	Name string
	Age int64

}

type Variable struct{
	Name string
	Value float64
}

var variable = Variable{}

var vm = otto.New()
func main(){

	vm.Run(`
		abc = 2 + 2;
		console.log("The value of abc is " + abc); // 4
	`)

	value, _ := vm.Get("abc")
	v1, _ := value.ToInteger()
	fmt.Println("v1:", v1)

	vm.Set("def", 11)
	vm.Run(`
	console.log("The value of def is " + def);
	// The value of def is 11
	`)

	vm.Set("xyzzy", "Nothing happens.")
	vm.Run(`
	console.log(xyzzy.length); // 16
	`)

	value, _ = vm.Run("xyzzy.length")
	{
		// value is an int64 with a value of 16
		v2, _ := value.ToInteger()
		fmt.Println("v2:" ,v2)
	}

	value, err1 := vm.Run("abcdefghijlmnopqrstuvwxyz.length")
	if err1 != nil {
		// err = ReferenceError: abcdefghijlmnopqrstuvwxyz is not defined
		// If there is an error, then value.IsUndefined() is true
		fmt.Println("value.IsUndefined():",  value.IsUndefined())
	}

	vm.Set("sayHello", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Hello, %s.\n", call.Argument(0).String())
		return otto.Value{}
	})

	vm.Set("twoPlus", func(call otto.FunctionCall) otto.Value {
		right, _ := call.Argument(0).ToInteger()
		result, _ := vm.ToValue(2 + right)

		return result
	})

	vm.Set("test", func(call otto.FunctionCall) otto.Value {
		a1, _ := call.Argument(0).ToInteger()
		a2, _ := call.Argument(1).ToInteger()
		right:=a1+a2
		result, _ := vm.ToValue(2 + right)
		return result
	})

	result, _ := vm.Run(`
    sayHello("Xyzzy");      // Hello, Xyzzy.
    sayHello();             // Hello, undefined
	

    result = twoPlus(2.0); // 4
 	resutl = test(1, 3);
	`)
	fmt.Println("result:", result)

	// value is a String object
	//value1, _ := vm.Call("Object", nil, "Hello, World.")

	// Likewise...
	//value1, _ := vm.Call("new Object", nil, "Hello, World.")

	// This will perform a concat on the given array and return the result
	// value is [ 1, 2, 3, undefined, 4, 5, 6, 7, "abc" ]
	value1, _ := vm.Call(`[ 1, 2, 3, undefined, 4 ].concat`, nil, 5, 6, 7, "abc")
	value11 , _ := value1.ToString()
	fmt.Println("value11:", value11)

	//校验脚本
	script, err := vm.Compile("", `var abc; if (!abc) abc = 0; abc += 2; abc;`)
	if err != nil {
		fmt.Println("出错了:", err.Error())
		return
	}
	vm.Run(script)

	//调用JavaScript函数
	s1 := `var x = function (a, b) {return a * b};`
	vm.Run(s1)
	value2 , _ := vm.Call("x",nil,2,2)
	fmt.Println("value2:", value2)


	// 将执行的结果转换为Golang对应的类型
	r, _ := vm.Run(`var x = function (a, b) {return a * b};x(2,3);`)
	v, _ := r.Export()
	switch v.(type) {
	case float64:
		fmt.Println("haha")
		fmt.Println(v.(float64))
	}

	var objJSON = `{
	  "bear": [
		"foo",
		"fooo"
	  ]
	}`

	var obj interface{}
	json.Unmarshal([]byte(objJSON), &obj)
	fmt.Println("obj:", obj)
/*
	deObj, _ := vm.ToValue(obj)
	vm.Set("deObj", deObj)
	xyz, err := vm.Run(`deObj.Xyz`)
	fmt.Println("xyz:", xyz)*/

	//传递对象
	de := Demo{Name:"david", Age: 32}
	deObj, _ := vm.ToValue(de)
	vm.Set("deObj", deObj)
	xyz, err := vm.Run(`deObj.Name`)
	fmt.Println("xyz:", xyz)

	vm.Set("test", myTest)
	vm.Run(`var a = test('张三', 23);  
            var b = a.Name;
 			a.Name='李四'
			console.log("The value of a.Name is " + a.Name);`)
	bb, _ := vm.Get("a")
	fmt.Println("bb:", bb)
	bbObj := bb.Object()
	bbValue, _:= bbObj.Get("Name")
	fmt.Println("bbValue:", bbValue)

	vm.Set("VR", VR)
	vm.Run(`VR('VR_DAYS_ALL').Value = 10 
VR('VR_DAYS_ALL').Value = 20
VR('VR_DAYS_ALL').Value = VR('VR_DAYS_ALL').Value + 88.3`)
	//vm.Run(`vv = VR('VR_DAYS_ALL'); vv.Value=22.33`)
	//vv, _ := vm.Get("vv")
	//vv1,_:= vv.Export()
	fmt.Println("vv:", variable)
}

func VR(call otto.FunctionCall) otto.Value {
	vrName, _ := call.Argument(0).ToString()
	variable.Name = vrName
	result, _ := vm.ToValue(&variable)
	return result
}

func myTest(call otto.FunctionCall) otto.Value {
	name, _ := call.Argument(0).ToString()
	age, _ := call.Argument(1).ToInteger()

	de := Demo{}
	de.Name = name
	de.Age = age

	result, _ := vm.ToValue(&de)
	return result
}