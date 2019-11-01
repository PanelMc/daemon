package daemon

type ServerMap map[string]*Server

var servers = make(ServerMap)

func GetServers() *ServerMap {
	return &servers
}

// GetServerByID - Get a server by it's id. If no server found, return nil
func GetServerByID(id string) *Server {
	for _, s := range servers {
		if s.ID == id {
			return s
		}
	}

	return nil
}

// GetServer - Get a server by it's id or name. If no server found, return nil
func GetServer(server string) *Server {
	for _, s := range servers {
		if s.ID == server || s.Name == server {
			return s
		}
	}

	return nil
}
