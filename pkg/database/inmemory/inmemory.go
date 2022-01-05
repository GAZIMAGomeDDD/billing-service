package inmemory

import (
	"sync"
	"time"
)

type Data struct {
	Key       string
	Value     interface{}
	Timestamp int64
}

type Heap struct {
	sync.Mutex
	data map[string]Data
}

func New() *Heap {
	return &Heap{
		data: make(map[string]Data),
	}
}

func (h *Heap) Set(key string, value interface{}, ttl int64) {
	if ttl == 0 {
		return
	}

	data := Data{
		Key:       key,
		Value:     value,
		Timestamp: time.Now().Unix(),
	}

	if ttl > 0 {
		data.Timestamp += ttl
	} else if ttl < 0 {
		data.Timestamp = -1
	}

	h.Lock()
	h.data[key] = data
	h.Unlock()
}

func (h *Heap) Get(key string) (interface{}, bool) {
	var data Data
	var val interface{}

	h.Lock()
	data, ok := h.data[key]
	h.Unlock()

	if ok {
		if data.Timestamp != -1 && data.Timestamp <= time.Now().Unix() {
			h.del(key)

			ok = false
		} else {
			val = data.Value
		}
	}

	return val, ok
}

func (h *Heap) del(key string) {
	h.Lock()
	_, ok := h.data[key]
	h.Unlock()

	if !ok {
		return
	}

	h.Lock()
	delete(h.data, key)
	h.Unlock()
}
