package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"golang.org/x/text/encoding/simplifiedchinese"
	"regexp"

	
)

func main() {
	
	// regdiv := regexp.MustCompile(`<div class="txtnav"(([\s\S])*?)(本章完)`)
	// div := regdiv.FindSubmatch([]byte(test))
	// rega := regexp.MustCompile(`(&emsp)(([\s\S])*?)(本章完)`)
	// a := rega.FindSubmatch(div[0])
	// regplace := regexp.MustCompile(`<br />`)
	// res := regplace.ReplaceAllString(string(a[2]), "");
	// fmt.Println(res)
  url := "https://www.ptwxz.com/html/15/15002/"
	p := Gethtml(url)
	res := Get(p.String())
	name := "火力为王"
	filename := "./"+name+".txt"
	os.Create(filename)
	for i:=0;i<len(res);i++ {
		fmt.Println(i)
		Write(res[i][1]+"\n", filename)
		out := Deala(url+res[i][0][1:len(res[i][0])-1])
		Write(out, filename)
	}
	
	// res := Get(p.String())
	// 
	// 
	
}

func Write(data, filePath string) {
    file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("文件打开失败", err)
    }
    //及时关闭file句柄
    defer file.Close()
    
    //写入文件时，使用带缓存的 *Writer
    write := bufio.NewWriter(file)
    
    write.WriteString(data)
    
    //Flush将缓存的文件真正写入到文件中
    write.Flush()
	//fmt.Printf("第%d章\n", i)
}

func Get(buf string)  [][]string {
	res := [][]string {}
	regdiv := regexp.MustCompile(`<div class="centent"(([\s\S])*?)<\/div>`)
	div := regdiv.FindSubmatch([]byte(buf))
	
	rega := regexp.MustCompile(`(<a href=)(([\s\S])*?)(>)(([\s\S])*?)(<)`)
	a := rega.FindAllSubmatch(div[0], -1)
	
	for i:=0;i<len(a);i++ {
        tail := []string{}
        tail = append(tail, string(a[i][2]), string(a[i][5]))
		res = append(res, tail)
	}
	//fmt.Println(res[0])
	return res
}

func Deala(url string) string {
	buf := Gethtml(url)
	//fmt.Println("buf", len(buf.Bytes()))
	regdiv := regexp.MustCompile(`(<br>)(([\s\S])*?)(</div>)`)
	div := regdiv.FindSubmatch(buf.Bytes())
	//fmt.Println("div=", len(div))
	
	a := div[0]
	
	regplace := regexp.MustCompile(`<br />`)
	regplace2 := regexp.MustCompile(`&emsp;`)
	regplace3 := regexp.MustCompile(`\s`)
	regplace4 := regexp.MustCompile(`。`)
	regplace5 := regexp.MustCompile(`&nbsp;&nbsp;&nbsp;&nbsp;`)
	regplace6 := regexp.MustCompile(`</div>`)
	regplace7 := regexp.MustCompile(`<br>`)
	res := regplace.ReplaceAllString(string(a), "")
	res2 := regplace2.ReplaceAllString(res, "")
	res3 := regplace3.ReplaceAllString(res2, "")
	res4 := regplace4.ReplaceAllString(res3,  "。\r\n")
	res5 := regplace5.ReplaceAllString(res4,  "")
	res6 := regplace6.ReplaceAllString(res5,  "")
	res7 := regplace7.ReplaceAllString(res6,  "")
	return res7
}

func Gethtml(url string) bytes.Buffer {
	var buf bytes.Buffer
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http err=", err)
		return buf
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read err=", err)
		return buf
	}
	utf8data, err := simplifiedchinese.GBK.NewDecoder().Bytes(body)
	buf.Write(utf8data)
	if err != nil {
		fmt.Println("utf8 transfer err=", err)
		return buf
	}
	return buf
}