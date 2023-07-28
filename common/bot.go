package common

import (
	"project_b/common/base"
	"project_b/common/ds"
	"project_b/common/log"
	"project_b/common/math"
	"project_b/common/time"
	"project_b/common_data"
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
	scene                         *SceneLogic
	tankInstId                    uint32
	searchRadius                  int32
	state                         BotStateType
	totalTick                     time.Duration
	timer                         time.Duration
	enemyId                       uint32
	enemyGetEvent, enemyLostEvent base.Event
}

func NewBot(id int32, scene *SceneLogic, tankId uint32) *Bot {
	bot := &Bot{}
	bot.init(id, scene, tankId)
	return bot
}

func (b *Bot) init(id int32, scene *SceneLogic, tankInstId uint32) {
	b.id = id
	b.scene = scene
	b.tankInstId = tankInstId
	b.searchRadius = common_data.DefaultSearchRadius
	b.state = BotStateIdle
	b.totalTick = 0
	b.timer = 0
	b.enemyId = 0
}

func (b *Bot) Update(tick time.Duration) {
	botTank := b.scene.GetTank(b.tankInstId)
	if botTank == nil {
		b.state = BotStateIdle
		log.Error("bot cant get tank by id %v", b.tankInstId)
		return
	}
	switch b.state {
	case BotStateIdle:
		if b.totalTick == 0 || b.totalTick-b.timer >= time.Second {
			b.state = BotStatePatrol
			b.timer = b.totalTick
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
		enemyTank := b.scene.GetTank(b.enemyId)
		if enemyTank == nil {
			b.state = BotStateIdle
			log.Debug("bot cant get enemy tank by id %v", b.enemyId)
			break
		}
		if botTank.IsMoving() {
			b.scene.TankFire(botTank.InstId(), 1)
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
	b.totalTick += tick
}

func (b *Bot) searchEnemyTank() uint32 {
	botTank := b.scene.GetTank(b.tankInstId)
	if botTank == nil {
		return 0
	}

	var (
		distance int32 = -1
		getId    uint32
		tcx, tcy = botTank.Pos()
	)
	var rect = math.NewRectObj(tcx-b.searchRadius, tcy-b.searchRadius, 2*b.searchRadius, 2*b.searchRadius)
	tankList := b.scene.GetTankListWithRange(&rect)
	for i := 0; i < len(tankList); i++ {
		searchTank := b.scene.GetTank(tankList[i])
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
}

func NewBotManager() *BotManager {
	return &BotManager{
		botList:     ds.NewMapListUnion[int32, *Bot](),
		botPool:     base.NewObjectPool[Bot](),
		enemyId2Bot: ds.NewMapListUnion[uint32, []int32](),
	}
}

func (bm *BotManager) NewBot(scene *SceneLogic, tankId uint32) *Bot {
	bm.idCounter++
	bot := bm.botPool.Get()
	bot.init(bm.idCounter, scene, tankId)
	bot.registerEnemyGetHandle(bm.onEmenyTankGet)
	bm.botList.Add(bm.idCounter, bot)
	return bot
}

func (bm *BotManager) RemoveBot(id int32) bool {
	bot, o := bm.botList.Get(id)
	if !o {
		return false
	}
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
	bm.botList.Remove(id)
	bm.botPool.Put(bot)
	return true
}

func (bm *BotManager) Update(tick time.Duration) {
	count := bm.botList.Count()
	for i := int32(0); i < count; i++ {
		_, bot := bm.botList.GetByIndex(i)
		if bot == nil {
			continue
		}
		bot.Update(tick)
	}
}

func (bm *BotManager) Clear() {
	bm.botList.Clear()
	bm.idCounter = 0
	bm.enemyId2Bot.Clear()
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
