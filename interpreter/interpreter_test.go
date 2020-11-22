package interpreter

import (
	"bytes"
	"fmt"
	"testing"
)



func TestInterpreter(t *testing.T)  {
	var funcList []*FuncBody
	var ok bool

	data := LoadFile("input.txt")
	for _,v := range data {
		fmt.Println(string(v))
	}
	if funcList, ok = DealData(data); ok {
		fmt.Println("deal data error.")
	}
	//找到程序入口-main函数
	f := &FuncBody{}
	//找到程序入口main函数
	for _,v := range funcList {
		if bytes.Compare(S2B(v.FuncName), S2B("main")) == 0 {
			f = v
			break
		}
	}
	Simulation(f, funcList)
	fmt.Println(funcList)
}
