package startup

import (
	"github.com/google/wire"
)

var DBProvider = wire.NewSet(InitDB)
