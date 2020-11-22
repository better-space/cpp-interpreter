package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

type FuncBody struct {
	FuncName string 			  //函数名
	//ReValue  *Variable          //返回值
	PipLine  []string             //存储函数体内所有语句
	VarMap   map[string]*Variable //存储函数局部变量
	ExpStack []InterPara          //模拟栈结构，用于表达式计算，存储临时变量
	SymStack []string             //模拟栈结构，用于表达式计算，存储运算符
}

type Variable struct {
	typ string		//类型
	val float64		//值
}

type InterPara struct {
	name  string	//变量名
	value float64	//变量值
}

//赋值语句处理
func (f *FuncBody) IsAssignment(i int) {
	strings.Split(f.PipLine[i], " ")
	s := strings.FieldsFunc(f.PipLine[i], func(r rune) bool { //按多个不同字符分割
		return r == ' '
	})
	val, _ := IsNum(s[3])
	f.VarMap[s[1]] = &Variable{
		typ: s[0],
		val: val,
	}
	Num++
	fmt.Fprintf(Outfile, "<%s,σ%v>→%v\n\n", s[1], Num, val)
}

//表达式语句处理
func (f *FuncBody) IsExpression(i int) {
	Num++
	str := strings.FieldsFunc(f.PipLine[i], func(r rune) bool { //按多个不同字符分割
		return r == ' '
	})
	sl := len(str)
	//分割表达式
	s := SplitExpression(str[sl-1])
	//进行*/%的运算
	for i := 0; i < len(s); i++ {
		v := s[i]
		switch v {
		case "":
		case "+":
			f.SymStack = append(f.SymStack, v)
		case "-":
			f.SymStack = append(f.SymStack, v)
		case "*":
			top := f.ExpStack[len(f.ExpStack)-1]
			f.ExpStack = f.ExpStack[:len(f.ExpStack)-1]
			next := s[i+1]
			var val float64
			if num, ok := IsNum(next); ok {
				val = num
			} else {
				val = f.VarMap[next].val
			}
			in := InterPara{
				name:  top.name + "*" + next,
				value: top.value * val,
			}
			f.ExpStack = append(f.ExpStack, in)
			fmt.Fprintf(Outfile, "<%s*%s,σ%v>→%v\n\n", top.name, next, Num, in.value)
		case "/":
			top := f.ExpStack[len(f.ExpStack)-1]
			f.ExpStack = f.ExpStack[:len(f.ExpStack)-1]
			next := s[i+1]
			var val float64
			if num, ok := IsNum(next); ok {
				val = num
			} else {
				val = f.VarMap[next].val
			}
			in := InterPara{
				name:  top.name + "/" + next,
				value: top.value / val,
			}
			f.ExpStack = append(f.ExpStack, in)
			fmt.Fprintf(Outfile, "<%s/%s,σ%v>→%v\n\n", top.name, next, Num, in.value)
		case "(":
			f.SymStack = append(f.SymStack, v)
		case ")":
			var val float64
			top0 := f.ExpStack[len(f.ExpStack)-1]
			top1 := f.ExpStack[len(f.ExpStack)-2]
			f.ExpStack = f.ExpStack[:len(f.ExpStack)-2]
			symbol0 := f.SymStack[len(f.SymStack)-1]
			symbol1 := f.SymStack[len(f.SymStack)-2]
			f.SymStack = f.SymStack[:len(f.SymStack)-2]
			if symbol0 == "+" {
				val = top0.value + top1.value
			} else {
				val = top1.value - top0.value
			}
			in := InterPara{
				name:  symbol1 + top1.name + symbol0 + top0.name + v,
				value: val,
			}
			f.ExpStack = append(f.ExpStack, in)
			fmt.Fprintf(Outfile, "<%s%s%s%s%s,σ%v>→%v\n\n", symbol1, top1.name, symbol0, top0.name, v, Num, in.value)
		case "%":
			top := f.ExpStack[len(f.ExpStack)-1]
			f.ExpStack = f.ExpStack[:len(f.ExpStack)-1]
			next := s[i+1]
			var val float64
			if num, ok := IsNum(next); ok {
				val = num
			} else {
				val = f.VarMap[next].val
			}
			in := InterPara{
				name:  top.name + "%" + next,
				value: float64(int(top.value) % int(val)),
			}
			f.ExpStack = append(f.ExpStack, in)
			fmt.Fprintf(Outfile, "<%s%%%s,σ%v>→%v\n\n", top.name, next, Num, in.value)
		default:
			var in InterPara

			if val, ok := IsNum(v); ok {
				in = InterPara{
					name:  v,
					value: val,
				}
			} else {
				in = InterPara{
					name:  v,
					value: f.VarMap[v].val,
				}
			}

			f.ExpStack = append(f.ExpStack, in)
		}
	}

	//进行加减运算
	//计算栈中的元素，为了输出顺序从前往后，此处栈转队列用
	var symFront string
	var first, second InterPara
	var ans float64
	for len(f.SymStack) != 0 {
		symFront = f.SymStack[0]
		f.SymStack = f.SymStack[1:]
		first, second = f.ExpStack[0], f.ExpStack[1]
		if symFront == "-" {
			ans = first.value - second.value
		} else {
			ans = first.value + second.value
		}
		fmt.Fprintf(Outfile, "<%s%s%s,σ%v>→%v\n\n", first.name, symFront, second.name, Num, ans)
		//将每次运算得到的值插入队列首部
		f.ExpStack[1].name = first.name + symFront + second.name
		f.ExpStack[1].value = ans
		f.ExpStack = f.ExpStack[1:]
	}

	fmt.Fprintf(Outfile, "<%s,σ%v>→%v\n\n", str[sl-3], Num, f.ExpStack[0].value)

	if f.VarMap[str[sl-3]] == nil {
		f.VarMap[str[sl-3]] = new(Variable)
	}
	f.VarMap[str[sl-3]].val = f.ExpStack[0].value
	f.ExpStack = f.ExpStack[:0] //运算时的中间变量栈清零
}

//选择结构处理
func (f *FuncBody) IsSelectStructure(line string) {

}

//循环结构处理
func (f *FuncBody) IsCycle(line string) {

}

func (f *FuncBody) GetNextValue(s string) interface{} {
	if s[0] >= '0' && s[0] <= '9' {
		val, _ := strconv.Atoi(s)
		return val
	} else {
		return f.VarMap[s]
	}
}
