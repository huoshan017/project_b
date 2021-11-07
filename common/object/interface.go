package object

type IObject interface {
	Id() int32 // 注意：这是配置id
	SetPos(x, y int32)
	Pos() (x, y int32)
	Width() int32
	Height() int32
	Left() int32
	Right() int32
	Top() int32
	Bottom() int32
	Update()
}

type IMovableObject interface {
	IObject
}
