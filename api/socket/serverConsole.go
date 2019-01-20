package socket

import (
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type ServerRoom struct {
	ServerId string `json:"server_id"`
}

type serverConsolePayload struct {
	ServerId string `json:"server_id"`

	Line string `json:"line"`
}

var _ io.Writer = ServerRoom{}

func (s ServerRoom) Write(b []byte) (n int, e error) {
	//l := make([]byte, len(b))
	//copy(l, b)
	line := string(b)

	if len(b) == 1 {
		// Is a key pressed, ignore
		return 1, nil
	} else if line == "\r\n" || strings.Contains(line, "\b") {
		return 2, nil
	}

	if strings.Contains(line, "\r") {
		line = strings.Replace(line, "\r", "", -1)
	}
	if strings.Contains(line, "\n") {
		line = strings.Replace(line, "\n", "", -1)
	}

	payload := serverConsolePayload{s.ServerId, line}
	logrus.WithField("server", s.ServerId).WithField("event", "Console").
		Infof("%#v", line)
	//Infof("%s", payload)
	Server.BroadcastTo(s.ServerId, "console_output", payload)

	return len(b), nil
}
