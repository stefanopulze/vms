package event

import (
	"sync"
	"testing"
)

func TestManager_Publish(t *testing.T) {
	emit := 0
	emit2 := 0
	m := NewManager()
	m.Start()
	defer m.Close()
	var wg sync.WaitGroup

	wg.Add(2)
	m.Subscribe("test", func(event Event) {
		emit++
		wg.Done()
	})

	m.Subscribe("test2", func(event Event) {
		emit2++
		wg.Done()
	})

	m.Publish("test", "Hello")
	m.Publish("test2", "Hello")

	wg.Wait()
	if emit != 1 {
		t.Error("event not emitted")
	}

	if emit2 != 1 {
		t.Error("event not emitted")
	}
}
