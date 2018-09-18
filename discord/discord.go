package discord

import (
	"bytes"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/aita/ghost/shell"
)

type DiscordHandler struct {
	sh *shell.Shell
}

func MakeDiscordHandler() *DiscordHandler {
	sh := &shell.Shell{}
	sh.Init()
	return &DiscordHandler{
		sh: sh,
	}
}

func (h *DiscordHandler) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	msg := m.Content
	buf := bytes.NewBuffer(nil)
	h.sh.Exec(buf, msg)
	out := buf.String()
	if out != "" {
		if _, err := s.ChannelMessageSend(m.ChannelID, out); err != nil {
			log.Println(err)
		}
	}
}
