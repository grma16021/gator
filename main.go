package main

import (
	"fmt"
	"os"

	"github.com/grma16021/gator/internal/config"
)

type state struct {
	conf *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmd map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmad command) error {

	handler, ok := c.cmd[cmad.name]
	if !ok {
		return fmt.Errorf("Unknown command")
	}

	return handler(s, cmad)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd[name] = f

}

func main() {

	conf, _ := config.Read()
	fmt.Println("config :", conf)

	stat := state{conf: &conf}

	cmds := commands{
		cmd: map[string]func(*state, command) error{},
	}

	args := os.Args
	if len(args) <= 1 {
		fmt.Println("Command not provided")
		os.Exit(1)
	}
	cmdName := args[1]
	if cmdName == "" {
		fmt.Errorf("Name not provided")
		os.Exit(1)
	}
	cmdArgs := args[2:]

	cmds.register("login", handlerLogin)
	c := command{name: cmdName, args: cmdArgs}

	err := cmds.run(&stat, c)
	if err != nil {
		fmt.Println("error ", err)
	}
	conf, _ = config.Read()
	fmt.Println("config user :", conf.Current_user_name)
	fmt.Println("config db :", conf.Db_url)

}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		fmt.Errorf("The login handler expects a username")
		os.Exit(1)
	}

	s.conf.Current_user_name = cmd.args[0]
	s.conf.SetUser(s.conf.Current_user_name)

	fmt.Println("The username has been set to: " + cmd.args[0])

	return nil
}
