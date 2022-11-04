package api

import "fmt"

func routeNodeEndpoints(prefix string, s *Server) {
	s.server.Get(fmt.Sprintf("%s/status", prefix), handleFunc(s, getStatusHandler))
}
