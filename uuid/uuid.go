// Package gorm is a wrapper for gorm.
package uuid

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/gofrs/uuid/v5"

	"github.com/sqids/sqids-go"

	"github.com/lithammer/shortuuid/v4"
)

// UnionUUID 统一的UUID生成入口
type UnionUUID struct {
	alphabet  string
	minLength uint8
	snowNode  *snowflake.Node
}

func New() *UnionUUID {
	node, _ := snowflake.NewNode(1)
	return &UnionUUID{
		alphabet:  "",
		minLength: 16,
		snowNode:  node,
	}
}

func NewWithAlphabet(alphabet string) *UnionUUID {
	node, _ := snowflake.NewNode(1)
	return &UnionUUID{
		alphabet:  alphabet,
		minLength: 16,
		snowNode:  node,
	}
}

// GenV4 产生一个36位的唯一字符串
func (s *UnionUUID) GenV4() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// GenShort16 产生一个16位的唯一字符串,通过测试，会重复,需要考虑去重，V4的重复概率低
func (s *UnionUUID) GenShort16() (string, error) {
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

	sid, err := sqids.New(sqids.Options{
		MinLength: s.minLength,
		Alphabet:  s.alphabet,
	})
	if err != nil {
		return "", err
	}

	nid, _ := sid.Encode(newArray)
	return nid[:16], nil
}

// GenShort16Array 生成指定长度不重复的数组，只能保障本次不重复
func (s *UnionUUID) GenShort16Array(length int) ([]string, error) {
	if length == 0 || length > 100000 {
		return nil, errors.New("len is to large")
	}

	idMap := make(map[string]bool)

	for len(idMap) < length {

		id, err := s.GenShort16()
		if err != nil {
			return nil, err
		}
		if _, ok := idMap[id]; !ok {
			idMap[id] = true
		}
	}

	keys := make([]string, 0, len(idMap))
	for k := range idMap {
		keys = append(keys, k)
	}
	return keys, nil
}

// GenShort24 以UUID为基础，再对字符进行压缩，该结果可以解密
func (s *UnionUUID) GenShort24() string {
	u := shortuuid.New()
	return fmt.Sprintf("%024s", u)
}

func (s *UnionUUID) GenShort24With() string {
	u := shortuuid.NewWithAlphabet("0123456789")
	return fmt.Sprintf("%024s", u)
}

func (s *UnionUUID) GenSnowflake16() (string, error) {
	var nodeIDMask int64 = 1023 << 12 // mask to get Node ID
	var sequenceIDMask int64 = 4095   // mask to get sequence ID

	id := s.snowNode.Generate()
	// base2 := id.Base2()

	val := id.Int64()
	// return id.String(), nil

	nodeIDPart := val &^ nodeIDMask        // remove Node ID only
	sequenceIDPart := val & sequenceIDMask // get sequence ID part

	// newNodeIDPart := nodeIDPart >> 10 // shift Node ID part right by 10 bits
	// 考虑到后面的序列号不会那么大，12位太多，减少到整体48位，需要去掉3位
	newNodeIDPart := nodeIDPart >> 13 // shift Node ID part right by 10 bits

	newVal := newNodeIDPart | sequenceIDPart // combine the two parts

	base2 := strconv.FormatInt(newVal, 2) // convert to binary

	lenBase2 := len(base2)
	// fmt.Printf("bit length : %d\n", lenBase2)
	// 每6个为一组，共8个数值
	size := 6
	if lenBase2%size != 0 {
		// 不足size位的补0
		for i := 0; i < size-lenBase2%size; i++ {
			base2 = base2 + "0"
		}
	}

	var result []uint64
	for i := 0; i < len(base2); i += size {
		num, err := strconv.ParseInt(base2[i:i+size], 2, 32)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		result = append(result, uint64(num))
	}

	sid, err := sqids.New(sqids.Options{
		MinLength: s.minLength,
		Alphabet:  s.alphabet,
	})
	if err != nil {
		return "", err
	}
	nid, err := sid.Encode(result)
	if err != nil {
		return "", err
	}

	if len(nid) < 16 {
		// fmt.Println("nid is " + fmt.Sprintf("%016s", nid))
		return fmt.Sprintf("%016s", nid), nil
	}
	// 前面几位基本一样，所以很容易重，这里直接取后16位
	return nid[len(nid)-16:], nil
}

func isInMap(m map[string]bool, target string) bool {
	_, ok := m[target]
	return ok
}
