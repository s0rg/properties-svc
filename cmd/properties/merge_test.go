package main

import (
	"testing"
)

func TestMergeBundles(t *testing.T) {
	var (
		a, b []Bundle
		c    []int
	)

	c = MergeBundles(a, b)
	if len(c) != 0 {
		t.Fatal("step 1 fail")
	}

	a = []Bundle{
		{ID: 1, Name: "a1"},
	}

	c = MergeBundles(a, b)
	if len(c) != 1 || c[0] != 1 {
		t.Fatal("step 2 fail")
	}

	b = []Bundle{
		{ID: 2, Name: "b1", ParentID: 1},
	}

	c = MergeBundles(a, b)
	if len(c) != 1 || c[0] != 2 {
		t.Fatal("step 3 fail")
	}

	a = append(a, Bundle{ID: 3, Name: "a2", ParentID: 4})

	c = MergeBundles(a, b)
	if len(c) != 2 {
		t.Fatal("step 4 fail")
	}

	b = append(b, Bundle{ID: 4, Name: "b2"})

	c = MergeBundles(a, b)
	if len(c) != 2 {
		t.Fatal("step 5 fail")
	}
}

func TestDropBundles(t *testing.T) {
	var (
		a, b []Bundle
		c    []int
	)

	c = DropBundles(a, b)
	if len(c) != 0 {
		t.Fatal("step 1 fail")
	}

	a = []Bundle{
		{ID: 1, Name: "a1"},
	}

	c = DropBundles(a, b)
	if len(c) != 1 || c[0] != 1 {
		t.Fatal("step 2 fail")
	}

	c = DropBundles(b, a)
	if len(c) != 0 {
		t.Fatal("step 3 fail")
	}

	b = []Bundle{
		{ID: 2, Name: "b1", ParentID: 1},
	}

	c = DropBundles(a, b)
	if len(c) != 1 || c[0] != 1 {
		t.Fatal("step 4 fail")
	}

	c = DropBundles(b, b)
	if len(c) != 1 || c[0] != 1 {
		t.Fatal("step 5 fail")
	}
}
