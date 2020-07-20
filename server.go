package main

import (
	"encoding/json"
	"fmt"
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

	words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore := [][]string{}, [][]string{}, [][]string{}, []string{}, []string{}, [][]bool{}, []int{}
	input := ""
	sendResponses(app, input, words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore)
	app.Listen(getPort())
}

func sendResponses(app *fiber.App, input string, words [][]string, closestScansion [][]string, islah [][]string, closestMeters []string, closestMeterNames []string, problematicWords [][]bool, ravaniScore []int) {

	app.Post("/words", func(c *fiber.Ctx) {

		if len(words) == 0 || c.Body() != input {
			input = c.Body()

			words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore = test(input)
		}

		wordsJSON := encodeJSON(words)
		c.Send(wordsJSON)
	})

	app.Post("/closestScansion", func(c *fiber.Ctx) {
		if len(words) == 0 || c.Body() != input {
			input = c.Body()
			words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore = test(input)
		}

		for i := range closestScansion {
			closestScansion[i] = getUrduNumerals(closestScansion[i])
		}
		closestScansionJSON := encodeJSON(closestScansion)
		c.Send(closestScansionJSON)
	})

	app.Post("/islah", func(c *fiber.Ctx) {

		if len(words) == 0 || c.Body() != input {
			input = c.Body()
			words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore = test(input)

		}

		for i := range islah {
			islah[i] = getUrduNumerals(islah[i])
		}
		closestMeterKeysJSON := encodeJSON(islah)
		c.Send(closestMeterKeysJSON)
	})

	app.Post("/closestMeters", func(c *fiber.Ctx) {

		if len(words) == 0 || c.Body() != input {
			input = c.Body()
			words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore = test(input)

		}

		closestMetersJSON := encodeJSON(closestMeters)
		c.Send(closestMetersJSON)
	})

	app.Post("/closestMeterNames", func(c *fiber.Ctx) {

		if len(words) == 0 || c.Body() != input {
			input = c.Body()
			words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore = test(input)

		}

		closestMeterNamesJSON := encodeJSON(closestMeterNames)
		c.Send(closestMeterNamesJSON)

	})

	app.Post("/problematicWords", func(c *fiber.Ctx) {
		if len(words) == 0 || c.Body() != input {
			input = c.Body()
			words, closestScansion, islah, closestMeters, closestMeterNames, problematicWords, ravaniScore = test(input)

		}
		problematicWordsJSON := encodeJSON(problematicWords)
		c.Send(problematicWordsJSON)
	})

	app.Get("/ravani", func(c *fiber.Ctx) {
		ravaniJSON := encodeJSON(ravaniScore)
		c.Send(ravaniJSON)
	})

}

func encodeJSON(response interface{}) interface{} {

	responseJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error:", err)
	}
	return responseJSON
}
