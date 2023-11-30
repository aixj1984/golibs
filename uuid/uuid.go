// Package gorm is a wrapper for gorm.
package uuid

import (
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/sqids/sqids-go"
)

type Options struct {
	Alphabet  string
	MinLength uint8
}

type ShortUUID struct {
	alphabet  string
	minLength uint8
}

func New() *ShortUUID {
	return &ShortUUID{}
}

// GenV4 产生一个36位的唯一字符串
func (s *ShortUUID) GenV4() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// GenV4 产生一个16位的唯一字符串
func (s *ShortUUID) GenShort() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	bytes := u.Bytes()
	if len(bytes) < 16 {
		return "", errors.New("byte array length error")
	}

	newArray := make([]uint64, 5)
	var sum uint64 = 0 // 用于计算累加和

	// 按每4位进行求和
	for i := 0; i < len(bytes); i += 4 {
		// 把4位的byte值转换为一个整数，然后加到新数组
		value := uint64(bytes[i]) + uint64(bytes[i+1]) + uint64(bytes[i+2]) + uint64(bytes[i+3])
		newArray[i/4] = value

		// 累加所有的byte值
		sum += value
	}
	newArray[4] = sum

	// fmt.Println("New array:", newArray)

	sid, err := sqids.New(sqids.Options{
		MinLength: 16,
		// Alphabet:  "FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE",
	})
	if err != nil {
		return "", err
	}

	nid, _ := sid.Encode(newArray)
	// fmt.Printf("nid : %s\n", nid)

	return nid[:16], nil
}
