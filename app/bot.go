package app

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/gorcon/rcon"
)

var s *discordgo.Session
var config Config
var dbPath string
var db MemberUsernameMap
var conn *rcon.Conn

func Run() {
	// parse CLI options and get config
	configPath := *(flag.String("config", "config.json", "path to config.json file"))
	dbPath = *(flag.String("db", "db.json", "path to db.json file"))
	flag.Parse()

	// load config
	config, err := GetConfig(configPath)
	if err != nil {
		log.Fatalf("Error when reading config file: %s", err)
	}

	// load db or create new empty db
	db, err = LoadDB(dbPath)
	if err != nil {
		log.Fatalf("Error when reading config file: %s", err)
	}

	// open rconn connection
	conn, err = rcon.Dial(config.RCON, config.Password)
	if err != nil {
		log.Fatalf("Failed to open rcon connection: %s", err.Error())
	}

	// create new session
	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	// attach slash command listener
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// attach session connected listener
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// overwrite registered commands
	log.Println("Bulk overwriting commands")
	_, err = s.ApplicationCommandBulkOverwrite(config.AppID, config.GuildID, commands)
	if err != nil {
		log.Fatalf("Failed to bulk overwrite commands: %s", err.Error())
	}

	// open session
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer s.Close()
	defer conn.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
