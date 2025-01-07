package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"takeout/config"
	"testing"
	"time"
)

type Dish struct {
	Id          uint64    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name        string    `json:"name"`
	CategoryId  uint64    `json:"categoryId"`
	Price       float64   `json:"price"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
	CreateUser  uint64    `json:"createUser"`
	UpdateUser  uint64    `json:"updateUser"`
	// 一对多
	Flavors []DishFlavor `json:"flavors"`
}

func (d *Dish) TableName() string {
	return "dish"
}

type DishFlavor struct {
	Id     uint64 `json:"id"`      //口味id
	DishId uint64 `json:"dish_id"` //菜品id
	Name   string `json:"name"`    //口味主题 温度|甜度|辣度
	Value  string `json:"value"`   //口味信息 可多个
}

func (d *DishFlavor) TableName() string {
	return "dish_flavor"
}

var (
	Config = config.InitLoadConfig() // 需要手动在dao目录下添加config目录并放入yaml配置文件
	d      = InitDatabase(Config.DataSource.Dsn())
	dish   Dish
)

func InitDatabase(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)  //设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) //打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	return db
}

func BenchmarkPreload(b *testing.B) {
	// 每次测试时重新初始化 dish
	for i := 0; i < b.N; i++ {
		dish.Id = 46

		err := d.Preload("Flavors").Find(&dish).Error
		if err != nil {
			b.Fatalf("Error fetching dish: %v", err)
		}
	}
}

func BenchmarkTwoQueries(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dish.Id = 46

		err := d.First(&dish).Error
		if err != nil {
			b.Fatalf("Error fetching dish: %v", err)
		}
		err = d.Where("dish_id = ?", dish.Id).Find(&dish.Flavors).Error
		if err != nil {
			b.Fatalf("Error fetching dish: %v", err)
		}
	}
}

// BenchmarkPreload-8      	    2346	    492142 ns/op
// BenchmarkPreload-8：表示 Preload 方法每次操作平均耗时 492142 纳秒 (ns)，运行了 2346 次。

// BenchmarkTwoQueries-8   	    2499	    511833 ns/op
// BenchmarkTwoQueries-8：表示注释中的两次查询方法每次操作平均耗时 511833 纳秒 (ns)，运行了 2499 次。

// Preload 方法的性能再次稍好，平均每次操作的耗时少了 19701 纳秒（约 4% 的差距）。
// 在这两个基准测试中，Preload 比 TwoQueries 更高效，查询的执行时间更短。
// 结果显示，Preload 方法在性能上明显优于 TwoQueries，但差距大约是 4%，这表明 Preload 方法的优化效果较为显著。
