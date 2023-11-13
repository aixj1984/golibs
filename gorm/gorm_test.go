package gorm

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/aixj1984/golibs/zlog"
)

var (
	dao *Engine
	ctx = context.Background()
)

func TestMain(m *testing.M) {
	/*
			zlog.InitLogger(&zlog.Config{
				LogPath:    "./log/test.log",
				AppName:    "log-sample",
				Level:      -1,
				MaxSize:    1024,
				MaxAge:     3,
				MaxBackups: 4,
				Compress:   false,
			})

		db := New(&Config{
			Alias:        "test",
			Type:         "mysql",
			Server:       "127.0.0.1",
			Port:         3306,
			Database:     "test_db",
			User:         "root",
			Password:     "my-secret-pw",
			MaxIdleConns: 2,
			MaxOpenConns: 10,
			Charset:      "utf8mb4",
			MaxLeftTime:  time.Second * 10,
		})*/

	db := GetEngine()

	if err := db.gorm.Statement.Error; err != nil {
		zlog.Fatalf("%s", err.Error())
	}
	zlog.Infof("db1数据库连接成功")
	db.SetLogMode(true)
	dao = db

	m.Run()
}

// 用户信息表
type UserInfo struct {
	Id          int       `gorm:"column:id"`
	UnionId     string    `gorm:"column:union_id"`
	UserId      int       `gorm:"column:user_id"`
	NickName    string    `gorm:"column:nick_name"`
	HeadImg     string    `gorm:"column:head_img"`
	Description string    `gorm:"column:description"`
	Tag         string    `gorm:"column:tag"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

func TestInsert(t *testing.T) {

	dao.gorm.Migrator().DropTable(UserInfo{}.TableName())

	dao.gorm.AutoMigrate(UserInfo{})

	userObj := &UserInfo{
		Id:          1,
		UnionId:     strconv.FormatInt(time.Now().UnixMicro(), 10),
		UserId:      1,
		NickName:    "abc",
		HeadImg:     "",
		Description: "",
		Tag:         "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := dao.Context(ctx).Table(UserInfo{}.TableName()).WithContext(ctx).Create(userObj).Error

	if err != nil {
		zlog.Error("db", zlog.Fields{"error": err.Error()})
	}
}

func TestSelect(t *testing.T) {
	res := make([]*UserInfo, 0)
	dao.Context(ctx).Table(UserInfo{}.TableName()).WithContext(ctx).
		Where("id in (?)", []int64{1, 2, 3}).
		Select("id, nick_name, `created_at`,`updated_at`").
		Scan(&res)
	newID := uuid.New().String()

	ctx = context.WithValue(ctx, "trace_id", newID)
	logger := zlog.Logger().WithContext(ctx)

	if bytes, err := json.Marshal(res); err != nil {
		t.Fatal(err)
	} else {
		logger.Sugar().Infof("over %s", string(bytes))
	}

}

func TestCondition(t *testing.T) {
	cond := &DBConditions{
		And: map[string]interface{}{
			"id IN (?)": []int{1, 96, 97},
		},
		Not: map[string]interface{}{
			"id": []int{96},
		},
		Limit:  1,
		Offset: 0,
		Order:  "id DESC",
	}
	var records []*UserInfo

	err := cond.Fill(dao.gorm).Model(&UserInfo{}).Find(&records).Error
	if err != nil {
		t.Fatal(err)
	}
	zlog.Info("query result", zlog.Fields{"records": records})

}
