package main

import "github.com/zzuRingo/ratestask/application"

func main() {
	app := application.GetApplicationInstance()
	app.Run()
}

