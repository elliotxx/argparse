package main

import (
	"./argparse"
	"bufio"
	"fmt"
	"os"
)

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func readLines(fp *os.File) ([]string, int) {
	// 从文件中读取所有行
	scanner := bufio.NewScanner(fp)
	lines := []string{}
	lineCnt := 0
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		lineCnt += 1
	}
	return lines, lineCnt
}

func main() {
	// cat 的参数设置
	isf := argparse.Bool("f", false, "倒序输出文本")
	isn := argparse.Bool("n", false, "输出带行号的文本")
	isb := argparse.Bool("b", false, "按照二进制模式显示文本")
	p 	:= argparse.Int("p", -1, "输出 n 行")
	ish := argparse.Bool("h", false, "帮助信息")
	ish2 := argparse.Bool("help", false, "帮助信息")
	err := argparse.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 首先判断是否要输出帮助信息
	// 如果是，那么输出完帮助信息就结束
	if *ish || *ish2 {
		argparse.Help()
		return
	}

	// cat 功能实现
	// 打开文件
	// 第一个未解析参数作为文件名，比如 cat -p 10 a.txt，那么 a.txt 作为文件名
	if argparse.OtherArgLen() <= 0 {
		fmt.Println("请提供文件名")
		return
	}
	filename, err := argparse.OtherArg(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fp, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()
	// 从文件中读取全部行
	lines, _ := readLines(fp)
	if *isf {
		// 倒序输出
		lines = reverse(lines)
	}
	for i, line := range lines {
		if *p != -1 && i >= *p {
			// 最多输出 p 行
			break
		}
		if *isb {
			// 二进制模式显示文本
			line = fmt.Sprintf("% x", line)
		}
		if *isn {
			// 输出带行号的文本
			line = fmt.Sprintf("%d\t%s", i+1, line)
		}
		fmt.Println(line)
	}
}
