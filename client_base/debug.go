package client_base

type Debug struct {
	showTankBoundingbox              bool
	showTankAABB                     bool
	showShellBoundingbox             bool
	showShellAABB                    bool
	showMapGrid                      bool
	showTankCollisionDetectionRegion bool
}

func (d Debug) IsShowTankBoundingbox() bool {
	return d.showTankBoundingbox
}

func (d *Debug) ShowTankBoundingbox() {
	d.showTankBoundingbox = true
}

func (d *Debug) HideTankBoundingbox() {
	d.showTankBoundingbox = false
}

func (d *Debug) IsShowTankAABB() bool {
	return d.showTankAABB
}

func (d *Debug) ShowTankAABB() {
	d.showTankAABB = true
}

func (d *Debug) HideTankAABB() {
	d.showTankAABB = false
}

func (d Debug) IsShowShellBoundingbox() bool {
	return d.showShellBoundingbox
}

func (d *Debug) ShowShellBoundingbox() {
	d.showShellBoundingbox = true
}

func (d *Debug) HideShellBoundingbox() {
	d.showShellBoundingbox = false
}

func (d Debug) IsShowShellAABB() bool {
	return d.showShellAABB
}

func (d *Debug) ShowShellAABB() {
	d.showShellAABB = true
}

func (d *Debug) HideShellAABB() {
	d.showShellAABB = false
}

func (d Debug) IsShowMapGrid() bool {
	return d.showMapGrid
}

func (d *Debug) ShowMapGrid() {
	d.showMapGrid = true
}

func (d *Debug) HideMapGrid() {
	d.showMapGrid = false
}

func (d Debug) IsShowTankCollisionDetectionRegion() bool {
	return d.showTankCollisionDetectionRegion
}

func (d *Debug) ShowTankCollisionDetectionRegion() {
	d.showTankCollisionDetectionRegion = true
}

func (d *Debug) HideTankCollisionDetectionRegion() {
	d.showTankCollisionDetectionRegion = false
}

var debug Debug

func GetDebug() *Debug {
	return &debug
}
