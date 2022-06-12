package virt_utils

import (
	"github.com/syunkitada/goapp2/pkg/lib/db_utils"
	"github.com/syunkitada/goapp2/pkg/lib/logger"
)

type VirtController struct {
	sqlClient *db_utils.SqlClient
}

func NewVirtContoller() (virtController *VirtController) {
	sqlClient := db_utils.NewSqlClient(&db_utils.Config{})
	return &VirtController{
		sqlClient: sqlClient,
	}
}

func (self *VirtController) Init() {
	tctx := logger.NewTraceContext()
	self.sqlClient.MustOpen(tctx)

	if tmpErr := self.sqlClient.DB.AutoMigrate(&Network{}).Error; tmpErr != nil {
		logger.Fatalf(tctx, "Failed Init: err=%s", tmpErr.Error())
	}
}
