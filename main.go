package main

import (
	"blogger/config"
	"blogger/src/router"
	"fmt"
	"log"
	"net/http"
)

// func init() {
// 	key := make([]byte, 64)
// 	if _, err := rand.Read(key); err != nil {
// 		log.Fatal(err)
// 	}

// 	stringBase64 := base64.StdEncoding.EncodeToString(key)
// 	fmt.Println("Suggested SECRET_KEY (put in .env file):", stringBase64)
// }

func main() {
	config.Initialize()
	r := router.Router()

	fmt.Printf("API on http://localhost:%d/api/v1/\n", config.Port)

	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", config.Port),
		r,
	); err != nil {
		log.Fatalln(err.Error())
	}
}
