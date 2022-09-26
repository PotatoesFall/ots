package form

import (
	"net/url"
	"testing"
)

type request struct {
	Int      int    `form:"intv"`
	Str      string `form:"strv"`
	Bool     bool   `form:"boolv"`
	NotFound int    `form:"notfound"`
}

func TestRead(t *testing.T) {
	form := url.Values{
		`intv`:     []string{"34"},
		`strv`:     []string{"blabla"},
		`boolv`:    []string{"true"},
		`whatever`: []string{"should be ignored"},
	}
	var req request
	err := Read(form, &req)
	if err != nil {
		t.Error(err)
	}

	expected := request{Int: 34, Str: "blabla", Bool: true, NotFound: 0}
	if req != expected {
		t.Errorf(`%#v should be equal to %#v`, req, expected)
	}
}
