package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/panelmc/daemon/api"
	"github.com/panelmc/daemon/daemon"
	"github.com/panelmc/daemon/infra"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var exitChan = make(chan os.Signal, 1)

var rootCmd = &cobra.Command{
	Use:   "panelmc",
	Short: "Starts the daemon master node",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func main() {
	infra.InitializeLogger()
	infra.InitializeCommand()

	infra.ServerConfig(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal()
	}
}

func run() {
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	daemon.Init()

	server := infra.NewServer(8080, infra.DebugMode)
	api.Init(server)

	daemon.Start()

	<-exitChan
	logrus.Infoln("Closing the daemon...")
	logrus.Exit(1)
}
