package main

import (
	"MusicApp/internal/auth"
	"MusicApp/internal/server"
	"fmt"
)

// todo - fix bugs + refactor
func main() {
	auth.NewAuth()
	srv := server.NewServer()
	err := srv.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
