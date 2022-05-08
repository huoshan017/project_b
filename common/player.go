package common

import (
	"project_b/common/object"

	gsnet_msg "github.com/huoshan017/gsnet/msg"
)

type PlayerState int32

const (
	PlayerNotEnter = PlayerState(0)
	PlayerEntered  = PlayerState(1)
	PlayerLeaving  = PlayerState(2)
)

// 玩家结构
type Player struct {
	id      uint64
	account string
	token   string
	state   PlayerState
	tank    *object.Tank
}

// 创建玩家
func NewPlayer(account string, id uint64) *Player {
	return &Player{account: account, id: id, state: PlayerNotEnter}
}

// ID
func (p *Player) Id() uint64 {
	return p.id
}

// 设置ID
func (p *Player) SetId(id uint64) {
	p.id = id
}

// 账号
func (p *Player) Account() string {
	return p.account
}

// 设置账号
func (p *Player) SetAccount(account string) {
	p.account = account
}

// 令牌
func (p *Player) Token() string {
	return p.token
}

// 设置令牌
func (p *Player) SetToken(token string) {
	p.token = token
}

// 已经进入游戏
func (p *Player) Entered() {
	p.state = PlayerEntered
}

// 离开游戏
func (p *Player) Left(force bool) {
	if force {
		p.state = PlayerNotEnter
	} else {
		p.state = PlayerLeaving
	}
}

// 是否已经进入
func (p *Player) IsEntered() bool {
	return p.state == PlayerEntered
}

// 是否已经离开
func (p *Player) IsLeft() bool {
	return p.state == PlayerNotEnter || p.state == PlayerLeaving
}

// 初始化坦克
//func (p *Player) InitTank(info *object.ObjStaticInfo) {
//	p.tank = object.NewTank(p.id, info)
//}

// 设置玩家坦克
func (p *Player) SetTank(tank *object.Tank) {
	p.tank = tank
}

// 获取坦克
func (p *Player) GetTank() *object.Tank {
	return p.tank
}

// 改变坦克
func (p *Player) ChangeTank(staticInfo *object.ObjStaticInfo) {
	p.tank.ChangeStaticInfo(staticInfo)
}

// 恢复坦克
func (p *Player) RestoreTank() {
	p.tank.RestoreStaticInfo()
}

// 会话玩家结构
type SPlayer struct {
	Player
	sess         *gsnet_msg.MsgSession
	disconnected bool
}

// 创建会话玩家结构
func NewSPlayer(account string, id uint64, sess *gsnet_msg.MsgSession) *SPlayer {
	return &SPlayer{
		Player: *NewPlayer(account, id),
		sess:   sess,
	}
}

// 重置会话
func (p *SPlayer) ResetSess(sess *gsnet_msg.MsgSession) {
	p.sess = sess
}

// 获得会话数据
func (p *SPlayer) GetSessData(k string) interface{} {
	return p.sess.GetData(k)
}

// 设置会话数据
func (p *SPlayer) SetSessData(k string, d interface{}) {
	p.sess.SetUserData(k, d)
}

// 离开游戏
func (p *SPlayer) Left(force bool) {
	p.Player.Left(force)
}

// 断开连接
func (p *SPlayer) Disconnect() {
	p.Left(true)
	p.sess.Close()
	p.disconnected = true
}

// todo 这里的断连表示有无调用过Disconnect函数
// 是否已断连
func (p *SPlayer) IsDisconencted() bool {
	return p.disconnected
}
