package main

// system packages
import (
	"os"
	"log"
	"encoding/json"
  "bytes"
)

type Config struct{
  BotTokenKey string
  GuildID string
}

// root -> config file
const CONFIG_FILE string = "./config.json"

// max data size
const CONFIG_FILE_MAX_SIZE int = 1024

func GetBotConfiguration(config *Config){
  rawData := make([]byte, CONFIG_FILE_MAX_SIZE)
  
	// open configuration file
	file, err := os.Open(CONFIG_FILE)
  if err != nil {
		log.Fatalf("Failed to open file '%s' [%v]", CONFIG_FILE, err)
	}
	defer file.Close()

	// read all file data
	_, err = file.Read(rawData)
	if err != nil {
		log.Fatalf("Failed to read file '%s' [%v]", CONFIG_FILE, err)
	}

  cleanData := bytes.Trim(rawData, "\x00")

  // convert byte data to a struct
  err = json.Unmarshal(cleanData, config)
  if err != nil {
    log.Fatalf("Failed converting file data to structure data [%v]", err)
  }

  // output the data
	log.Printf("Outputting configuration data\n" +
  "                     -> BotTokenKey [%s]\n" +
  "                     -> GuildID [%s]",
   config.BotTokenKey,
   config.GuildID)
}
