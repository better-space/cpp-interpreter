package interpreter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const Output string = "./output.txt"
var Outfile = GetOutputFile()
var Num int

func GetOutputFile() *os.File {
	output, err := os.OpenFile(Output, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		log.Println("open output file error.")
	}
	return output
}

func LoadFile(filename string) [][]byte {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		fmt.Println("open input.txt error")
	}
	defer file.Close()
	//out,err := os.OpenFile(output, os.O_WRONLY, 0666)
	//if err != nil {
	//	fmt.Println( "open output.txt error")
	//}

	data := make([][]byte, 0)
	r := bufio.NewReader(file)
	for {
		content, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil
			}
		}
		line := bytes.TrimSpace(content)
		if len(line) == 0 || len(line) >= 2 && bytes.Compare(line[:2], []byte(`//`)) == 0 {
			continue
		}
		data = append(data, CopyBytes(line))
	}
	return data
}

func CopyBytes(b []byte) []byte {
	t := make([]byte, len(b))
	for i, v := range b {
		t[i] = v
	}
	return t
}

//正则匹配是否是函数行
func regular(s []byte) bool {
	r, err := regexp.Compile(`[^a-zA-Z_](\S*?)(\s*?)[a-zA-Z_](\S*?)(\s*?)\(.*?\)(.*?)(\{|;)`)
	if err != nil {
		fmt.Println("error when compile string")
	}
	return r.Match(s)
}

//正则匹配函数行的函数名
func regularString(s []byte) string {
	r := regexp.MustCompile(` (.*\()`)
	return r.FindString(string(s))
}

//正则匹配赋值表达式
func regularAssignment(s []byte) bool {
	r,err := regexp.Compile(`((\d+(\.(\d+))?)((\+|\-)(\d+(\.(\d+))?))*)`)
	if err != nil {
		fmt.Println("error when compile string")
	}
	return r.Match(s)

	//非正则匹配
	//list := strings.Split(string(s), " ")
	//if len(list) == 3 && list[2][0] >= '0' && list[2][0] <= '9' {
	//	return true
	//}
	//return false
}

//正则匹配运算表达式
func regularExpression(s []byte) bool {
	r,err := regexp.Compile(`[\+\-\*/%]|= [A-Za-z]`)
	if err != nil {
		fmt.Println("error when compile string")
	}
	return r.Match(s)
}

func regularFuncNameInFunc(s []byte) bool {
	r,err := regexp.Compile(`(static\s*){0,1}\w{1,}\s{1,}\w{1,}\s*\(.*\)[^;]`)
	if err != nil {
		fmt.Println("error when compile string")
	}
	return r.Match(s)
}

//正则匹配if语句
func regularIf(s []byte) bool {
	r, err := regexp.Compile(`if\s*(.*?)`)
	if err != nil {
		fmt.Println("error when compile if")
	}
	return r.Match(s)
}

//正则匹配for循环
func regularFor(s []byte) bool {
	r, err := regexp.Compile(`for(.*)(\()(.*)(;)(.*)(\s*)(.*)(;)(\s*)(.*{)`)
	if err != nil {
		fmt.Println("error when compile if")
	}
	return r.Match(s)
}

//正则匹配while循环
func regularWhile(s []byte) bool {
	r, err := regexp.Compile(`while(.*)`)
	if err != nil {
		fmt.Println("error when compile if")
	}
	return r.Match(s)
}

func DealData(data [][]byte) ([]*FuncBody, bool) {
	//返回体
	var funcList []*FuncBody

	for _, v := range data {
		l := len(funcList)
		if regular(v) {
			if name := regularString(v); len(name) != 0 {
				f := &FuncBody{
					FuncName: name[1 : len(name)-1],
					//ReValue:  new(Variable),
					PipLine:  make([]string, 0),
					VarMap:   make(map[string]*Variable),
					ExpStack: make([]InterPara, 0),
					SymStack: make([]string, 0),
				}
				funcList = append(funcList, f)
			}
		} else if bytes.Compare(v, []byte("}")) == 0 || v[len(v)-1] == '{'{
			funcList[l-1].PipLine = append(funcList[l-1].PipLine, string(v))
		} else if bytes.Compare(v[:2], []byte("//")) == 0 || len(v) == 0{
			continue
		} else {
			funcList[l-1].PipLine = append(funcList[l-1].PipLine, string(v[:len(v)-1]))
			fmt.Printf("%#v", funcList)
		}
	}
	return funcList, true
}

//字符串转字符数组
func S2B(s string) []byte {
	return []byte(s)
}

//判断字符串是否是数字
func IsNum(s string) (float64, bool) {
	val,err := strconv.ParseFloat(s,64)
	return val, err == nil
}

//分割数学表达式
func SplitExpression(exp string) (s []string) {
	var tmp string
	for i:=0;i<len(exp);i++ {
		if exp[i] == '+' || exp[i] == '-' || exp[i] == '*' || exp[i] == '/' || exp[i] == '%' || exp[i] == '(' || exp[i] == ')' {
			if len(tmp) != 0 {
				s = append(s, tmp)
				s = append(s, string(exp[i]))
				tmp = ""
			}
		} else if exp[i] == '>' || exp[i] == '<'{
			s = append(s, tmp)
			tmp = ""
			tmp = tmp + string(exp[i])
			i++
			if !(exp[i]>='a'&&exp[i]<='z'||exp[i]>='A'&&exp[i]<='Z'||exp[i]>='0'||exp[i]<='9'||exp[i]=='_') {
				tmp = tmp + string(exp[i])
			} else {
				i--
			}
			s = append(s, tmp)
			tmp = ""
		} else {
			tmp = tmp+string(exp[i])
		}
	}
	if len(tmp) != 0 {
		s = append(s, tmp)
	}
	return
}

//判断if条件表达式
func JudgeCompareExpression(f *FuncBody, s string, ) bool {
	s = strings.TrimFunc(s, func(r rune) bool {
		return r == '(' || r == ')'
	})
	var (
		leftV,rightV string
		leftI,rightI float64
		tmp string
	)
	for i:=0;i<len(s);i++ {
		if s[i] == '>' {
			leftV = tmp
			if s[i+1] != '=' {
				rightV = s[i+1:]
			} else {
				rightV = s[i+2:]
			}
			if val, ok := f.VarMap[leftV]; ok {
				leftI = val.val
			} else {
				lv,_ := strconv.Atoi(leftV)
				leftI = float64(lv)
			}
			if val, ok := f.VarMap[rightV]; ok {
				rightI = val.val
			} else {
				rv,_ := strconv.Atoi(rightV)
				rightI = float64(rv)
			}
			return leftI > rightI
		} else if s[i] == '<' {
			leftV = tmp
			if s[i+1] != '=' {
				rightV = s[i+1:]
			} else {
				rightV = s[i+2:]
			}
			if val, ok := f.VarMap[leftV]; ok {
				leftI = val.val
			} else {
				lv,_ := strconv.Atoi(leftV)
				leftI = float64(lv)
			}
			if val, ok := f.VarMap[rightV]; ok {
				rightI = val.val
			} else {
				rv,_ := strconv.Atoi(rightV)
				rightI = float64(rv)
			}
			return leftI < rightI
		} else if s[i] == '=' && s[i+1] == '=' {
			leftV = tmp
			rightV = s[i+2:]
			if val, ok := f.VarMap[leftV]; ok {
				leftI = val.val
			} else {
				lv,_ := strconv.Atoi(leftV)
				leftI = float64(lv)
			}
			if val, ok := f.VarMap[rightV]; ok {
				leftI = val.val
			} else {
				rv,_ := strconv.Atoi(rightV)
				rightI = float64(rv)
			}
			return leftI == rightI
		} else {
			tmp += string(s[i])
		}
	}
	if val, ok := f.VarMap[rightV]; ok {
		leftI = val.val
	} else {
		lv,_ := strconv.Atoi(rightV)
		leftI = float64(lv)
	}
	return leftI != 0
}

//返回for循环的[low,high]和变量名
func ForJudgeLowHigh(s []string) (paraName string, low int, high int) {
	first := strings.Split(s[0], "=")
	second := SplitExpression(s[1])
	switch second[1] {
	case "<":
		low,_ = strconv.Atoi(first[1])
		high,_ = strconv.Atoi(second[2])
		high--
	case "<=":
		low,_ = strconv.Atoi(first[1])
		high,_ = strconv.Atoi(second[2])
	case ">":
		low,_ = strconv.Atoi(second[2])
		high,_ = strconv.Atoi(first[2])
		low++
	case ">=":
		low,_ = strconv.Atoi(second[2])
		high,_ = strconv.Atoi(first[2])
	}
	paraName = first[0]
	return
}

//返回for循环的[low,high]和变量名
func WhileJudgeLowHigh(s string,f *FuncBody) (paraName string, low int, high int) {
	second := SplitExpression(s)
	switch second[1] {
	case "<":
		low = int(f.VarMap[second[0]].val)
		high,_ = strconv.Atoi(second[2])
		high--
	case "<=":
		low = int(f.VarMap[second[0]].val)
		high,_ = strconv.Atoi(second[2])
	case ">":
		low,_ = strconv.Atoi(second[2])
		high = int(f.VarMap[second[0]].val)
		low++
	case ">=":
		low,_ = strconv.Atoi(second[2])
		high = int(f.VarMap[second[0]].val)
	}
	paraName = second[0]
	return
}

//得到字符串值
//1.直接string<-int
//2.从变量池里取
func GetNumValue()  {

}