package util

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlClient *MySQLClient

type MySQLClientConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	Charset  string
}

func NewMySQLClientConfig(dbName string) *MySQLClientConfig {
	charset := os.Getenv("MYSQL_CHARSET")
	if charset == "" {
		charset = "utf8mb4"
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}

	return &MySQLClientConfig{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     port,
		DBName:   dbName,
		Charset:  charset,
	}
}

func (c *MySQLClientConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset,
	)
}

type MySQLClient struct {
	SqlDB  *sql.DB
	GormDB *gorm.DB
}

func NewMySQLClient() (*MySQLClient, error) {

	config := NewMySQLClientConfig(os.Getenv("MYSQL_USER_DBNAME"))
	dsn := config.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("open db failed", "error", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("get sql db failed", "error", err)
		return nil, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	if err := sqlDB.Ping(); err != nil {
		slog.Error("ping db failed", "error", err)
		_ = sqlDB.Close()
		return nil, err
	}

	mysqlClient = &MySQLClient{SqlDB: sqlDB, GormDB: db}

	return mysqlClient, nil
}

func GetMySQLClient() *MySQLClient {
	return mysqlClient
}
