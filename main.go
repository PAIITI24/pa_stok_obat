package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hakushigo/stok_obat/controller"
	"github.com/hakushigo/stok_obat/helper"
)

func main() {
	helper.Migrator()

	srv := fiber.New(
		fiber.Config{
			Immutable: false,
			AppName:   "Stok_Management_Obat",
		})

	srv.Put("/obat/stok/add", controller.AddStok)
	srv.Put("/obat/stok/reduce", controller.ReduceStok)

	err := srv.Listen(":3010")

	if err != nil {
		panic(err)
	}
}
