package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ۱. تعریف ساختار جدول (معادل SQLModel در پایتون)
type Hero struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	SecretName string `json:"secret_name"`
	Age        int    `json:"age"`
}

var db *gorm.DB

func main() {
	// ۲. اتصال به دیتابیس داکر
	dsn := "host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("❌ خطا در اتصال به دیتابیس: " + err.Error())
	}

	// ۳. ساخت خودکار جدول در دیتابیس (اگر وجود نداشته باشه)
	db.AutoMigrate(&Hero{})
	fmt.Println("✅ اتصال به دیتابیس و ساخت جدول موفقیت‌آمیز بود.")

	r := gin.Default()

	// ۴. API برای ساختن یک قهرمان جدید (ذخیره در دیتابیس)
	r.POST("/heroes", func(c *gin.Context) {
		var newHero Hero
		// گرفتن دیتای JSON از کاربر
		if err := c.ShouldBindJSON(&newHero); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ذخیره در دیتابیس
		db.Create(&newHero)
		c.JSON(http.StatusOK, newHero)
	})

	// ۵. API برای خواندن تمام قهرمان‌ها از دیتابیس
	r.GET("/heroes", func(c *gin.Context) {
		var heroes []Hero
		db.Find(&heroes)
		c.JSON(http.StatusOK, heroes)
	})

	fmt.Println("🚀 سرور گولنگ روی پورت 8080 اجرا شد...")
	r.Run(":8080")
}
