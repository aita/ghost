package discord

import (
	"bytes"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/aita/ghost/shell"
)

type Bot struct {
	sh      *shell.Shell
	session *discordgo.Session
}

func MakeBot(token string) (bot Bot, err error) {
	bot.sh = &shell.Shell{}
	bot.sh.Init()

	bot.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return
	}

	bot.session.AddHandler(bot.OnMessageCreate)
	return
}

func (bot Bot) Start() error {
	return bot.session.Open()

}

func (bot Bot) Close() error {
	return bot.session.Close()
}

func (bot Bot) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	msg := m.Content
	buf := bytes.NewBuffer(nil)
	bot.sh.Exec(buf, msg)
	out := buf.String()
	if strings.TrimSpace(out) == "" {
		out = "`no output`"
	}
	if _, err := s.ChannelMessageSend(m.ChannelID, out); err != nil {
		log.Println(err)
	}
}
