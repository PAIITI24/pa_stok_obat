package controller

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/hakushigo/stok_obat/model"
	"gorm.io/gorm"
)

func AddStok(ctx *fiber.Ctx) error {
	var request adreqbody
	var obat model.Obat
	var newInsertRecord model.StokMasuk

	// get the request body data
	err := json.Unmarshal(ctx.Body(), &request)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status_code": fiber.StatusInternalServerError,
			"error":       err,
		})
	}

	// then we decided to fetch the obat data.
	search := db.Where(&model.Obat{ID: request.ObatID}).First(&obat)
	if search.RowsAffected == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status_code": fiber.StatusNotFound,
			"error":       "The Obat can't be found",
		})
	}

	// because it's adding we don't need to amount checking beforehand
	// and then we added to record

	err = db.Transaction(func(tx *gorm.DB) error { // transaction as a condom of SQL
		obat.JumlahStok += request.Amount // first we update the amount of obat in the obat field
		if err := tx.Save(&obat).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status_code": fiber.StatusInternalServerError,
				"error":       err,
			})
		}

		// add item to the add stock table
		// get the obat that was meant

		newInsertRecord = model.StokMasuk{
			StokMasuk:   request.Amount,
			Obat:        obat,
			ExpiredDate: request.ExpiredDate,
		}
		if err := tx.Create(&newInsertRecord).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status_code": fiber.StatusInternalServerError,
				"error":       err,
			})
		}

		return nil
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status_code": fiber.StatusInternalServerError,
			"error":       err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"obat":        obat,
		"record_info": newInsertRecord,
	})
}

func ReduceStok(ctx *fiber.Ctx) error {
	var request redreqbody

	var obat model.Obat
	var stokMasuk model.StokMasuk
	var newOutRecord model.StokKeluar

	err := json.Unmarshal(ctx.Body(), &request)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status_code": fiber.StatusInternalServerError,
			"error":       err,
		})
	}

	// then we decided to fetch the obat data.
	search := db.Where(&model.Obat{ID: request.ObatID}).First(&obat)
	if search.RowsAffected == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status_code": fiber.StatusNotFound,
			"error":       "The Obat can't be found",
		})
	}

	// then we decided to fetch the obat data.
	search = db.Where(&model.StokMasuk{ID: request.StokMasukID}).First(&stokMasuk)
	if search.RowsAffected == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status_code": fiber.StatusNotFound,
			"error":       "The stok record can't be found",
		})
	}

	if obat.JumlahStok < request.Amount {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"error":       "The amount taken larger than what available",
		})
	}

	if stokMasuk.StokMasuk < request.Amount {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status_code": fiber.StatusBadRequest,
			"error":       "The amount taken larger than the stock",
		})
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		obat.JumlahStok -= request.Amount
		if err := tx.Save(&obat).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status_code": fiber.StatusInternalServerError,
				"error":       err,
			})
		}

		stokMasuk.StokMasuk -= request.Amount
		if err := tx.Save(&stokMasuk).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status_code": fiber.StatusInternalServerError,
				"error":       err,
			})
		}

		newOutRecord = model.StokKeluar{
			StokMasuk:  stokMasuk,
			StokKeluar: request.Amount,
			Obat:       obat,
		}

		if err := tx.Create(&newOutRecord).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status_code": fiber.StatusInternalServerError,
				"error":       err,
			})
		}

		return nil
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status_code": fiber.StatusInternalServerError,
			"error":       err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"obat":        obat,
		"record_info": newOutRecord,
	})
}

func ListStokMasuk(ctx *fiber.Ctx) error {
	var stokMasukList []model.StokMasuk

	// Fetch the list of stock additions, ordered by their creation date
	err := db.Preload("Obat").Order("created_at desc").Find(&stokMasukList).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status_code": fiber.StatusInternalServerError,
			"error":       err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(stokMasukList)
}

func ListStokMasukOfObat(ctx *fiber.Ctx) error {
	var stokMasukList []model.StokMasuk
	var id, _ = ctx.ParamsInt("id")

	// Fetch the list of stock additions, ordered by their creation date
	err := db.Preload("Obat").Order("created_at desc").Where(&model.StokMasuk{
		ObatID: id,
	}).Find(&stokMasukList).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status_code": fiber.StatusInternalServerError,
			"error":       err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(stokMasukList)
}

func ListStokKeluar(ctx *fiber.Ctx) error {
	var stokKeluarList []model.StokKeluar

	// Fetch the list of stock reductions, ordered by their creation date
	err := db.Preload("Obat").Order("created_at desc").Find(&stokKeluarList).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status_code": fiber.StatusInternalServerError,
			"error":       err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(stokKeluarList)
}
func ListStokKeluarOfObat(ctx *fiber.Ctx) error {
	var stokKeluarList []model.StokKeluar
	var id, _ = ctx.ParamsInt("id")

	// Fetch the list of stock reductions, ordered by their creation date
	err := db.Preload("Obat").Order("created_at desc").Where(&model.StokKeluar{
		ObatID: id,
	}).Find(&stokKeluarList).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status_code": fiber.StatusInternalServerError,
			"error":       err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(stokKeluarList)
}
