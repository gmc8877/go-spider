package tool

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/text/encoding/simplifiedchinese"
)

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

func Geta(buf bytes.Buffer, head, n int) ([][]string, int) {
	res := [][]string{}
	regdiv := regexp.MustCompile(`<div class="catalog" id="catalog"(([\s\S])*?)<\/div>`)
	div := regdiv.FindSubmatch(buf.Bytes())
	rega := regexp.MustCompile(`(<a href=)(([\s\S])*?)(>)(([\s\S])*?)(<)`)
	a := rega.FindAllSubmatch(div[0], -1)
	if n == -1 {
		n = len(a)
	}
	for i := head; i < n; i++ {
		tail := []string{}
		tail = append(tail, string(a[i][2]), string(a[i][5]))
		res = append(res, tail)
	}
	return res, n
	//fmt.Println(res[0])
}

func Deala(url string) string {
	buf := Gethtml(url)
	//fmt.Println("buf", len(buf.Bytes()))
	regdiv := regexp.MustCompile(`<div class="txtnav"(([\s\S])*?)(<div class="bottom-ad">)`)
	div := regdiv.FindSubmatch(buf.Bytes())
	if len(div) == 0 {
		fmt.Println("无法访问", url)
		file, err := os.OpenFile("D:/GOWORK/project/我家娘子不是妖/err.txt", os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
		//及时关闭file句柄
		defer file.Close()

		//写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(file)

		write.WriteString(url+"\n")

		//Flush将缓存的文件真正写入到文件中
		write.Flush()
		return "空"
	}
	//fmt.Println("div=", len(div))
	//fmt.Println(string(div[0]))
	rega := regexp.MustCompile(`(&emsp)(([\s\S])*?)(<div class="bottom-ad">)`)
	a := rega.FindSubmatch(div[0])
	if len(a) == 0 {
		fmt.Println("这章为空", url)
		return "空"
	}
	//fmt.Println("a", len(a))
	regplace := regexp.MustCompile(`<br />`)
	regplace2 := regexp.MustCompile(`&emsp;`)
	res := regplace.ReplaceAllString(string(a[2]), "")
	res2 := regplace2.ReplaceAllString(res, "	")
	return res2
	//如果空格有问题用下面这个
	// regplace := regexp.MustCompile(`<br />`)
	// regplace2 := regexp.MustCompile(`&emsp;`)
	// regplace3 := regexp.MustCompile(`\s`)
	// regplace4 := regexp.MustCompile(`。`)
	// res := regplace.ReplaceAllString(string(a[2]), "")
	// res2 := regplace2.ReplaceAllString(res, "")
	// res3 := regplace3.ReplaceAllString(res2, "")
	// res4 := regplace4.ReplaceAllString(res3,  "。\r\n")
	// return res4
}

func WriteData(data, filePath string) {
	os.Create(filePath)
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
