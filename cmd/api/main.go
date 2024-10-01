package main

import (
	"MusicApp/internal/server"
	"fmt"
)

// todo - fix bugs + refactor
func main() {
	srv := server.NewServer()
	err := srv.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
