package base

import "reflect"

type EventHandleT[T any] func(t T)
type EventHandleT2[T1 any, T2 any] func(t1 T1, t2 T2)
type EventHandleT3[T1 any, T2 any, T3 any] func(t1 T1, t2 T2, t3 T3)
type EventHandleT4[T1 any, T2 any, T3 any, T4 any] func(t1 T1, t2 T2, t3 T3, t4 T4)
type EventHandleT5[T1 any, T2 any, T3 any, T4 any, T5 any] func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5)
type EventHandleT6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any] func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6)
type EventHandleT7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any] func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7)
type EventHandleT8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any] func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8)
type EventHandleT9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any] func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9)

type EventT[T any] struct {
	handles []EventHandleT[T]
}

func NewEventT[T any]() *EventT[T] {
	return &EventT[T]{}
}

func (e *EventT[T]) Register(handle EventHandleT[T]) {
	for _, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			return
		}
	}
	e.handles = append(e.handles, handle)
}

func (e *EventT[T]) Unregister(handle EventHandleT[T]) {
	for i, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			e.handles = append(e.handles[:i], e.handles[i+1:]...)
			break
		}
	}
}

func (e *EventT[T]) Call(t T) {
	for _, h := range e.handles {
		if h != nil {
			h(t)
		}
	}
}

type EventT2[T1 any, T2 any] struct {
	handles []EventHandleT2[T1, T2]
}

func NewEventT2[T1 any, T2 any]() *EventT2[T1, T2] {
	return &EventT2[T1, T2]{}
}

func (e *EventT2[T1, T2]) Register(handle EventHandleT2[T1, T2]) {
	for _, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			return
		}
	}
	e.handles = append(e.handles, handle)
}

func (e *EventT2[T1, T2]) Unregister(handle EventHandleT2[T1, T2]) {
	for i, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			e.handles = append(e.handles[:i], e.handles[i+1:]...)
			break
		}
	}
}

func (e *EventT2[T1, T2]) Call(t1 T1, t2 T2) {
	for _, h := range e.handles {
		if h != nil {
			h(t1, t2)
		}
	}
}

type EventT3[T1 any, T2 any, T3 any] struct {
	handles []EventHandleT3[T1, T2, T3]
}

func NewEventT3[T1 any, T2 any, T3 any]() *EventT3[T1, T2, T3] {
	return &EventT3[T1, T2, T3]{}
}

func (e *EventT3[T1, T2, T3]) Register(handle EventHandleT3[T1, T2, T3]) {
	for _, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			return
		}
	}
	e.handles = append(e.handles, handle)
}

func (e *EventT3[T1, T2, T3]) Unregister(handle EventHandleT3[T1, T2, T3]) {
	for i, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			e.handles = append(e.handles[:i], e.handles[i+1:]...)
			break
		}
	}
}

func (e *EventT3[T1, T2, T3]) Call(t1 T1, t2 T2, t3 T3) {
	for _, h := range e.handles {
		if h != nil {
			h(t1, t2, t3)
		}
	}
}

type EventT4[T1 any, T2 any, T3 any, T4 any] struct {
	handles []EventHandleT4[T1, T2, T3, T4]
}

func NewEventT4[T1 any, T2 any, T3 any, T4 any]() *EventT4[T1, T2, T3, T4] {
	return &EventT4[T1, T2, T3, T4]{}
}

func (e *EventT4[T1, T2, T3, T4]) Register(handle EventHandleT4[T1, T2, T3, T4]) {
	for _, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			return
		}
	}
	e.handles = append(e.handles, handle)
}

func (e *EventT4[T1, T2, T3, T4]) Unregister(handle EventHandleT4[T1, T2, T3, T4]) {
	for i, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			e.handles = append(e.handles[:i], e.handles[i+1:]...)
			break
		}
	}
}

func (e *EventT4[T1, T2, T3, T4]) Call(t1 T1, t2 T2, t3 T3, t4 T4) {
	for _, h := range e.handles {
		if h != nil {
			h(t1, t2, t3, t4)
		}
	}
}

type EventT5[T1 any, T2 any, T3 any, T4 any, T5 any] struct {
	handles []EventHandleT5[T1, T2, T3, T4, T5]
}

func NewEventT5[T1 any, T2 any, T3 any, T4 any, T5 any]() *EventT5[T1, T2, T3, T4, T5] {
	return &EventT5[T1, T2, T3, T4, T5]{}
}

func (e *EventT5[T1, T2, T3, T4, T5]) Register(handle EventHandleT5[T1, T2, T3, T4, T5]) {
	for _, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			return
		}
	}
	e.handles = append(e.handles, handle)
}

func (e *EventT5[T1, T2, T3, T4, T5]) Unregister(handle EventHandleT5[T1, T2, T3, T4, T5]) {
	for i, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			e.handles = append(e.handles[:i], e.handles[i+1:]...)
			break
		}
	}
}

func (e *EventT5[T1, T2, T3, T4, T5]) Call(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) {
	for _, h := range e.handles {
		if h != nil {
			h(t1, t2, t3, t4, t5)
		}
	}
}

type EventT6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any] struct {
	handles []EventHandleT6[T1, T2, T3, T4, T5, T6]
}

func NewEventT6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any]() *EventT6[T1, T2, T3, T4, T5, T6] {
	return &EventT6[T1, T2, T3, T4, T5, T6]{}
}

func (e *EventT6[T1, T2, T3, T4, T5, T6]) Register(handle EventHandleT6[T1, T2, T3, T4, T5, T6]) {
	for _, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			return
		}
	}
	e.handles = append(e.handles, handle)
}

func (e *EventT6[T1, T2, T3, T4, T5, T6]) Unregister(handle EventHandleT6[T1, T2, T3, T4, T5, T6]) {
	for i, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			e.handles = append(e.handles[:i], e.handles[i+1:]...)
			break
		}
	}
}

func (e *EventT6[T1, T2, T3, T4, T5, T6]) Call(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) {
	for _, h := range e.handles {
		if h != nil {
			h(t1, t2, t3, t4, t5, t6)
		}
	}
}

type EventT7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any] struct {
	handles []EventHandleT7[T1, T2, T3, T4, T5, T6, T7]
}

func NewEventT7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any]() *EventT7[T1, T2, T3, T4, T5, T6, T7] {
	return &EventT7[T1, T2, T3, T4, T5, T6, T7]{}
}

func (e *EventT7[T1, T2, T3, T4, T5, T6, T7]) Register(handle EventHandleT7[T1, T2, T3, T4, T5, T6, T7]) {
	for _, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			return
		}
	}
	e.handles = append(e.handles, handle)
}

func (e *EventT7[T1, T2, T3, T4, T5, T6, T7]) Unregister(handle EventHandleT7[T1, T2, T3, T4, T5, T6, T7]) {
	for i, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			e.handles = append(e.handles[:i], e.handles[i+1:]...)
			break
		}
	}
}

func (e *EventT7[T1, T2, T3, T4, T5, T6, T7]) Call(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) {
	for _, h := range e.handles {
		if h != nil {
			h(t1, t2, t3, t4, t5, t6, t7)
		}
	}
}

type EventT8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any] struct {
	handles []EventHandleT8[T1, T2, T3, T4, T5, T6, T7, T8]
}

func NewEventT8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any]() *EventT8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return &EventT8[T1, T2, T3, T4, T5, T6, T7, T8]{}
}

func (e *EventT8[T1, T2, T3, T4, T5, T6, T7, T8]) Register(handle EventHandleT8[T1, T2, T3, T4, T5, T6, T7, T8]) {
	for _, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			return
		}
	}
	e.handles = append(e.handles, handle)
}

func (e *EventT8[T1, T2, T3, T4, T5, T6, T7, T8]) Unregister(handle EventHandleT8[T1, T2, T3, T4, T5, T6, T7, T8]) {
	for i, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			e.handles = append(e.handles[:i], e.handles[i+1:]...)
			break
		}
	}
}

func (e *EventT8[T1, T2, T3, T4, T5, T6, T7, T8]) Call(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) {
	for _, h := range e.handles {
		if h != nil {
			h(t1, t2, t3, t4, t5, t6, t7, t8)
		}
	}
}

type EventT9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any] struct {
	handles []EventHandleT9[T1, T2, T3, T4, T5, T6, T7, T8, T9]
}

func NewEventT9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any]() *EventT9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return &EventT9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{}
}

func (e *EventT9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Register(handle EventHandleT9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) {
	for _, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			return
		}
	}
	e.handles = append(e.handles, handle)
}

func (e *EventT9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Unregister(handle EventHandleT9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) {
	for i, h := range e.handles {
		if reflect.DeepEqual(h, handle) {
			e.handles = append(e.handles[:i], e.handles[i+1:]...)
			break
		}
	}
}

func (e *EventT9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Call(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) {
	for _, h := range e.handles {
		if h != nil {
			h(t1, t2, t3, t4, t5, t6, t7, t8, t9)
		}
	}
}
