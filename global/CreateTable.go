package global

import (
	"fmt"
	"log"
	"miaoshaSystem/sql"
)

/*type User struct {
	Name     string `gorm:"type:varchar(20);not_null;unique" json:"name"`
	Age      int    `gorm:"type:int;not_null" json:"age"`
	Address  string `gorm:"type:varchar(100);not_null" json:"address"`
	Avatar   string `gorm:"type:varchar(100);not_null" json:"avatar"`
	ID       int    `gorm:"type:int;not_null;auto_increment" json:"id"`
	Password string `gorm:"type:varchar(100);not_null" json:"password"`
} // 购物车打算用redis存储，记一笔免得后面忘记了
type Product struct {
	Name            string `gorm:"type:varchar(20);not_null" json:"name"`
	ID              string `gorm:"type:int;not_null;auto_increment" json:"id"`
	Num             int    `gorm:"type:int" json:"num"`
	Producter       string `gorm:"type:varchar(100);not_null" json:"producter"`
	TimeBegintokill int64  `gorm:"type:int;not_null" json:"time_begintokill"`
	TimeEndkill     int64  `gorm:"type:int;not_null" json:"time_endkill"`
	//因为这次考核的核心是秒杀系统，并且每个用户只买一件商品，所以将秒杀的逻辑和商品逻辑写在一起了
}*/

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(20);not null;uniqueIndex" json:"name"`
	Age      int    `gorm:"not null" json:"age"`
	Address  string `gorm:"type:varchar(100);not null" json:"address"`
	Avatar   string `gorm:"type:varchar(100);not null" json:"avatar"`
	Password string `gorm:"type:varchar(100);not null" json:"password"`
}

type Product struct {
	ID              uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string `gorm:"type:varchar(20);not null" json:"name"`
	Num             int    `json:"num"`
	Producer        string `gorm:"type:varchar(100);not null" json:"producer"`
	TimeBegintokill int64  `gorm:"not null" json:"time_begintokill"`
	TimeEndkill     int64  `gorm:"not null" json:"time_endkill"`
}

func CreateTable() {
	// 检查数据库连接是否有效
	if sql.DB == nil {
		log.Fatalf("Database connection is not initialized")
	}

	err := sql.DB.AutoMigrate(&User{}).Error
	if err != nil {
		fmt.Println("Failed to migrate:", err)
		return
	}
	err2 := sql.DB.AutoMigrate(&Product{}).Error
	if err2 != nil {
		fmt.Println("Failed to migrate:", err)
		return
	}
}
