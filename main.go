package main

import (
	"commmunity/app/routes"
	"commmunity/app/zlog"
)

func main() {
	zlog.Info("程序启动")
	routes.Routes()
}
