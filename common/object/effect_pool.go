package object

type EffectPool struct {
	effects   []*Effect
	emap      map[*Effect]struct{}
	idCounter uint32
}

func NewEffectPool() *EffectPool {
	return &EffectPool{
		emap: make(map[*Effect]struct{}),
	}
}

func (ep *EffectPool) Get(staticInfo *EffectStaticInfo, effectFunc func(...any), args ...any) *Effect {
	var (
		e *Effect
		l = len(ep.effects)
	)
	ep.idCounter += 1
	if l > 0 {
		e = ep.effects[l-1]
		e.init(ep.idCounter, staticInfo, effectFunc, args...)
		ep.effects = ep.effects[:l-1]
	} else {
		e = NewEffect(ep.idCounter, staticInfo, effectFunc, args...)
	}
	ep.emap[e] = struct{}{}
	return e
}

func (ep *EffectPool) Put(e *Effect) bool {
	if _, o := ep.emap[e]; !o {
		return false
	}
	e.uninit()
	ep.effects = append(ep.effects, e)
	delete(ep.emap, e)
	return true
}

func (ep *EffectPool) Clear() {
	for _, e := range ep.effects {
		e.uninit()
	}
	ep.effects = ep.effects[:0]
	clear(ep.emap)
	ep.idCounter = 0
}
