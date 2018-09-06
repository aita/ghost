package ghost

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func DiscordHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Pong!"); err != nil {
			log.Println(err)
		}
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Ping!"); err != nil {
			log.Println(err)
		}
	}
}
