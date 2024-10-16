package main

import (
	"MusicApp/internal/server"
	"fmt"
)

func main() {
	srv := server.NewServer()
	if err := srv.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))

	}

}
