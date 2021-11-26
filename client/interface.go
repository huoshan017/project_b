package main

import (
	"project_b/common/base"
)

type IGame interface {
	GetMode() Mode
	EventMgr() base.IEventManager
}
