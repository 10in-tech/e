package e

import (
	"fmt"
	"math"
	"sort"
)

type algorithm struct{}

// 通过数学规律快速计算关联关系

// GZ2O 根据天干索引转化为对应的阴阳索引
// 如天干戊的索引为4，代入函数得到1->阳，则戊为阳干
func (e *algorithm) GZ2O(gi int) int {
	return (gi + 1) % 2
}

// G2E 根据天干索引转化为对应的五行索引
// 如天干壬为8，代入函数得0->水，则壬对应五行的水
func (e *algorithm) G2E(gi int) int {
	return G2E[gi]
}

// Z2E 根据地支索引转化为对应的五行索引
// 如地支未为7，代入函数得3->土，则未对应五行的土
func (e *algorithm) Z2E(zi int) int {
	return Z2E[zi]
}

// Spirit 根据两五行元素获取对应的五神关系索引
// 如"我"日干为壬->水->0，"他"支为巳->火->2, 代入公式得3->才局
func (e *algorithm) Spirit(ie, oe int) int {
	return ((6-ie)%5 + oe) % 5
}

// Strong 格局是否是属于命强的局
// 五神中印(0)比(1)为命强，伤(2)才(3)杀(4)为命弱
func (e *algorithm) Strong(spirit int) bool {
	return spirit < 2
}

// Born 五行相生
// 水0->木1->火2->土3->金4->水5->···, 代入任意两个五行元素，可以判断前者是否生后者
func (e *algorithm) Born(parent, child int) bool {
	return (parent+1)%5 == child
}

// Restrain 五行相克
// 水0->火2->金4->木1->土3->水0->···，代入任意两个原属，可以判断这前者是否克后者
func (e *algorithm) Restrain(active, passive int) bool {
	return (active+2)%5 == passive
}

// Gh 天干相合
// 天干合是阴阳相吸，阴干配阳干，同时索引间隔4个。代入甲0和己5，可以得出两者为合的关系
func (e *algorithm) Gh(a, b int) bool {
	return b-a == 5
}

// Gc 天干相冲
// 天干冲为两个阴干或者阳干相斥，同时索引间隔5个。代入乙1和辛7，可以得出两者为冲的关系
func (e *algorithm) Gc(a, b int) bool {
	return b-a == 6
}

// Zh 地支相合
// 地支合是阴阳相吸，阴支配阳支，按环形排列，则以0，1和6，7分别为分界中线，两两相合，如代入卯3->戌10可以得出两者为合的关系
func (e *algorithm) Zh(a, b int) bool {
	return ((a+b)%12 == 1) && (b > a)
}

// Zc 地支相冲
// 地支冲为两个阴支或者阳支相斥，按环形排列，则以中心为对称点，两两相冲。如代入子0，午6，得出两者为冲的关系
func (e *algorithm) Zc(a, b int) bool {
	return (b - a) == 6
}

// Is3He 三合(增强中间)
// 三合仅地支有，按环形排列地支，则每隔三个选出一个，能选出三个地支组成三合
// ips bool Is Position Sensitive 位置敏感
func (e *algorithm) Is3He(a, b, c int, ips bool) bool {
	if ips {
		return inArray(fmt.Sprintf("%d.%d.%d", a, b, c), []string{"2.6.10", "5.9.1", "8.0.4", "11.3.7"})
	} else {
		r := []int{a, b, c}
		sort.Ints(r)
		return (r[0]+4) == r[1] && (r[0]+8) == r[2]
	}
}

// Is3Hui 三会(增强前两个)
// 三会仅地支有，按环形排列地支，从2开始，每连续三个为一个三会
// ips bool Is Position Sensitive 位置敏感
func (e *algorithm) Is3Hui(a, b, c int, ips bool) bool {
	if ips {
		return inArray(fmt.Sprintf("%d.%d.%d", a, b, c), []string{"2.3.4", "5.6.7", "8.9.10", "11.0.1"})
	} else {
		r := []int{a, b, c}
		sort.Ints(r)
		if r[0] == 2 || r[0] == 5 || r[0] == 8 {
			return (r[0]+1) == r[1] && (r[0]+2) == r[2]
		} else {
			return r[0] == 0 && r[1] == 1 && r[2] == 11
		}
	}
}

// GZh 干支暗合(天地鸳鸯媾合)
// 干支合，本身是天干合的变体。地支的主气天干和天干具有合的关系，如地支子的主气为癸，而天干戊癸合，则地支子和天干戊媾合简称戊子合
func (e *algorithm) GZh(g, z int) bool {
	g1 := ZwG[z][0]
	return math.Abs(float64(g-g1)) == 5
}

// GZc 干支暗冲(天地鸳鸯媾冲)
// 干支冲，本身是天干冲的变体。地支的主气天干和天干具有冲的关系，如地支子的主气为癸，而天干丁癸冲，则地支子和天干丁媾冲简称丁子冲
func (e *algorithm) GZc(g, z int) bool {
	g1 := ZwG[z][0]
	return math.Abs(float64(g-g1)) == 6
}

// Zah 地支暗媾合
// 天干合的变体，两个地支的主气天干相合。如巳的主气丙和酉的主气辛相合，则巳酉媾(暗)合
func (e *algorithm) Zah(z1, z2 int) bool {
	g1 := ZwG[z1][0]
	g2 := ZwG[z2][0]
	return math.Abs(float64(g1-g2)) == 5
}

// NextG 下一个天干索引
// 环形排列天干，获取当前天干的下一个天干索引
func (e *algorithm) NextG(i int) int {
	return e.next(i, 10)
}

// PrevG 上一个天干索引
// 环形排列天干，获取当前天干的上一个天干索引
func (e *algorithm) PrevG(i int) int {
	return e.prev(i, 10)
}

// NextZ 下一个地支
// 环形排列地支，获取当前地支的下一个地支索引
func (e *algorithm) NextZ(i int) int {
	return e.next(i, 12)
}

// PrevZ 上一个地支
// 环形排列地支，获取当前地支的上一个地支索引
func (e *algorithm) PrevZ(i int) int {
	return e.prev(i, 12)
}

// Year2month 年上起月法
func (e *algorithm) Year2month(yearGan int) int {
	return ((yearGan%5 + 1) * 2) % 10
}

// Day2hour 日上起时法
func (e *algorithm) Day2hour(dayGan int) int {
	return (dayGan % 5) * 2
}

// GZ2ci 天干地址索引转为60甲子索引
func (e *algorithm) GZ2ci(g, z int) int {
	return 5*((g+12-z)%12) + g
}

// 环形偏移下一个
func (e *algorithm) next(current, size int) int {
	return (current + 1) % size
}

// 环形偏移上一个
func (e *algorithm) prev(current, size int) int {
	return (current + size - 1) % size
}
