package test

import (
	"fmt"
	"github.com/fulldog/utools/diyTypes"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestDiy(t *testing.T) {
	type location struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	type tx struct {
		User        diyTypes.DiyListAny[string]    `json:"User"`
		Location    diyTypes.DiyListAny[location]  `json:"Location"`
		LocationPtr diyTypes.DiyListAny[*location] `json:"LocationPtr"`
	}

	var ttt = tx{}
	s := `{"User":["xxxx","yyyyy"],"Location":[{"x":1,"y":2},{"x":2,"y":2}],"LocationPtr":[{"x":1,"y":2},{"x":2,"y":2}]}`
	err := jsoniter.Unmarshal([]byte(s), &ttt)
	fmt.Println(err)
	fmt.Println(ttt)
	b, err := jsoniter.Marshal(ttt)
	fmt.Println(string(b), err)
}
