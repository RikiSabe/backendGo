package services

import (
	"backend/internal/db"
	"backend/internal/models"
)

// GrupoService define los métodos para trabajar con grupos
type GrupoService struct{}

// GetAll obtiene todos los grupos
func (g GrupoService) GetAll(grupos *[]models.Grupo) error {
	return db.GDB.Find(grupos).Error
}

// GetById obtiene un grupo por su código
func (g GrupoService) GetById(grupo *models.Grupo, cod string) error {
	return db.GDB.Where("cod = ?", cod).First(grupo).Error
}

// Save guarda un nuevo grupo o actualiza uno existente
func (g GrupoService) Save(grupo *models.Grupo) error {
	return db.GDB.Save(grupo).Error
}

// Delete elimina un grupo por su código
func (g GrupoService) Delete(cod string) error {
	return db.GDB.Where("cod = ?", cod).Delete(&models.Grupo{}).Error
}

// Instancia global del servicio de grupos
var Grupo GrupoService
