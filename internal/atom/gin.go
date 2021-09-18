package atom

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func (atom *Atom) InitGin() error {
	atom.GinEngine = gin.New()
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	atom.GinEngine.Use(ginzap.Ginzap(zaplogger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	atom.GinEngine.Use(ginzap.RecoveryWithZap(zaplogger, true))
	return nil
}
