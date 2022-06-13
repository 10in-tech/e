package e

type str struct{}

// O 把阴阳索引转为字符串
func (s *str) O(arg ...int) string {
	return stringify(arg, Opposite)
}

// E 把五行索引转为字符串
func (s *str) E(arg ...int) string {
	return stringify(arg, Element)
}

// G 把天干索引转为字符串
func (s *str) G(arg ...int) string {
	return stringify(arg, Gan)
}

// Z 把地支索引转为字符串
func (s *str) Z(arg ...int) string {
	return stringify(arg, Zhi)
}

// Spirits 把五神索引转为字符串
func (s *str) Spirits(arg ...int) string {
	return stringify(arg, Spirits5)
}

func stringify(ps []int, maps []string) (str string) {
	for _, p := range ps {
		str = join(str, maps[p])
	}
	return str
}
