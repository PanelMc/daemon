package daemon

import (
	"time"

	"github.com/panelmc/daemon/api/socket"
	"github.com/panelmc/daemon/types"
)

func (s *Server) onPlayerJoin(player *types.Player) {
	s.Stats.OnlinePlayers = append(s.Stats.OnlinePlayers, player)

	socket.Broadcast("player_join", socket.ServerPlayerJoinPayload{
		ServerID: s.ID,
		Player:   player,
	})
}

func (s *Server) onPlayerLeave(name string) {
	player := &types.Player{
		Name: name,
	}
	for i, v := range s.Stats.OnlinePlayers {
		if v.Name == name {
			s.Stats.OnlinePlayers = append(s.Stats.OnlinePlayers[:i], s.Stats.OnlinePlayers[i+1:]...)
			player = v
		}
	}

	socket.Broadcast("player_leave", socket.ServerPlayerLeavePayload{
		ServerID: s.ID,
		Player:   player,
	})
}

func (s *Server) onStatusUpdate(status types.ServerStatus) {
	socket.Broadcast("server_status_update", socket.ServerStatusUpdatePayload{
		ServerID: s.ID,
		Status:   status.String(),
	})
}

func (s *Server) onDie() {
	s.UpdateStatus(types.ServerStatusOffline)

	if s.willRestart {
		payload := socket.ServerConsolePayload{
			ServerID: s.ID,
			Line:     "Server stopped! Restarting in 5 seconds...",
		}
		s.Logger().WithField("event", "Console").Info(payload.Line)
		socket.BroadcastTo(s.ID, "console_output", payload)

		go func() {
			time.Sleep(5 * time.Second)
			if s.Stats.Status == types.ServerStatusOffline {
				if err := s.Start(); err != nil {
					s.Logger().WithError(err).Error("There was an error while trying to start the server.")
				}
			}
		}()
	}
}
