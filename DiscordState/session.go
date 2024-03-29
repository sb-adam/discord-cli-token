//Package DiscordState is an abstraction layer that gives proper structs and functions to get and set the current state of the cli server
package DiscordState

import (
	"fmt"

	"github.com/adam-psycho/discordgo"
)

//!----- Session -----!//

//NewSession Creates a new Session
func NewSession(Token string) *Session {
	Session := new(Session)
	Session.Token = Token

	return Session
}

//Start attaches a discordgo listener to the Sessions and fills it.
func (Session *Session) Start() error {

	fmt.Printf("*Starting Session...")

	dg, err := discordgo.New(Session.Token)
	if err != nil {
		return err
	}

	// Open the websocket and begin listening.
	dg.Open()

	//Retrieve GuildID's from current User
	UserGuilds, err := dg.UserGuilds(100, "", "")
	if err != nil {
		return err
	}

	Session.Guilds = UserGuilds

	Session.DiscordGo = dg

	Session.User, _ = Session.DiscordGo.User("@me")

	fmt.Printf(" PASSED!\n")

	return nil
}

//NewState (constructor) attaches a new state to the Guild inside a Session, and fills it.
func (Session *Session) NewState(GuildID string, MessageAmount int) (*State, error) {
	State := new(State)

	//Disable Event Handling
	State.Enabled = false

	//Set Session
	State.Session = Session

	//Set Guild
	for _, guildID := range Session.Guilds {
		if guildID.ID == GuildID {
			Guild, err := State.Session.DiscordGo.Guild(guildID.ID)
			if err != nil {
				return nil, err
			}

			State.Guild = Guild
		}
	}

	//Retrieve Members

	State.Members = make(map[string]*discordgo.Member)

	for _, Member := range State.Guild.Members {
		State.Members[Member.User.Username] = Member
	}

	//Set MessageAmount
	State.MessageAmount = MessageAmount

	//Init Messages
	State.Messages = []*discordgo.Message{}

	//Retrieve Channels

	State.Channels = State.Guild.Channels

	return State, nil
}

//Update updates the session, this reloads the Guild list
func (Session *Session) Update() error {
	UserGuilds, err := Session.DiscordGo.UserGuilds(100, "", "")
	if err != nil {
		return err
	}

	Session.Guilds = UserGuilds
	return nil
}
