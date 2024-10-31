package services

import (
	"backend/internal/db"
	"backend/internal/models"
)

type medidor struct {
}

var Medidor medidor

func (m *medidor) CountActivos(total *int64) error {
	// Inicia una transacción solo para la consulta
	tx := db.GDB.Begin()

	// Cuenta los medidores activos
	if err := tx.Model(&models.Medidor{}).Where("estado = ?", "activo").Count(total).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *medidor) GetAll(l *[]models.Medidor) error {
	tx := db.GDB.Begin()
	// Buscar los medidores asegurarse de que esté activo
	if err := tx.Where("estado = ?", "activo").Preload("Ruta").Find(&l).Error; err != nil {
		tx.Rollback()
		return err
	}
	for i := range *l {
		(*l)[i].NombreRuta = (*l)[i].Ruta.Nombre
	}
	tx.Commit()
	return nil
}

func (m *medidor) GetByCod(i *models.Medidor, id string) error {
	tx := db.GDB.Begin()
	// Buscar el medidor por ID y asegurarse de que esté activo
	if err := tx.Where("cod = ?", id).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// func (m *medidor) GetByRuta(i *[]models.Medidor, id string) error {
// 	tx := db.GDB.Begin()

// 	if err := tx.Where("cod_ruta = ? and estado = ?", id, "activo").Find(&i).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}
// 	tx.Commit()
// 	return nil
// }

func (m *medidor) GetByRuta(i *[]models.Medidor, id string) error {
	tx := db.GDB.Begin()

	// Realiza el join entre Medidor y Direccion para obtener las coordenadas
	if err := tx.Preload("CodDireccion").Where("cod_ruta = ? and estado = ?", id, "activo").Find(&i).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (m *medidor) Save(i *models.Medidor) error {
	tx := db.GDB.Begin()

	// Intentar guardar el medidor en la base de datos
	if err := tx.Create(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *medidor) Delete(id string) error {
	var medidor models.Medidor
	tx := db.GDB.Begin()

	// Buscar el medidor por ID y asegurarse de que esté activo
	if err := tx.Where("cod = ? AND estado = ?", id, "activo").First(&medidor).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Marcar el medidor como inactivo en lugar de eliminarlo físicamente
	if err := tx.Model(&medidor).Update("estado", "inactivo").Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *medidor) Update(cod string, i *models.Medidor) error {
	//medidor := &i
	tx := db.GDB.Begin()
	if err := tx.Where("cod = ?", cod).First(&i).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil

}
