package model

// интерфейсы издатель/субъекта
type Publisher interface {
	register(Observer) // Регистрация наблюдателя
	notify()           // Оповещение
}

// Интерфейс наблюдателя
type Observer interface {
	Update(Audit)
}
