package uuid

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestV4(t *testing.T) {
	uid := New()
	id, err := uid.GenV4()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("t: %v\n", id)
}

func TestShort16(t *testing.T) {
	uid := New()
	id, err := uid.GenShort16()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("t: %v\n", id)
}

func TestShort16WithAlphabet(t *testing.T) {
	uid := NewWithAlphabet("0123456789")
	id, err := uid.GenShort16()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("t: %v\n", id)
}

func TestShort16_10(t *testing.T) {
	for index := 1; index < 10; index++ {
		uid := New()
		id, err := uid.GenShort16()
		if err != nil {
			t.Error(err.Error())
		}
		fmt.Printf("t: %v\n", id)
	}
}

func TestShort16WithAlphabet_10(t *testing.T) {
	for index := 1; index < 10; index++ {
		uid := NewWithAlphabet("0123456789")
		id, err := uid.GenShort16()
		if err != nil {
			t.Error(err.Error())
		}
		fmt.Printf("t: %v\n", id)
	}
}

func TestShort16WithAlphabet_Array(t *testing.T) {
	uid := NewWithAlphabet("0123456789")
	start := time.Now()
	ids, err := uid.GenShort16Array(1000)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("len: %d\n", len(ids))

	elapsed := time.Since(start)
	fmt.Printf("该函数执行完成耗时: %s\n", elapsed)
}

func TestShort16WithAlphabet_10000(t *testing.T) {
	idMap := make(map[string]bool)

	start := time.Now()

	uid := NewWithAlphabet("0123456789")
	for index := 1; index < 100000; index++ {

		id, err := uid.GenShort16()
		if err != nil {
			t.Error(err.Error())
		}
		// fmt.Printf("t: %v\n", id)
		idMap[fmt.Sprintf("%s", id)] = true
	}
	elapsed := time.Since(start)
	fmt.Printf("该函数执行完成耗时: %s\n", elapsed)

	start = time.Now()

	for index := 1; index < 100000; index++ {
		id, err := uid.GenShort16()
		if err != nil {
			t.Error(err.Error())
		}
		// fmt.Printf("t: %v\n", id)
		if isInMap(idMap, id) {
			fmt.Println("find same id : " + id)
		}
	}
	elapsed = time.Since(start)
	fmt.Printf("该函数执行完成耗时: %s\n", elapsed)
}

func TestShort24(t *testing.T) {
	u := New()
	fmt.Println(u.GenShort24()) // KwSysDpxcBU9FNhGkn2dCf
}
