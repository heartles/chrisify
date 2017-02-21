// Declare this file to be part of the main package so it can be compiled into
// an executable.
package main

// Import all Go packages required for this file.
import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"os/exec"

	"bytes"

	"github.com/bwmarrin/discordgo"
)

// Version is a constant that stores the Disgord version information.
const Version = "v0.0.0-alpha"

// Session is declared in the global space so it can be easily used
// throughout this program.
// In this use case, there is no error that would be returned.
var Session *discordgo.Session

// Read in all configuration options from both environment variables and
// command line arguments.
func init() {

	// Discord Authentication Token
	token := os.Getenv("DG_TOKEN")
	if token == "" {
		flag.StringVar(&token, "t", "", "Discord Authentication Token")
		flag.Parse()
	}

	Session, _ = discordgo.New("Bot " + token)
}

func main() {

	// Declare any variables needed later.
	var err error

	// Print out a fancy logo!
	fmt.Printf(` 
	________  .__                               .___
	\______ \ |__| ______ ____   ___________  __| _/
	||    |  \|  |/  ___// ___\ /  _ \_  __ \/ __ | 
	||    '   \  |\___ \/ /_/  >  <_> )  | \/ /_/ | 
	||______  /__/____  >___  / \____/|__|  \____ | 
	\_______\/        \/_____/   %-16s\/`+"\n\n", Version)

	// Parse command line arguments
	flag.Parse()

	// Verify a Token was provided
	if Session.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Verify the Token is valid and grab user information
	Session.State.User, err = Session.User("@me")
	if err != nil {
		log.Printf("error fetching user information, %s\n", err)
	}

	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
	}

	Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		mentioned := false

		for _, u := range m.Mentions {
			if u.ID == s.State.User.ID {
				mentioned = true
			}
		}

		if !mentioned {
			return
		}

		for _, att := range m.Attachments {

			temp, _ := ioutil.TempFile("", "")
			defer os.Remove(temp.Name())
			defer temp.Close()

			resp, _ := http.Get(att.URL)
			io.Copy(temp, resp.Body)

			img, err := exec.Command(
				"mykify",
				"--faces=faces/mike", "--haar=haar.xml",
				temp.Name()).Output()

			if err != nil {

			}

			s.ChannelFileSendWithMessage(m.ChannelID, "Here ya go boi!",
				"png", bytes.NewReader(img))
		}
	})

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	Session.Close()

	// Exit Normally.
}
