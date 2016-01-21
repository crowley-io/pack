package main

import (
	"testing"
)

func TestHello(t *testing.T) {
	e := "Hello world!"
	v := hello()

	if v != e {
		t.Errorf("'%s' is not equals to '%s'", v, e)
	}

}
