package app

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	AppID    string `json:"app-id"`
	GuildID  string `json:"guild-id"`
	Token    string `json:"token"`
	RCON     string `json:"mc-rcon"`
	Password string `json:"mc-rcon-password"`
}

func GetConfig(configPath string) (Config, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

type MemberUsernameMap map[string]string

func LoadDB(dbPath string) (MemberUsernameMap, error) {
	var db MemberUsernameMap

	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("Did not find db.json file, making new one at %s", dbPath)
		_, err = os.Create(dbPath)
		if err != nil {
			log.Fatalf("Failed to create empty db.json file: %s", err.Error())
		}
		SaveDB(dbPath, db)
	}

	content, err := os.ReadFile(dbPath)
	if err != nil {
		return MemberUsernameMap{}, err
	}

	err = json.Unmarshal(content, &db)
	if err != nil {
		return MemberUsernameMap{}, err
	}
	return db, nil
}

func SaveDB(dbPath string, db MemberUsernameMap) error {
	content, err := json.Marshal(db)
	if err != nil {
		log.Fatalf("Failed to marshal db as json: %s", err.Error())
		return err
	}

	err = os.WriteFile(dbPath, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Failed to write to db.json file: %s", err.Error())
		return err
	}

	return nil
}

func simpleResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func checkMcUsernameValid(username string) bool {
	res, _ := regexp.MatchString("^[abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_]{3,16}$", username) // ^ and $ ensure that the whole string must match
	return res
}
