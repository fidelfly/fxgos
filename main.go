package main

import "github.com/lyismydg/fxgos/app"

func main() {
	err := app.StartService()
	if err != nil {
		return
	}
}
