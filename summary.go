package e

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type summary struct{}

// Columns 根据四柱八字，计算出相关信息
// typ 五神
// ele 五行
// season 季节
// ch 冲/合
func (s *summary) Columns(gs []int, zs []int) (typ []int, ele []int, season int, ch []string) {
	dge := Algorithm.G2E(gs[2]) // 计算出日干的五行

	water := false
	fire := false

	for _, g := range gs {
		ge := Algorithm.G2E(g)
		ele = append(ele, ge)                        // 五行
		typ = append(typ, Algorithm.Spirit(dge, ge)) // 五神

		if ge == 0 {
			water = true
		}
		if ge == 2 {
			fire = true
		}
	}
	for _, z := range zs {
		ze := Algorithm.Z2E(z)
		ele = append(ele, ze)                        // 五行
		typ = append(typ, Algorithm.Spirit(dge, ze)) // 五神

		if ze == 0 {
			water = true
		}
		if ze == 2 {
			fire = true
		}
	}

	season = s.season(zs[1], water, fire) // 取出月支计算寒燥

	gs = s.sortAndUnique(gs)
	zs = s.sortAndUnique(zs)

	// 冲合
	ch = append(ch, s.couple(gs, func(a, b int) string { // 天干冲合
		if Algorithm.Gh(a, b) {
			return join("g5h", s.String(a), s.String(b))
		}
		if Algorithm.Gc(a, b) {
			return join("g4c", s.String(a), s.String(b))
		}
		return ""
	})...)
	ch = append(ch, s.couple(zs, func(a, b int) string { // 地支冲合
		if Algorithm.Zh(a, b) {
			return join("z6h", s.String(a), s.String(b))
		}
		if Algorithm.Zc(a, b) {
			return join("z6c", s.String(a), s.String(b))
		}
		return ""
	})...)
	ch = append(ch, s.hh3(zs)...) // 三合三会

	for _, g := range gs { // 干支暗冲合
		for _, z := range zs {
			if Algorithm.GZh(g, z) {
				ch = append(ch, join("gzh", s.String(g), s.String(z)))
				break
			}
			if Algorithm.GZc(g, z) {
				ch = append(ch, join("gzc", s.String(g), s.String(z)))
			}
		}
	}

	ch = append(ch, s.couple(zs, func(a, b int) string { // 地支暗合
		if Algorithm.Zah(a, b) {
			return join("zah", s.String(a), s.String(b))
		}
		return ""
	})...)

	return
}

// 寒燥
func (s *summary) season(monthZ int, water, fire bool) (season int) {
	switch monthZ {
	case 11, 0, 1: // 亥子丑
		season = -2
	case 5, 6, 7: // 巳午未
		season = 2
	case 2, 10: // 寅戌
		if water && !fire { // 格局中有水无火
			season = -1
		} else {
			season = 0
		}
	case 4, 8: // 卯酉
		if !water && fire { // 格局中无水有火
			season = 1
		} else {
			season = 0
		}
	default:
		season = 0
	}
	return
}

// 成对比较冲合
func (s *summary) couple(list []int, f func(a, b int) string) []string {
	var ch []string
	l := len(list)
	for i := 0; i < l-1; i++ {
		for j := i + 1; j < l; j++ {
			flag := f(list[i], list[j])
			if flag != "" {
				ch = append(ch, flag)
			}
		}
	}
	return ch
}

// 把干支索引转化为字符串
func (s *summary) String(i int) string {
	if -1 < i && i < 10 {
		return strconv.Itoa(i)
	}
	switch i {
	case 10:
		return "A"
	case 11:
		return "B"
	default:
		panic(errors.New("Index of Gan or Zhi is out of range "))
	}
}

func (s *summary) Encode(unique [8]int) string {
	var buffer bytes.Buffer
	for _, c := range unique {
		buffer.WriteString(s.String(c))
	}
	return buffer.String()
}

// Decode 把编码的四柱八字，转换为未编码状态
func (s *summary) Decode(columns string) [8]int {
	chars := strings.Split(columns, "")
	unique := [8]int{}
	for c := 0; c < 8; c++ {
		switch chars[c] {
		case "0":
			unique[c] = 0
		case "1":
			unique[c] = 1
		case "2":
			unique[c] = 2
		case "3":
			unique[c] = 3
		case "4":
			unique[c] = 4
		case "5":
			unique[c] = 5
		case "6":
			unique[c] = 6
		case "7":
			unique[c] = 7
		case "8":
			unique[c] = 8
		case "9":
			unique[c] = 9
		case "A":
			unique[c] = 10
		case "B":
			unique[c] = 11
		}
	}
	return unique
}

// Check 检查四柱八字是否正确
func (s *summary) Check(cs [8]int) bool {
	for i := 0; i < 4; i++ { // 检查阴阳相对
		if cs[i]%2 != cs[i+4]%2 {
			return false
		}
	}

	if (Algorithm.Year2month(cs[0])+(cs[5]+10)%12)%10 != cs[1] { // 年上起月
		return false
	}

	hg := Algorithm.Day2hour(cs[2])
	if cs[7] == 0 { // 子时有两种情况
		if (hg+12)%10 == cs[3] { // 当为晚子时时
			return true
		}
	}
	if (hg+cs[7])%10 != cs[3] {
		return false
	}
	return true
}

// 三合三会计算
func (s *summary) hh3(list []int) []string {
	var buffer bytes.Buffer
	for _, i := range list {
		buffer.WriteString(s.String(i))
	}
	zStr := buffer.String()
	if len(zStr) >= 3 {
		for k, n := range HH3 {
			if strings.Contains(zStr, k) {
				return []string{join("3hh", n)}
			}
		}
	}
	return []string{}
}

func (s *summary) sortAndUnique(list []int) (nl []int) {
	var sl [12]bool
	for _, i := range list {
		sl[i] = true
	}
	for i := 0; i < 12; i++ {
		if sl[i] == true {
			nl = append(nl, i)
		}
	}
	return
}
