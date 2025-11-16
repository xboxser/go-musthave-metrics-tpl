package main

import "metrics/reset"

func main() {
	root := "../../" // Корневая директория проекта относительно cmd/reset
	resetService := reset.NewServiceReset(root)

	err := resetService.Run()
	if err != nil {
		panic(err)
	}

}
