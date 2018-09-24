package discord

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/aita/ghost/shell"
)

type Bot struct {
	sh      *shell.Shell
	session *discordgo.Session
	option  BotOption
}

type BotOption struct {
	Prefix string
}

func MakeBot(token string, option BotOption) (bot Bot, err error) {
	bot.option = option
	bot.sh = &shell.Shell{}
	bot.sh.Init()
	bot.sh.In = bytes.NewReader(nil)
	bot.sh.Out = bytes.NewBuffer(nil)

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
	if !strings.HasPrefix(msg, bot.option.Prefix) {
		return
	}
	script := msg[len(bot.option.Prefix):]
	bot.sh.Exec(script)

	buf, _ := ioutil.ReadAll(bot.sh.Out.(*bytes.Buffer))
	result := string(buf)
	if strings.TrimSpace(result) == "" {
		result = "`no output`"
	}
	if _, err := s.ChannelMessageSend(m.ChannelID, result); err != nil {
		log.Println(err)
	}
}
