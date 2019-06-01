package main

import "github.com/fidelfly/fxgos/app"

func main() {
	err := app.StartService()
	if err != nil {
		return
	}
}
