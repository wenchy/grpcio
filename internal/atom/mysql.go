package atom

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func (atom *Atom) InitMysqlDB(address, database, username, password string) error {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	// schema: dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, address, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		panic(err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(3)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(3)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// suggester: < 5m
	// refer:
	//  - https://github.com/go-sql-driver/mysql#important-settings
	//  - http://www.zhangjiee.com/blog/2020/go-mysql-closing-bad-idle-connection.html
	sqlDB.SetConnMaxLifetime(4 * time.Minute)
	atom.MysqlDB = db
	return nil
}
