package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/axgle/mahonia"
)

type CustomJson struct {
	Returncode int    `json:"returncode"`
	Message    string `json:"message"`
	Result     Result `json:"result"`
}
type Result struct {
	Specid         int             `json:"specid"`
	Paramtypeitems []Paramtypeitem `json:"paramtypeitems"`
}
type Paramtypeitem struct {
	Name       string      `json:"name"`
	Paramitems []Paramitem `json:"paramitems"`
}
type Paramitem struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// 反序列化为结构体对象
func ParseJson(url string) (CustomJson, error) {

	//fmt.Printf("原始字符串: %s\n", a)
	var c CustomJson
	resp, err1 := http.Get(url)
	if err1 != nil {
		return c, err1
	}
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	n, err2 := resp.Body.Read(buf)
	if err2 != nil && err2 != io.EOF {
		return c, err2
	}
	res1 := ConvertToString(string(buf[:n]), "gbk", "utf-8")
	if err := json.Unmarshal([]byte(res1), &c); err != nil {
		return c, err
	}
	return c, nil
}

func HttpGet(url string) (bytes.Buffer, error) {
	var res bytes.Buffer
	resp, err1 := http.Get(url)
	if err1 != nil {
		return res, err1
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			return res, err2
		}
		res.Write(buf[:])
	}
	return res, nil
}

func ConvertToString(src string, srcCode string, tagCode string) string {

	srcCoder := mahonia.NewDecoder(srcCode)

	srcResult := srcCoder.ConvertString(src)

	tagCoder := mahonia.NewDecoder(tagCode)

	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

	result := string(cdata)

	return result

}

var Needid = map[int]int{93: 1, 86: 1, 94: 1, 91: 1, 90: 1, 87: 1, 88: 1, 116: 1, 96: 1, 103: 1, 28: 1, 29: 1, 30: 1, 3: 1, 42: 1, 26: 1, 14: 1, 15: 1, 16: 1}

func (c *CustomJson) DealName() []string {
	wd := c.Result.Paramtypeitems
	name := []string{}
	for _, v := range wd {
		p := v.Paramitems
		for _, p_v := range p {
			_, ok := Needid[p_v.Id]
			if ok {
				Needid[p_v.Id] = len(name)
				name = append(name, p_v.Name)
			}
		}
	}
	return name
}
func (c *CustomJson) DealCont() []string {
	wd := c.Result.Paramtypeitems
	value := []string{}
	for _, v := range wd {
		p := v.Paramitems
		for _, p_v := range p {
			v, ok := Needid[p_v.Id]
			if ok {
				for len(value) < v {
					value = append(value, "")
				}
				value = append(value, p_v.Value)
			}
		}
	}
	return value
}

// id 里程 价格 日期
func SearchUrl(text []byte) [][]string {
	res := [][]string{}
	res_id := []string{}
	res_km := []string{}
	res_price := []string{}
	res_date := []string{}
	//过滤ul
	reg := regexp.MustCompile(`(<ul class="viewlist_ul" nocitycarposition="0")(([\s\S])*?)(<\/ul>)`)
	sub := reg.FindSubmatch(text)
	//过滤<li>
	reg2 := regexp.MustCompile(`(<li class="cards-li list-photo-li ")(([\s\S]*?))(<\/li>)`)
	if len(sub) == 0 {
		return res
	}
	sub2 := reg2.FindAllSubmatch(sub[2], -1)
	//过滤id
	regid := regexp.MustCompile(`(<li class="cards-li list-photo-li ")(([\s\S]*?))(specid=")(([\s\S]*?))(")`)
	//过滤div
	regdiv := regexp.MustCompile(`(<div class="cards-bottom">)(([\s\S])*?)(</div>)`)
	//过滤p
	regp := regexp.MustCompile(`(<p class="cards-unit">)(([\s\S])*?)(<)`)
	//过滤公里数
	regKM := regexp.MustCompile(`(<p class="cards-unit">)(([\s\S])*?)(公里)`)
	//过滤价格
	regPrice := regexp.MustCompile(`(<span class="pirce"><em>)(([\s\S])*?)(</em>)`)
	//过滤上牌时间
	regTime := regexp.MustCompile(`\d{4}-\d{2}`)
	for _, li := range sub2 {
		//if len(str[2])
		if len(li[0]) != 0 {
			div := regdiv.FindSubmatch(li[0])
			if len(div) == 0 {
				fmt.Println("div=0")
				continue
			}
			Id := regid.FindSubmatch(li[0])
			p := regp.FindSubmatch(div[2])
			if len(p) == 0 {
				fmt.Println("p=0")
				continue
			}

			Km := regKM.FindSubmatch(p[0])
			Price := regPrice.FindSubmatch(div[0])
			//fmt.Println(string(p[0]))
			Ntime := regTime.FindSubmatch(p[0])
			res_id = append(res_id, string(Id[5]))
			if len(Km) == 0 {
				res_km = append(res_km, "0")
				//fmt.Println(string(p[0]))
			} else {
				res_km = append(res_km, string(Km[2]))
			}
			if len(Price) == 0 {
				res_price = append(res_price, "")
				//fmt.Println(string(p[0]))
			} else {
				res_price = append(res_price, string(Price[2]))
			}

			if len(Ntime) != 0 {
				res_date = append(res_date, string(Ntime[0]))
			} else {
				res_date = append(res_date, "")
			}

		}
	}
	res = append(res, res_id, res_km, res_price, res_date)
	return res
}
func SearchId(text bytes.Buffer) string {

	//过滤ul
	reg := regexp.MustCompile(`(<input id="car_specid")(([\s\S])*?)(\/)`)
	sub := reg.FindSubmatch(text.Bytes())
	if len(sub) == 0 {
		fmt.Println("sub=0")
		return "1"
	}
	reg2 := regexp.MustCompile(`(value)(([\s\S])*?)(\/)`)
	sub2 := reg2.FindSubmatch(sub[0])
	if len(sub2) == 0 {
		fmt.Printf("%q\n", sub[0])
		return ""
	}
	res := "https://cacheapigo.che168.com/CarProduct/GetParam.ashx?specid=" + string(sub2[2][2:len(sub2[2])-2])
	return res
}
