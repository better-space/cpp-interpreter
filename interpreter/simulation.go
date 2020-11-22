package interpreter

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

func Simulation(f *FuncBody, funcList []*FuncBody)  {

	for i:=0;i<len(f.PipLine); {
		if regularFor(S2B(f.PipLine[i])) {
			log.Println("this is a for cycle")
			elem := strings.FieldsFunc(f.PipLine[i], func(r rune) bool {
				return r == ';' || r == '(' || r == ')'
			}) //根据;（）分割字符串
			_, low, high := ForJudgeLowHigh(elem[1:4])
			index := low
			for n := i + 1; ; n++ {
				if f.PipLine[n] == "}" {
					index++
					if index > high {
						i=n
						break
					} else {
						n = i+1
					}
				}
				f.IsExpression(n)
			}
		} else if regularIf(S2B(f.PipLine[i])) {	//if else 条件语句
			elem := strings.Split(f.PipLine[i], " ")	//分号被去掉
			comparer := elem[1]
			if JudgeCompareExpression(f, comparer) {
				flag := false
				i++
				for ;f.PipLine[i] != "}";i++ {
					if f.PipLine[i][0] == '}' || flag{
						flag = true
						continue
					}
					f.IsExpression(i)
				}
			} else {
				flag := false
				i++
				for ;f.PipLine[i] != "}";i++ {
					if f.PipLine[i][0] == '}'{
						flag = true
						continue
					}
					if flag {
						f.IsExpression(i)
					}
				}
			}
		} else if f.PipLine[i][len(f.PipLine[i])-1] == ')'{	//匹配函数调用
			var sf *FuncBody
			funLine := strings.FieldsFunc(f.PipLine[i], func(r rune) bool {
				return r == '(' || r == ')' || r == ' '
			})
			for _,v := range funcList {
				if bytes.Compare(S2B(v.FuncName),S2B(funLine[0])) == 0 {
					sf = v
					break
				}
			}
			Simulation(sf, funcList)
		} else if regularWhile(S2B(f.PipLine[i])) {
			fmt.Println("this is a while")
			log.Println("this is a for cycle")
			elem := strings.FieldsFunc(f.PipLine[i], func(r rune) bool {
				return r == ';' || r == '(' || r == ')'
			}) //根据;（）分割字符串
			name, low, high := WhileJudgeLowHigh(elem[1],f)
			for n := i + 1; ; n++ {
				if f.PipLine[n] == "}" {
					if int(f.VarMap[name].val) > high || int(f.VarMap[name].val) < low{
						i=n
						break
					} else {
						n = i+1
					}
				}
				f.IsExpression(n)
			}
		} else if regularExpression(S2B(f.PipLine[i])) {	//运算式
			log.Println("this is expression")
			f.IsExpression(i)
		} else if regularAssignment(S2B(f.PipLine[i])) { //赋值表达式
			f.IsAssignment(i)
		}
		i++
 	}
}
