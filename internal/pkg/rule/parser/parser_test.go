package parser

import (
	"context"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	p := New()
	res, err := p.ParseDirectory(context.Background(), "./example/rules")
	if err != nil {
		t.Log(err.Error())
	}
	t.Log(res)
}
