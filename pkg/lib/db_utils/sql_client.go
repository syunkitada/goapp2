package db_utils

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
	"github.com/syunkitada/goapp2/pkg/lib/struct_utils"
)

type Config struct {
	Driver     string
	Connection string
}

var conf = Config{
	Driver:     "sqlite3",
	Connection: "/tmp/sqlite3.db",
}

type SqlClient struct {
	DB   *gorm.DB
	conf Config
}

func NewSqlClient(conf2 *Config) (client *SqlClient) {
	struct_utils.MergeStruct(conf, conf2)
	client = &SqlClient{
		conf: conf,
	}
	return
}

func (self *SqlClient) MustOpen(tctx *logger.TraceContext) {
	if db, tmpErr := gorm.Open("sqlite3", self.conf.Connection); tmpErr != nil {
		logger.Fatalf(tctx, "Failed Open: err=%s", tmpErr.Error())
	} else {
		self.DB = db
		return
	}

	if db, tmpErr := gorm.Open("mysql", self.conf.Connection); tmpErr != nil {
		// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		logger.Fatalf(tctx, "Failed Open: err=%s", tmpErr.Error())
	} else {
		// db.LogMode(self.conf.EnableDatabaseLog)
		self.DB = db
		return
	}
}

func (self *SqlClient) MustClose(tctx *logger.TraceContext) {
	if tmpErr := self.DB.Close(); tmpErr != nil {
		logger.Fatalf(tctx, "Failed Close: err=%s", tmpErr.Error())
	}
}

func (self *SqlClient) Transact(tctx *logger.TraceContext, txFunc func(tx *gorm.DB) (err error)) (err error) {
	tx := self.DB.Begin()
	if err = tx.Error; err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			if tmpErr := tx.Rollback().Error; tmpErr != nil {
				logger.Errorf(tctx, "failed rollback on recover: err=%s", tmpErr.Error())
			} else {
				logger.Infof(tctx, "transaction was rollbacked because of recover")
				err = fmt.Errorf("transaction was rollbacked because of recover")
			}
		} else if err != nil {
			if tmpErr := tx.Rollback().Error; tmpErr != nil {
				logger.Errorf(tctx, "failed rollback after err: err=%s", tmpErr.Error())
			} else {
				logger.Infof(tctx, "transaction was rollbacked because of err:: err=%s", err.Error())
			}
		} else {
			if err = tx.Commit().Error; err != nil {
				if tmpErr := tx.Rollback().Error; tmpErr != nil {
					logger.Errorf(tctx, "failed rollback on commit: err=%s", tmpErr.Error())
				} else {
					logger.Infof(tctx, "transaction was rollbacked because of commit err:: err=%s", err.Error())
				}
			}
		}
	}()
	err = txFunc(tx)
	return
}
