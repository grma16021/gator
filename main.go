package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/grma16021/gator/internal/config"
	"github.com/grma16021/gator/internal/database"
	_ "github.com/lib/pq"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type state struct {
	db   *database.Queries
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

var rssfeed = RSSFeed{}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd[name] = f

}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Println("error sending req")
	}

	req.Header.Set("User-Agent", "Gator/0.1")

	res, err := client.Do(req)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading body")
	}
	defer res.Body.Close()
	xml.Unmarshal(body, &rssfeed)

	html.UnescapeString(rssfeed.Channel.Title)
	html.UnescapeString(rssfeed.Channel.Description)
	for _, item := range rssfeed.Channel.Item {
		html.UnescapeString(item.Title)
		html.UnescapeString(item.Description)
	}

	return &rssfeed, nil

}

func main() {

	conf, _ := config.Read()
	fmt.Println("config :", conf)

	db, err := sql.Open("postgres", conf.Db_url)

	dbQueries := database.New(db)

	stat := state{db: dbQueries, conf: &conf}

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
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	c := command{name: cmdName, args: cmdArgs}

	err = cmds.run(&stat, c)
	if err != nil {
		fmt.Println("error ", err)
	}
	conf, _ = config.Read()
	fmt.Println("config user :", conf.Current_user_name)
	fmt.Println("config db :", conf.Db_url)

}

func handlerRegister(s *state, cmd command) error {
	if cmd.args == nil {
		return fmt.Errorf("The login handler expects a username")
	}
	userName := cmd.args[0]
	fuser, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		fmt.Errorf("error searching for user")
	}
	if fuser.Name == userName {
		os.Exit(1)
	}
	user := database.CreateUserParams{}
	id := uuid.New()

	user.ID = id.String()
	user.Name = cmd.args[0]
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	cUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		fmt.Println("error creating user ", err)
	}
	fmt.Println("created user: " + cUser.Name)

	s.conf.Current_user_name = cmd.args[0]
	s.conf.SetUser(s.conf.Current_user_name)
	fmt.Println("User with name: " + userName + " Has been created")

	return nil
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		fmt.Errorf("The login handler expects a username")
		os.Exit(1)
	}
	uName, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Errorf("error fetching user", err)

	}
	if uName.Name == "" {
		fmt.Println("user does not exist")
		os.Exit(1)
	}

	s.conf.Current_user_name = uName.Name
	s.conf.SetUser(s.conf.Current_user_name)

	fmt.Println("The username has been set to: " + cmd.args[0])

	return nil
}

func handlerReset(s *state, cmd command) error {

	fmt.Println("Deleting users")
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Println("Error deleting users", err)
	}
	fmt.Println("Deleted all users")

	return nil
}

func handlerUsers(s *state, cmd command) error {

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("Error getting users")
		os.Exit(1)
	}

	for _, user := range users {

		if user.Name == s.conf.Current_user_name {
			fmt.Println("- " + user.Name + " (current)")
		} else {
			fmt.Println("- " + user.Name)
		}

	}

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		os.Exit(1)
	}

	user, err := s.db.GetUser(context.Background(), s.conf.Current_user_name)
	if err != nil {
		fmt.Errorf("error fetching user", err)
	}

	feed := database.CreateFeedParams{}
	feed.ID = uuid.NewString()
	feed.UserID = user.ID
	feed.Name.String = cmd.name
	feed.Url.String = cmd.args[0]
	feed.CreatedAt = time.Now()
	feed.UpdatedAt = time.Now()
	createdFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return err
	}

	fmt.Println(createdFeed)
	return nil
}

func handlerAgg(s *state, cmd command) error {
	fmt.Println("Fetching stuff")
	//fmt.Println(cmd.args[0])
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(feed)

	return nil
}
