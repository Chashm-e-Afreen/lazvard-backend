package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
)

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return p
	}
	return "3000"
}
func main() {
	app := fiber.New()
	app.Settings.Prefork = true
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           0,
	}))

	words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore, dict := [][]string{}, [][]string{}, [][]string{}, []string{}, []string{}, [][]bool{}, []int{}, map[string][]string{}
	input := ""
	sendResponses(app, input, words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore, dict)
	cer, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	app.Listen(getPort(), config)
}

func sendResponses(app *fiber.App, input string, words [][]string, closestScansion [][]string, islah [][]string, closestMeters []string, closestMeterNames []string, problematicWords [][]bool, ravaniScore []int, dict map[string][]string) {

	app.Post("/words", func(c *fiber.Ctx) {

		if len(words) == 0 || c.Body() != input {
			input = c.Body()

			words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore, dict = test(input)
		}
		for i := range islah {
			islah[i] = getUrduNumerals(islah[i])
		}
		for i := range closestScansion {
			closestScansion[i] = getUrduNumerals(closestScansion[i])
		}

		wordsJSON := encodeJSON(words)
		closestScansionJSON := encodeJSON(closestScansion)
		closestMetersJSON := encodeJSON(closestMeters)
		closestMeterKeysJSON := encodeJSON(islah)
		closestMeterNamesJSON := encodeJSON(closestMeterNames)
		problematicWordsJSON := encodeJSON(problematicWords)
		ravaniJSON := encodeJSON(ravaniScore)
		newline := "\n"
		c.Send(wordsJSON + newline + closestScansionJSON + newline + closestMeterKeysJSON + newline + closestMetersJSON + newline + closestMeterNamesJSON + newline + problematicWordsJSON + newline + ravaniJSON)
	})

}

func encodeJSON(response interface{}) string {

	responseJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(responseJSON)
}
