package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garciademarina/verse/pkg/server"
	"github.com/go-chi/jwtauth"
)

var port = flag.Int("port", 8080, "-port=<port> sets the server's listening port. 8080 by default.")
var env = flag.String("env", "prod", "-env=<environment> specifies the environment. prod by default.")

func init() {
	flag.Parse()

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:{id-here}}` here:
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n", tokenString)
	_, tokenString, _ = tokenAuth.Encode(jwt.MapClaims{"user_id": "01D3XZ3ZHCP3KG9VT4FGAD8KDR"})
	fmt.Printf("DEBUG: a User 1 jwt is %s\n", tokenString)
	_, tokenString, _ = tokenAuth.Encode(jwt.MapClaims{"user_id": "01D3XZ7CN92AKS9HAPSZ4D5DP9"})
	fmt.Printf("DEBUG: a User 2 jwt is %s\n", tokenString)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	ensureInterruptionsStopApplication(cancel, logger)

	config := server.NewConfig(*port, *env)
	s := server.NewServer(config, logger)

	err := s.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func ensureInterruptionsStopApplication(cancelFunc context.CancelFunc, logger *log.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-c
		logger.Println(fmt.Sprintf("Got signal %s. Stopping server...", s))
		cancelFunc()

		os.Exit(1)
	}()
}
