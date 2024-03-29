package common

import (
	"project_b/common/base"
	"project_b/common/ds"
	"project_b/common/math"
	"project_b/common_data"
	"project_b/log"
)

type BotStateType int

const (
	BotStateIdle      = iota
	BotStatePatrol    = 1
	BotStateAttacking = 2
)

// Bot
type Bot struct {
	id                            int32
	world                         *World
	tankInstId                    uint32
	searchRadius                  int32
	state                         BotStateType
	totalTickMs                   uint32
	timerMs                       uint32
	enemyId                       uint32
	enemyGetEvent, enemyLostEvent base.Event
}

func NewBot(id int32, world *World, tankId uint32) *Bot {
	bot := &Bot{}
	bot.init(id, world, tankId)
	return bot
}

func (b *Bot) init(id int32, world *World, tankInstId uint32) {
	b.id = id
	b.world = world
	b.tankInstId = tankInstId
	b.searchRadius = common_data.DefaultSearchRadius
	b.state = BotStateIdle
	b.totalTickMs = 0
	b.timerMs = 0
	b.enemyId = 0
}

func (b *Bot) Update(tickMs uint32) {
	botTank := b.world.GetTank(b.tankInstId)
	if botTank == nil {
		b.state = BotStateIdle
		log.Error("bot cant get tank by id %v", b.tankInstId)
		return
	}
	switch b.state {
	case BotStateIdle:
		if b.totalTickMs == 0 || b.totalTickMs-b.timerMs >= 1000 {
			b.state = BotStatePatrol
			b.timerMs = b.totalTickMs
		}
	case BotStatePatrol:
		b.enemyId = b.searchEnemyTank()
		if b.enemyId == 0 {
			// 沒有找到敵人，繼續巡邏
			break
		}
		b.enemyGetEvent.Call(b.id, b.enemyId)
		b.state = BotStateAttacking
	case BotStateAttacking:
		enemyTank := b.world.GetTank(b.enemyId)
		if enemyTank == nil {
			b.state = BotStateIdle
			log.Debug("bot cant get enemy tank by id %v", b.enemyId)
			break
		}
		if botTank.IsMoving() {
			b.world.TankFire(botTank.InstId())
			break
		}
		bx, by := botTank.Pos()
		ex, ey := enemyTank.Pos()
		dx := ex - bx
		dy := ey - by
		var orientation int16
		if dx < dy {
			if dy < 0 {
				orientation = 180 // 敵人在左下，且X軸距離遠于Y軸
			} else {
				if dx >= 0 {
					orientation = 90 // 在右上，且Y軸距離遠于X軸
				} else {
					if -dx < dy {
						orientation = 90 // 在左上，且Y軸距離遠于X軸
					} else {
						orientation = 180 // 在左上，且X軸距離遠于X軸
					}
				}
			}
		} else {
			if dy > 0 {
				orientation = 0 // 敵人在右上，且X軸距離遠于Y軸
			} else {
				if dx <= 0 {
					orientation = 270 // 在左下，且Y軸距離遠于X軸
				} else {
					if -dy > dx {
						orientation = 270 // 在右下，且Y軸距離大於X軸
					} else {
						orientation = 0 // 在右下，且X軸距離大於Y軸
					}
				}
			}
		}
		angle := base.NewAngle(orientation, 0)
		botTank.Move(angle)
	}
	b.totalTickMs += tickMs
}

func (b *Bot) searchEnemyTank() uint32 {
	botTank := b.world.GetTank(b.tankInstId)
	if botTank == nil {
		return 0
	}

	var (
		distance int32 = -1
		getId    uint32
		tcx, tcy = botTank.Pos()
	)
	var rect = math.NewRectObj(tcx-b.searchRadius, tcy-b.searchRadius, 2*b.searchRadius, 2*b.searchRadius)
	tankList := b.world.GetTankListWithRange(&rect)
	for i := 0; i < len(tankList); i++ {
		searchTank := b.world.GetTank(tankList[i])
		if tankList[i] == b.tankInstId || searchTank == nil {
			continue
		}
		if searchTank.Camp() == botTank.Camp() {
			continue
		}
		scx, scy := searchTank.Pos()
		d := (tcx-scx)*(tcx-scx) + (tcy-scy)*(tcy-scy)
		if distance < 0 || distance > d {
			distance = d
			getId = tankList[i]
		}
	}
	return getId
}

func (b *Bot) clearEnemy() {
	if b.enemyId > 0 {
		b.enemyLostEvent.Call(b.id, b.enemyId)
		b.enemyId = 0
		b.state = BotStateIdle
	}
}

func (b *Bot) registerEnemyGetHandle(handle func(...any)) {
	b.enemyGetEvent.Register(handle)
}

func (b *Bot) unregisterEnemyGetHandle(handle func(...any)) {
	b.enemyGetEvent.Unregister(handle)
}

// BotManager
type BotManager struct {
	botList     *ds.MapListUnion[int32, *Bot]
	botPool     *base.ObjectPool[Bot]
	idCounter   int32
	enemyId2Bot *ds.MapListUnion[uint32, []int32]
	pause       bool
}

func NewBotManager() *BotManager {
	return &BotManager{
		botList:     ds.NewMapListUnion[int32, *Bot](),
		botPool:     base.NewObjectPool[Bot](),
		enemyId2Bot: ds.NewMapListUnion[uint32, []int32](),
	}
}

func (bm *BotManager) NewBot(world *World, tankId uint32) *Bot {
	bm.idCounter++
	bot := bm.botPool.Get()
	bot.init(bm.idCounter, world, tankId)
	bot.registerEnemyGetHandle(bm.onEmenyTankGet)
	bm.botList.Add(bm.idCounter, bot)
	return bot
}

func (bm *BotManager) RemoveBot(id int32) bool {
	bot, o := bm.botList.Get(id)
	if !o {
		return false
	}
	bm.removeBot(bot)
	bm.botList.Remove(id)
	return true
}

func (bm *BotManager) removeBot(bot *Bot) {
	bot.unregisterEnemyGetHandle(bm.onEmenyTankGet)
	bots, o := bm.enemyId2Bot.Get(bot.enemyId)
	if o {
		for i := 0; i < len(bots); i++ {
			if bots[i] == bot.id {
				bots = append(bots[:i], bots[i+1:]...)
				bm.enemyId2Bot.Set(bot.enemyId, bots)
				break
			}
		}
	}
	bm.botPool.Put(bot)
}

func (bm *BotManager) Update(tickMs uint32) {
	if bm.pause {
		return
	}
	count := bm.botList.Count()
	for i := int32(0); i < count; i++ {
		_, bot := bm.botList.GetByIndex(i)
		if bot == nil {
			continue
		}
		bot.Update(tickMs)
	}
}

func (bm *BotManager) Clear() {
	for i := int32(0); i < bm.botList.Count(); i++ {
		_, bot := bm.botList.GetByIndex(i)
		bm.removeBot(bot)
	}
	bm.botList.Clear()
	bm.idCounter = 0
	bm.enemyId2Bot.Clear()
}

func (bm *BotManager) Pause() {
	bm.pause = true
}

func (bm *BotManager) Resume() {
	bm.pause = false
}

func (bm *BotManager) onEmenyTankGet(args ...any) {
	botId := args[0].(int32)
	enemyId := args[1].(uint32)
	bots, o := bm.enemyId2Bot.Get(enemyId)
	if !o {
		bm.enemyId2Bot.Add(enemyId, []int32{botId})
	} else {
		bots = append(bots, botId)
		bm.enemyId2Bot.Set(enemyId, bots)
	}
}

func (bm *BotManager) onEnemyTankDestoryed(args ...any) {
	enemyId := args[0].(uint32)
	botIdList, o := bm.enemyId2Bot.Get(enemyId)
	if o {
		for _, botId := range botIdList {
			bot, o := bm.botList.Get(botId)
			if o {
				bot.clearEnemy()
			}
		}
		bm.enemyId2Bot.Remove(enemyId)
	}
}
