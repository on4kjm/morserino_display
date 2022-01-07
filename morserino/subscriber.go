package morserino

import "sync"

type EventDispatcher struct {
	handlersMu sync.RWMutex
	handlers   []EventHandler
}

func (e *EventDispatcher) Handle(event Event) error {
	e.handlersMu.RLock()
	defer e.handlersMu.RUnlock()

	for _, hdl := range e.handlers {
		if err := hdl.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

func (e *EventDispatcher) Register(handler EventHandler) {
	e.handlersMu.Lock()
	defer e.handlersMu.Unlock()

	e.handlers = append(e.handlers, handler)
}

func (e *EventDispatcher) Unregister(handler EventHandler) {
	idx, ok := e.indexOf(handler)

	if !ok {
		return
	}

	e.handlersMu.Lock()
	defer e.handlersMu.Unlock()

	e.handlers[idx] = e.handlers[len(e.handlers)-1]
	e.handlers = e.handlers[:len(e.handlers)-1]
}

func (e *EventDispatcher) indexOf(handler EventHandler) (int, bool) {
	e.handlersMu.RLock()
	defer e.handlersMu.RUnlock()

	for i, hdl := range e.handlers {
		if hdl == handler {
			return i, true
		}
	}

	return 0, false
}
