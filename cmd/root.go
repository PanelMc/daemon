package cmd

import (
	"fmt"
	"github.com/panelmc/daemon/api"
	"github.com/panelmc/daemon/config"
	"github.com/panelmc/daemon/daemon"
	"os"
	"os/signal"
	"syscall"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "panelmc",
	Short: "A brief description of your application",
}
var ExitChannel = make(chan os.Signal, 1)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		NoColors:    true,
	})
	log.SetLevel(log.TraceLevel)
	signal.Notify(ExitChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	config.Init()
	daemon.Init()

	api.Init()

	daemon.Start()

	<-ExitChannel
	fmt.Print("\nClosing the daemon...")
	os.Exit(1)
}
