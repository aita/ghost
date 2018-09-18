package ghost

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"

	"github.com/aita/ghost/discord"
)

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
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		die(err)
	}

	h := discord.MakeDiscordHandler()
	dg.AddHandler(h.OnMessageCreate)

	err = dg.Open()
	if err != nil {
		die(err)
	}
	defer dg.Close()

	fmt.Println("GHOST is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
