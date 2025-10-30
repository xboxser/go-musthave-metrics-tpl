package audit

import "metrics/internal/audit/model"

// Реализует интерфейс Publisher
type Event struct {
	observers []model.Observer
	audit     model.Audit
}

func (e *Event) Register(o model.Observer) {
	e.observers = append(e.observers, o)
}

func (e *Event) notify() {
	for _, observer := range e.observers {
		observer.Update(e.audit)
	}
}

func (e *Event) Update(audit model.Audit) {
	e.audit = audit
	e.notify()
}
