package main

import (
	"fmt"

	"github.com/grma16021/gator/internal/config"
)

func main() {

	conf, _ := config.Read()
	fmt.Println("config :", conf)

	conf.SetUser("Mathias")

	conf, _ = config.Read()
	fmt.Println("config user :", conf.Current_user_name)
	fmt.Println("config db :", conf.Db_url)
}
