package common

import (
	"project_b/common/base"
	"project_b/common/ds"
	"project_b/common/log"
	"project_b/common/math"
	"project_b/common/object"
	"project_b/common/time"
	"project_b/common_data"
)

type BotStateType int

const (
	BotStateIdle      = iota
	BotStatePatrol    = 1
	BotStateAttacking = 2
)

type Bot struct {
	id           int32
	scene        *SceneLogic
	tankInstId   uint32
	searchRadius int32
	state        BotStateType
	totalTick    time.Duration
	timer        time.Duration
	enemyId      uint32
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
			break
		}
		b.state = BotStateAttacking
	case BotStateAttacking:
		enemyTank := b.scene.GetTank(b.enemyId)
		if enemyTank == nil {
			b.state = BotStateIdle
			log.Error("bot cant get enemy tank by id %v", b.enemyId)
			break
		}
		bx, by := botTank.Center()
		ex, ey := enemyTank.Center()
		dx := ex - bx
		dy := ey - by
		var dir object.Direction
		if dx < dy {
			if dy < 0 {
				dir = object.DirLeft // 敵人在左下，且X軸距離遠于Y軸
			} else {
				if dx >= 0 {
					dir = object.DirUp // 在右上，且Y軸距離遠于X軸
				} else {
					if -dx < dy {
						dir = object.DirUp // 在左上，且Y軸距離遠于X軸
					} else {
						dir = object.DirLeft // 在左上，且X軸距離遠于X軸
					}
				}
			}
		} else {
			if dy > 0 {
				dir = object.DirRight // 敵人在右上，且X軸距離遠于Y軸
			} else {
				if dx <= 0 {
					dir = object.DirDown // 在左下，且Y軸距離遠于X軸
				} else {
					if -dy > dx {
						dir = object.DirDown // 在右下，且Y軸距離大於X軸
					} else {
						dir = object.DirRight // 在右下，且X軸距離大於Y軸
					}
				}
			}
		}
		botTank.Move(dir)
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
		tcx, tcy = botTank.Center()
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
		scx, scy := searchTank.Center()
		d := (tcx-scx)*(tcx-scx) + (tcy-scy)*(tcy-scy)
		if distance < 0 || distance > d {
			distance = d
			getId = tankList[i]
		}
	}
	return getId
}

type BotManager struct {
	botList   *ds.MapListUnion[int32, *Bot]
	botPool   *base.ObjectPool[Bot]
	idCounter int32
}

func NewBotManager() *BotManager {
	return &BotManager{
		botList: ds.NewMapListUnion[int32, *Bot](),
		botPool: base.NewObjectPool[Bot](),
	}
}

func (bm *BotManager) NewBot(scene *SceneLogic, tankId uint32) *Bot {
	bm.idCounter++
	bot := bm.botPool.Get()
	bot.init(bm.idCounter, scene, tankId)
	bm.botList.Add(bm.idCounter, bot)
	return bot
}

func (bm *BotManager) RemoveBot(id int32) bool {
	bot, o := bm.botList.Get(id)
	if !o {
		return false
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
}
