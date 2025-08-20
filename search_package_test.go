package main

import (
	"testing"
)

func Test_SearchHandler(t *testing.T) {
	res, err, shouldTerminate := search_package_core(t.Context(), "asdf")
	if err != nil {
		t.Log("Error:", err)
		t.Error(err)
	}
	t.Log("Result:", res, "shouldTerminate:", shouldTerminate)
	t.Fail()
}
