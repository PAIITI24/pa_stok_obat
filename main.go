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

	stokGroup := srv.Group("/obat/stok")

	stokGroup.Put("/add", controller.AddStok)
	stokGroup.Put("/reduce", controller.ReduceStok)
	stokGroup.Get("/add/history", controller.ListStokMasuk)
	stokGroup.Get("/reduce/history", controller.ListStokKeluar)
	stokGroup.Get("/add/history/:id", controller.ListStokMasukOfObat)
	stokGroup.Get("/reduce/history/:id", controller.ListStokKeluarOfObat)

	err := srv.Listen(":3010")

	if err != nil {
		panic(err)
	}
}
