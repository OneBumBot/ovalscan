package main

import (
	"fmt"
	"log"

	"github.com/onebumbot/ovalscan/internal/config"
	"github.com/onebumbot/ovalscan/internal/ssh"
)

func main() {
	fmt.Println("Welcome to ovalscan!")

	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal("unable to load config.toml: ", err)
	}

	client, err := ssh.GetClient(cfg)

	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

	session, err := ssh.CreateSession(client)
	if err != nil {
		log.Fatal("failed to create ssh session: ", err)
	}
	defer session.Close()

	if err := session.Run("touch test"); err != nil {
		log.Fatal("failed to create remote file: ", err)
	}

	out, err := ssh.ExecuteCommand(session, "touch test1", "touch test2")

	if err != nil {
		log.Fatal("failed execution commands: ", err)
	}

	fmt.Print(string(out))
}
