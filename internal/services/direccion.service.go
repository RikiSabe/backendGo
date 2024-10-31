package services

import (
	"backend/internal/db"
	"backend/internal/models"
)

type direccion struct {
}

var Direccion direccion

func (d *direccion) GetByCod(i *models.Direccion, id string) error {
	tx := db.GDB.Begin()
	// Buscar el medidor por ID y asegurarse de que esté activo
	if err := tx.Where("cod = ?", id).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
