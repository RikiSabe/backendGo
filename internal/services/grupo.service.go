package services

import (
	"backend/internal/db"
	"backend/internal/models"
)

type GrupoService struct{}

func (g GrupoService) GetAll(grupos *[]models.Grupo) error {
	return db.GDB.Find(grupos).Error
}

func (g GrupoService) GetById(grupo *models.Grupo, cod string) error {
	return db.GDB.Where("cod = ?", cod).First(grupo).Error
}

func (g GrupoService) Save(grupo *models.Grupo) error {
	return db.GDB.Save(grupo).Error
}

func (g GrupoService) Delete(cod string) error {
	return db.GDB.Where("cod = ?", cod).Delete(&models.Grupo{}).Error
}

var Grupo GrupoService
