# argparse
A 200-line command-line parsing library

blog: [200 行代码实现一个简单的命令行解析库](http://yangyingming.com/article/456/)

## Installation
```
$ go get -u -v github.com/elliotxx/argparse
```

## Usage
```
package main

import (
  "fmt"
  "github.com/elliotxx/argparse"
)

func main() {
    isr   := argparse.Bool("r", false, "Output text in reverse order")
    n     := argparse.Int("n", -1, "Output n lines")
    ish   := argparse.Bool("h", false, "Help information")
    isH   := argparse.Bool("help", false, "Help information")
    
    err := argparse.Parse()
    if err != nil {
        fmt.Println(err)
        return
    }

    if *ish || *isH {
        argparse.Help()
        return
    }
    fmt.Printf("isr=%v, n=%d\n", *isr, *n)
}
```
Ouput:
```
$ go run cat.go -h
Usage of cat
    -h,--help	bool	Help information
    -n	int	Output n lines
    -r	bool	Output text in reverse order
```
```
$ go run cat.go -nr 10
isr=true, n=10
```
