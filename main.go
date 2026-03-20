package main

import (
	"boilerplate/config"
	"boilerplate/internal/middleware"
	"boilerplate/internal/model"
	"boilerplate/internal/router"
	"boilerplate/pkg/mysql"
	"boilerplate/pkg/redis"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func setup() {
	config.Init()
	mysql.Init()
	redis.Init()
}

func b() {
	fmt.Println("b before")
	defer func() { fmt.Println("b after") }()
}
func a() {
	fmt.Println("a before")
	b()
	fmt.Println("a after")
}

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	if 1 == 1 {
		a()
		fmt.Printf("%#v", model.EbLotteryRule{})
		return
	}
	setup()
	if os.Getenv("ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware(), middleware.RequestIDMiddleware())
	router.SetupRouter(r)
}
