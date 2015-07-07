package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var (
	splitByWeek = regexp.MustCompile(`(周[一二三四五六日])([\w\W]{0,}?)<\/div>`)
	spiltAnime  = regexp.MustCompile(`<strong>[\w\W]{0,}?<em[^>]*>([0-9:]{1,}?)<\/em>[\w\W]{0,}?<u>(.{1,}?)<\/u><\/strong>`)
)

func getQuarter() string {
	now := time.Now()
	result := strconv.Itoa(now.Year())
	month := now.Month()
	switch {
	case 1 <= month && month < 4:
		result += "01"
	case 4 <= month && month < 7:
		result += "04"
	case 7 <= month && month < 10:
		result += "07"
	case 10 <= month && month <= 12:
		result += "10"
	}
	return result
}

type dayList [][2]string

func (l dayList) Len() int           { return len(l) }
func (l dayList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l dayList) Less(i, j int) bool { return l[i][0] < l[j][0] }

func getData() (*[7][][2]string, error) {
	// 1, 4, 7, 10
	resp, err := http.Get("http://acgdb.com/" + getQuarter() + "/bangumi")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	s := string(b)
	var result [7][][2]string
	week := splitByWeek.FindAllStringSubmatch(s, -1)
	for _, v := range week {
		anime := spiltAnime.FindAllStringSubmatch(v[2], -1)
		var i int
		switch v[1] {
		case "周日":
			i = 0
		case "周一":
			i = 1
		case "周二":
			i = 2
		case "周三":
			i = 3
		case "周四":
			i = 4
		case "周五":
			i = 5
		case "周六":
			i = 6
		}

		var r [][2]string
		for _, vv := range anime {
			time, name := vv[1], vv[2]
			r = append(r, [2]string{time, name})
		}
		sort.Sort(dayList(r))
		result[i] = r
	}
	return &result, nil
}

// 判断是否为半角字符
func isDbcCase(c rune) bool {
	if c >= 32 && c <= 127 {
		return true
	} else if c >= 65377 && c <= 65439 {
		return true
	}
	return false
}

// 获取字符串显示的长度
func getStrLen(s string) (length int) {
	for _, c := range s {
		if isDbcCase(c) {
			length += 1
		} else {
			length += 2
		}
	}
	return length
}

func makeStr(s string, l int) (str string) {
	for i := 0; i < l; i++ {
		str += s
	}
	return str
}

var week = [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}

func gen(i int, v [][2]string) string {
	var longest int
	for _, vv := range v {
		l := getStrLen(vv[1])
		if l > longest {
			longest = l
		}
	}

	half := (longest + 6 - 4) / 2
	left := makeStr(" ", half)
	right := makeStr(" ", half)
	if (longest+6-4)%2 != 0 {
		right += " "
	}

	s := "┌──────" + makeStr("─", longest) + "┐\n" +
		"│" + left + week[i] + right + "│\n" +
		"├─────┬" + makeStr("─", longest) + "┤\n"
	for _, vv := range v {
		s += fmt.Sprintf("│%v│%v%v│\n", vv[0], vv[1], makeStr(" ", longest-getStrLen(vv[1])))
	}
	s += "└─────┴" + makeStr("─", longest) + "┘"
	return s
}

func main() {
	data, err := getData()
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 2 && os.Args[1] == "all" {
		for i, v := range data {
			fmt.Println(gen(i, v))
		}
	} else {
		w := int(time.Now().Weekday())
		fmt.Println(gen(w, data[w]))
	}
	// fmt.Println("数据提供: http://acgdb.com")
}
