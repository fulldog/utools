package test

import (
	"fmt"
	"github.com/fulldog/utools/diyTypes"
	"github.com/fulldog/utools/timex"
	jsoniter "github.com/json-iterator/go"
	"testing"
	"time"
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
		UnixTime    diyTypes.UnixTime              `json:"UnixTime"`
		DateTime    diyTypes.DiyTime               `json:"DateTime,omitempty"`
	}

	var ttt = tx{}
	s := `{"User":["xxxx","yyyyy"],"Location":[{"x":1,"y":2},{"x":2,"y":2}],"LocationPtr":[{"x":1,"y":2},{"x":2,"y":2}],"UnixTime":1682006515,"DateTime":null}`
	err := jsoniter.Unmarshal([]byte(s), &ttt)
	fmt.Println(err)
	fmt.Println(ttt)
	ttt.UnixTime.Layout = timex.DateOnly
	b, err := jsoniter.Marshal(ttt)
	fmt.Println(string(b), err)

	tm, _ := time.Parse(timex.DateOnly, "2020-02-02")
	fmt.Println(tm.String())
}
