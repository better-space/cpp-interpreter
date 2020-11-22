package main

import (
	"bytes"
	"fmt"
	"interpreter/interpreter"
)

func main() {
	var funcList []*interpreter.FuncBody
	var ok bool
	data := interpreter.LoadFile("./input.txt")
	for _, v := range data {
		fmt.Println(string(v))
	}
	if funcList, ok = interpreter.DealData(data); ok {
		fmt.Println("deal data error.")
	}
	//找到程序入口-main函数
	f := &interpreter.FuncBody{}
	//找到程序入口main函数
	for _, v := range funcList {
		if bytes.Compare(interpreter.S2B(v.FuncName), interpreter.S2B("main")) == 0 {
			f = v
			break
		}
	}
	interpreter.Simulation(f, funcList)
	fmt.Println(funcList)
}
