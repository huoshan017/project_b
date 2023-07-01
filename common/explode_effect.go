package common

import (
	"project_b/common/math"
	"project_b/common/object"
)

type BulletExplodeEffect struct {
	object.Effect
}

type TankExplodeEffect struct {
	object.Effect
}

func bulletExplodeEffect(args ...any) {
	pmap := args[0].(*PartitionMap)
	bullet := args[1].(*object.Bullet)
	pmap.GetLayerObjsWithRange(math.NewRect(bullet.Left(), bullet.Bottom(), bullet.Width(), bullet.Height()))
}

func bigBulletExplodeEffect(args ...any) {

}
