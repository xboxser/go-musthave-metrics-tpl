package main

import (
	"fmt"
	"metrics/internal/reset"
)

func main() {
	root := "../../" // Корневая директория проекта относительно cmd/reset
	resetService := reset.NewServiceReset(root)

	err := resetService.Run()
	if err != nil {
		panic(err)
	}

	//Проверяем работу iter22
	pool := reset.NewPool(func() *MyStruct { return &MyStruct{} })

	obj := pool.Get()
	obj.Field = 100

	pool.Put(obj)

	another := pool.Get()
	fmt.Println(another)

}

type MyStruct struct {
	Field int
}

func (m *MyStruct) Reset() {
	m.Field = 0
}
