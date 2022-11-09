package api

import "fmt"

func routeNodeEndpoints(prefix string, s *Server) {
	s.server.Post(fmt.Sprintf("%s", prefix), handleFunc(s, postNewNodeHandler))
	s.server.Get(fmt.Sprintf("%s/status", prefix), handleFunc(s, getStatusHandler))
}
