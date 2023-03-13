package gorm

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"testing"
	"time"
)

type User struct {
	gorm.Model
	Name         string
	Email        *string
	Age          uint8 `gorm:"column:age"` //指定列名
	Birthday     time.Time
	MemberNumber sql.NullString
}

// 指定表名
func (u User) TableName() string {
	return "user_t"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	println("before create")
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	println("after create")
	return
}

func Test_CRUD(t *testing.T) {
	// 连接数据库
	dsn := "root:wang@tcp(127.0.0.1:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 迁移 schema
	db.AutoMigrate(&User{})

	// 增
	result := db.Create(&User{Name: "Jinzhu", Age: 18, Birthday: time.Now()})
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	//result.Error  //执行结果error
	//result.RowsAffected  // 影响的行数

	//tx := db.Session(&gorm.Session{DryRun: true})

	// 查
	var user User
	db.Where("name = ?", "Jinzhu").First(&user)

	// 打印 SQL，但不执行
	tx := db.Session(&gorm.Session{
		DryRun: true,
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 改
	tx.Model(&user).Update("name", "wangyanlei")
	// 一次改多个
	tx.Model(&user).Updates(User{Name: "Linda", Age: 20}) // non-zero fields
	tx.Model(&user).Updates(map[string]interface{}{"Name": "Linda", "Age": 22})

	// 删
	tx.Delete(&user)
}
