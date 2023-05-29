package main

import (
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/siwacarn/cash-flow-line-bot/database"
	"github.com/siwacarn/cash-flow-line-bot/models"
)

func main() {
	// create connection of firebase firestore
	ctx := context.Background()
	db, err := database.NewFirestoreDatabase(ctx, "aomsin")
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	api := app.Group("/aomsin")

	api.Get("/all", func(c *fiber.Ctx) error {
		records, err := db.ReadAll()
		if err != nil {
			return err
		}
		return c.JSON(records)
	})

	// <localhost>/all/:id
	api.Get("/:id", func(c *fiber.Ctx) error {
		idstr := c.Params("id")

		idi, err := strconv.Atoi(idstr)
		if err != nil {
			log.Println(err)
			return fiber.ErrBadRequest
		}

		record, err := db.ReadOne(idi)
		if err != nil {
			log.Println(err)
			return fiber.ErrNotFound
		}

		return c.JSON(record)
	})

	api.Post("/create_tx", func(c *fiber.Ctx) error {
		model := new(models.Details)
		if err := c.BodyParser(model); err != nil {
			return err
		}

		err = db.Create(model.Txname, model.Amount)
		if err != nil {
			return err
		}

		return c.SendString("Create database successful!")
	})

	app.Listen(":3000")
}

// func CreateTransaction(db *database.FirestoreDatabase) func(c *fiber.Ctx) error {
// 	err := db.Create()
// 	if err != nil {
// 		log.Println(err)
// 		return func(c *fiber.Ctx) error {
// 			return errors.New("cannot create database on firestore")
// 		}
// 	}

// 	return func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.StatusOK)
// 	}
// }
