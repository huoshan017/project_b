package main

import (
	"project_b/common"
)

type IGame interface {
	GetMode() Mode
	EventMgr() common.IEventManager
}
