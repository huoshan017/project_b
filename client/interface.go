package main

import (
	"project_b/common/base"
)

type IGame interface {
	GetState() GameState
	EventMgr() base.IEventManager
}
