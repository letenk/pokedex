package tests

import (
	"log"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/letenk/pokedex/router"
	"github.com/letenk/pokedex/util"
	"gorm.io/gorm"
)

var ConnTest *gorm.DB
var RouteTest *gin.Engine

func TestMain(m *testing.M) {
	// Load Config
	config, err := util.LoadConfig("../.")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Open connection to postgres
	db := util.SetupDB(config.DB_SOURCE_TEST)
	ConnTest = db

	// Setup router
	RouteTest = router.SetupRouter(db)

	m.Run()
}
