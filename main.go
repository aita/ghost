package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"

	"github.com/aita/ghost/discord"
)

func init() {
	viper.SetDefault("shell.prefix", "%")
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/ghost/")
	viper.AddConfigPath("$HOME/.ghost")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		die(err)
	}

	token := viper.GetString("discord.token")
	opt := discord.BotOption{
		Prefix: viper.GetString("shell.prefix"),
	}
	bot, err := discord.NewBot(token, opt)
	if err != nil {
		die(err)
	}

	err = bot.Start()
	if err != nil {
		die(err)
	}
	defer func() {
		err := bot.Close()
		if err != nil {
			die(err)
		}
	}()

	fmt.Println("GHOST is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
