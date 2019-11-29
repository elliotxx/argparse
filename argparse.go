// argparse 特性：
// 支持 int、bool、string 三种类型的参数，比如 -p 10 或者 -f xxx.txt
// 支持同时识别多个参数，比如 -nfp 10
// 支持同时定义长/短参数，比如 -h 或者 --help，你需要给它们定义相同的 usage，比如 "帮助信息"
// 		"-" 后面的是短参数，"--" 后面的是长参数，错误的使用会报错，比如 -help 是错误的
// 输出帮助信息（--help）时
// 		长/短参数会自动合并在一起，像这样：-h,--help，它们必须有相同的 usage
//		会根据参数名进行排序

package argparse

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var argMap = make(map[string]ArgItem) // 参数配置
var otherArg []string                 // 没解析到的其它参数（不包含在参数配置中的参数）

// 自定义类型，方便添加 Set 方法，以实现 Value 接口
type boolVal 	bool
type intVal 	int
type stringVal 	string

// 该接口类型的变量可调用 Set 方法进行赋值
type Value interface {
	Set(v string) error			// 设置值
}

// 存储参数配置和当前值
type ArgItem struct {
	Type 	string		// 类型
	Value 	Value		// 参数值
	Usage 	string		// 参数用法
}

// 可以解析三类参数：bool、int、string
func (x *boolVal) Set(v string) error {
	r, err := strconv.ParseBool(v)
	if err != nil {
		return err
	}
	*x = boolVal(r)
	return nil
}

func (x *intVal) Set(v string) error {
	r, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	*x = intVal(r)
	return nil
}

func (x *stringVal) Set(v string) error {
	*x = stringVal(v)
	return nil
}

// 添加参数初始配置到 argMap
func Bool(name string, defaultValue bool, usage string) (*bool) {
	x := new(bool)
	*x = defaultValue
	v := (*boolVal)(x)
	arg := ArgItem{
		Type: 	"bool",
		Value:  v,
		Usage:	usage,
	}
	argMap[name] = arg
	return x
}

func Int(name string, defaultValue int, usage string) (*int) {
	x := new(int)
	*x = defaultValue
	v := (*intVal)(x)
	arg := ArgItem{
		Type: 	"int",
		Value:  v,
		Usage:	usage,
	}
	argMap[name] = arg
	return x
}

func String(name string, defaultValue string, usage string) (*string) {
	x := new(string)
	*x = defaultValue
	v := (*stringVal)(x)
	arg := ArgItem{
		Type: 	"string",
		Value:  v,
		Usage:	usage,
	}
	argMap[name] = arg
	return x
}

func fullArg(arg string) string {
	// 返回参数带 "-" 的完整形式，短参数用"-"，长参数用"--"，比如 -h，或者 --help
	if len(arg) <= 1 {
		return "-" + arg
	} else {
		return "--" + arg
	}
}

func Help() {
	// 输出帮助信息
	// 排序并合并同 usage 的参数，比如 -h --help 合并为 -h,--help
	keyIndex := []string{}
	usageMap := map[string][]string{}
	// 合并同 usage 的参数
	for k, v := range argMap {
		// 合并所有同 usage 的参数，usageMap 记录 usage => 该 usage 对应的参数列表
		// 比如 "帮助信息" => [h,help]
		if keys, ok := usageMap[v.Usage]; ok {
			// 遇到同 usage 参数，累加到列表后面
			keys = append(keys, k)
			usageMap[v.Usage] = keys
		} else {
			usageMap[v.Usage] = []string{k}
		}
		// 提取所有参数名，以便之后排序
		keyIndex = append(keyIndex, k)
	}
	// 排序
	sort.Strings(keyIndex)
	// 按顺序输出，并合并同 usage 参数
	fmt.Printf("Usage of %s\n", os.Args[0])
	for _, oldKey := range keyIndex {
		usage := argMap[oldKey].Usage
		argType := argMap[oldKey].Type
		if newKeys, ok := usageMap[usage]; ok {
			sort.Strings(newKeys)
			if len(newKeys) <= 1 {
				fmt.Printf("    %s\t%s\t%s\n", fullArg(oldKey), argType, usage)
			} else {
				// 遇到了同 usage 参数，合并输出，并删除对应 key
				fullKeys := []string{}
				for _, k := range newKeys {
					fullKeys = append(fullKeys, fullArg(k))
				}
				fmt.Printf("    %s\t%s\t%s\n", strings.Join(fullKeys, ","), argType, usage)
				delete(usageMap, usage)
			}
		}
	}
}

// 没解析到的其它参数（不包含在参数配置中的参数）
func OtherArgLen() int {
	return len(otherArg)
}

func OtherArg(i int) (s string, err error) {
	defer func() {
		if e,ok := recover().(error); ok {
			s = ""
			err = e
		}
	}()
	return otherArg[i], nil
}

func Parse() error {
	// args parse
	length := len(os.Args)
	for i:=1; i<length; i++ {
		// 依次读取命令行参数
		v := []byte(os.Args[i])
		if len(v) > 1 && v[0] == '-' {
			// 解析每个参数
			argGroup := []string{}
			if len(v) > 2 && v[1] == '-' {
				// 解析 --name 这类参数
				argGroup = append(argGroup, string(v[2:]))
			} else {
				// 解析 -afn 这类参数组
				for _, arg := range v[1:] {
					argGroup = append(argGroup, string(arg))
				}
			}
			// 处理解析到的每个参数
			readNext := false		// 是否读取了下一个值
			for _, arg := range argGroup {
				if ai, ok := argMap[arg]; ok { // 判断是否是需要解析的参数
					if ai.Type == "bool" {		// 如果该参数是布尔型
						ai.Value.Set("true")	// 那么值应该是 true
					} else {					// 如果该参数不是布尔型
						if readNext {			// 已经读取过下一个值，报错
							return fmt.Errorf("参数解析错误，有多个需要后续值的参数")
						} else {				// 没读取过，将下一个参数作为该参数的值，比如 -p 10，把 10 作为 p 参数的值
							readNext = true
							i += 1
							if i >= length {
								return fmt.Errorf("参数 -%s 需要值", arg)
							}
							if err := ai.Value.Set(os.Args[i]); err != nil {
								return fmt.Errorf("%s 无法作为参数 -%s 的值", os.Args[i], arg)
							}
						}
					}
				} else {
					return fmt.Errorf("参数 -%s 无法识别", arg)
				}
			}
		} else {
			otherArg = append(otherArg, string(v))
		}
	}
	return nil
}