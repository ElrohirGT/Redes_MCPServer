package main

import (
	"testing"
)

func Test_SearchHandler(t *testing.T) {
	search_package_core(t.Context(), "asdf")
	t.Fail()
}
