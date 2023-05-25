package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"GOWORK/project/spider/tool"
)
type OutChan struct {
	Num int
	Data bytes.Buffer
}
type InChan struct {
	Num int
	Contents []string
}
func main() {
	head := 0
	n := -1
	var url, name string
	// fmt.Println("url = ")
	// fmt.Scanln(&url)
	// fmt.Println("name = ")
	// fmt.Scanln(&name)
	url = "https://www.69shu.com/A48044/"
	name = "我的属性修行人生"
	path := "D:/GOWORK/project/"+name+"/"
	err := os.Mkdir(path, 0666)
	if err != nil {
	   fmt.Println(err)
	}
	buf := tool.Gethtml(url)
	intchan := make(chan InChan, 200)
	
	num := 24
	as, n := tool.Geta(buf, head, n)
	outchan := make(chan OutChan, n)
	exitchan := make(chan bool, num)
	go Puta(as, intchan, head, n)
	for i:=0;i<num;i++ {
		go Deal(intchan, outchan, exitchan)
	}
	//go Write(intchan, path, exitchan)
	go func()  {
		for i:=0;i<num;i++ {
			<-exitchan
		}
		close(outchan)
	}()
	//取结果
	for {
		res, ok := <- outchan
		if !ok {
			break
		}
		//fmt.Println(res.Num)
		tool.WriteData(res.Data.String(), path+ strconv.Itoa(res.Num)+".txt")
	}
	
}
func Puta(as [][]string, intchan chan InChan, head, n int) {
	for i:=head;i<n;i++ {
		intchan<-InChan{
			Num: i,
			Contents: as[i-head],
		}
	}
	close(intchan)
}

func Deal(intchan chan InChan, outchan chan OutChan, exitchan chan bool) {
	for {
		v, ok := <-intchan
		if !ok {
			break
		}
		a := v.Contents
		var Buf bytes.Buffer
		Buf.WriteString("\r\n")
		Buf.WriteString(a[1])
		Buf.WriteString("\r\n")
		res := tool.Deala(a[0][1:len(a[0])-1])
		Buf.WriteString(res)
		fmt.Println(v.Num)
		outchan<-OutChan{
			Num: v.Num,
			Data: Buf,
		}
	}
	fmt.Println("一个协程关闭")
	exitchan<-true
}






