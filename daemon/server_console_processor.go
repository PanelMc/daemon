package daemon

import (
	"regexp"

	"github.com/panelmc/daemon/types"
)

type serverRegexMatch struct {
	Port    *regexp.Regexp
	Login   *regexp.Regexp
	LogOut  *regexp.Regexp
	Start   *regexp.Regexp
	Stop    *regexp.Regexp
	Message *regexp.Regexp
}

var vanillaRegex = &serverRegexMatch{
	Port:    regexp.MustCompile(`^\[[\d:]{8}] \[Server thread/INFO]: Starting Minecraft server on (.*?):([0-9]{5}|[0-9]{4})`),
	Login:   regexp.MustCompile(`^\[[\d:]{8}] \[Server thread/INFO]: (\w+)\[/([\d.:]+)] logged in`),
	LogOut:  regexp.MustCompile(`^\[[\d:]{8}] \[Server thread/INFO]: (\w+) lost connection`),
	Start:   regexp.MustCompile(`^\[[\d:]{8}] \[Server thread/INFO]: Done \((.*?)s\)! For help, type "help`),
	Stop:    regexp.MustCompile(`^\[[\d:]{8}] \[Server thread/INFO]: Stopping server`),
	Message: regexp.MustCompile(`^\[[\d:]{8}] \[Server thread/(\w+)]: (.*?)`),
}

var bukkitRegex = &serverRegexMatch{
	Port:    regexp.MustCompile(`^\[[\d:]{8} INFO]: Starting Minecraft server on (.*?):([0-9]{5}|[0-9]{4})`),
	Login:   regexp.MustCompile(`^\[[\d:]{8} INFO]: (\w+)\[/([\d.:]+)] logged in (.*?)`),
	LogOut:  regexp.MustCompile(`^\[[\d:]{8} INFO]: (\w+) lost connection`),
	Start:   regexp.MustCompile(`^\[[\d:]{8} INFO]: Done \((.*?)s\)! For help, type "help"`),
	Stop:    regexp.MustCompile(`^\[[\d:]{8} INFO]: Stopping server`),
	Message: regexp.MustCompile(`^\[[\d:]{8} (\w+)]: (.*?)`),
}

func processConsoleOutput(s *Server, output string) {
	if s.Type == "" {
		if vanillaRegex.Message.MatchString(output) {
			s.Type = "VANILLA"
			s.Save()
		} else if bukkitRegex.Message.MatchString(output) {
			s.Type = "BUKKIT"
			s.Save()
		}
	}

	var r *serverRegexMatch
	if s.Type == "VANILLA" {
		r = vanillaRegex
	} else if s.Type == "BUKKIT" {
		r = bukkitRegex
	} else {
		return
	}

	match := r.Login.FindAllStringSubmatch(output, -1)
	if len(match) > 0 {
		player := &types.Player{
			Name: match[0][1],
			Ip:   match[0][2],
		}

		s.onPlayerJoin(player)
	}

	match = r.LogOut.FindAllStringSubmatch(output, -1)
	if len(match) > 0 {
		s.onPlayerLeave(match[0][1])
	}

	match = r.Start.FindAllStringSubmatch(output, -1)
	if len(match) > 0 {
		s.UpdateStatus(types.ServerStatusOnline)
	}

	match = r.Stop.FindAllStringSubmatch(output, -1)
	if len(match) > 0 {
		s.UpdateStatus(types.ServerStatusStopping)
	}
}
