package uuid

import (
	"fmt"
	"testing"
	"time"

	"github.com/bwmarrin/snowflake"
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

func TestShort16WithAlphabet_1000(t *testing.T) {
	idMap := make(map[string]bool)

	start := time.Now()

	uid := NewWithAlphabet("0123456789")
	for index := 1; index < 1000; index++ {

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

	for index := 1; index < 1000; index++ {
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

func TestSnow(t *testing.T) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())

	// Print out the ID's timestamp
	fmt.Printf("ID Time  : %d\n", id.Time())

	// Print out the ID's node number
	fmt.Printf("ID Node  : %d\n", id.Node())

	// Print out the ID's sequence number
	fmt.Printf("ID Step  : %d\n", id.Step())

	// Generate and print, all in one.
	fmt.Printf("ID       : %d\n", node.Generate().Int64())
}

func TestShortSnowflake(t *testing.T) {
	uid := NewWithAlphabet("0123456789")
	id, err := uid.GenSnowflake16()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("t: %v\n", id)

	idMap := make(map[string]bool)

	start := time.Now()

	for index := 1; index < 10000; index++ {
		id, err := uid.GenSnowflake16()
		if err != nil {
			t.Error(err.Error())
		}
		// fmt.Printf("t: %v\n", id)
		if _, ok := idMap[id]; ok {
			fmt.Printf("same id : %s\n", id)
			continue
		}
		idMap[fmt.Sprintf("%s", id)] = true
	}
	elapsed := time.Since(start)
	fmt.Printf("该函数执行完成耗时: %s\n", elapsed)
}
