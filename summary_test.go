package e

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSummaryColunms(t *testing.T) {

	typ, ele, season, ch := Summary.Columns([]int{5, 2, 7, 5}, []int{9, 0, 5, 1})

	fmt.Println(typ)
	fmt.Println(ele)
	fmt.Println(season)
	fmt.Println(ch)

	if !reflect.DeepEqual(typ, []int{0, 4, 1, 0, 1, 2, 4, 0}) {
		t.Error("五神计算错误")
	}

	if !reflect.DeepEqual(ele, []int{3, 2, 4, 3, 4, 0, 2, 3}) {
		t.Error("五行计算错误")
	}

	if season != -2 {
		t.Error("格局错误")
	}

	if !reflect.DeepEqual(ch, []string{"g5h27", "z6h01", "3hh591", "gzh29", "gzh75", "zah59"}) {
		t.Error("冲合关系计算错误")
	}

}

func TestChar(t *testing.T) {
	char := "9100110A"
	cs := Summary.Decode(char)

	if !Summary.Check(cs) {
		t.Fail()
	}
}
