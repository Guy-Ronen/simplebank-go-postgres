package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {

	// set gin to test mode
	gin.SetMode(gin.TestMode)
	
	// run the tests
	os.Exit(m.Run())
}