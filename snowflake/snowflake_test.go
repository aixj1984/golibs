package snowflake

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"testing"

	"github.com/aixj1984/golibs/zlog"
)

func TestMain(m *testing.M) {

	zlog.InitLogger(&zlog.Config{
		LogPath:    "./log/test.log",
		AppName:    "log-sample",
		Level:      -1,
		MaxSize:    1024,
		MaxAge:     3,
		MaxBackups: 4,
		Compress:   false,
	})

	m.Run()
}

func TestSnowFlake(t *testing.T) {

	id, err := NewID()

	if err != nil {
		t.Error("生成ID失败：", err)
	}

	if id == 0 {
		t.Error("生成ID失败：", err)
	}

	zlog.Infof(strconv.FormatUint(id, 10))
	zlog.Infof("%d", len(strconv.FormatUint(id, 10)))
	v := fmt.Sprintf("%x", md5.Sum([]byte(strconv.FormatUint(id, 10))))
	fmt.Println(v)
}
