package core

import (
	"net/http"
	"sync"
)

type (
	logLevel int8
	// MapData map like any value datatype
	MapData map[string]interface{}
	// Garden go garden framework class
	Garden struct {
		container      sync.Map
		services       map[string]*service
		cfg            cfg
		serviceManager chan serviceOperate
	}
)

const (
	httpOk       = http.StatusOK
	httpFail     = http.StatusInternalServerError
	httpNotFound = http.StatusNotFound

	infoSuccess       = "Success"
	infoServerError   = "Server Error"
	infoServerLimiter = "Server limit flow"
	infoServerFusing  = "Server fusing flow"
	infoNoAuth        = "No access permission"
	infoNotFound      = "The resource could not be found"
	infoTimeout       = "Request timeout"
)
