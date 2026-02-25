// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package helper_test

import (
	"os"
	"testing"

	"github.com/cinar/indicator/v2/helper"
)

func TestRemove(t *testing.T) {
	f, err := os.Create("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	helper.Remove(t, "test.txt")

	if _, err := os.Stat("test.txt"); !os.IsNotExist(err) {
		t.Fatal("file not removed")
	}
}

func TestRemoveAll(t *testing.T) {
	err := os.Mkdir("testdir", 0700)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Create("testdir/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	helper.RemoveAll(t, "testdir")

	if _, err := os.Stat("testdir"); !os.IsNotExist(err) {
		t.Fatal("directory not removed")
	}
}
