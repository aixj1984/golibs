package uuid

import (
	"fmt"
	"testing"
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

func TestShort24(t *testing.T) {
	u := New()
	fmt.Println(u.GenShort24()) // KwSysDpxcBU9FNhGkn2dCf
}
