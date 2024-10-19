package test

import (
	"log"
	"testing"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/core"
)

// CriticaPDF es una función que realiza una operación en el reporte.
// Puedes implementar la lógica de esta función aquí.
func TestCriticaPDF(t *testing.T) {
	m := GetMaroto()
	document, err := m.Generate()
	if err != nil {
		log.Fatal(err)
	}

	err = document.Save("simplestv2.pdf")
	if err != nil {
		log.Fatal(err)
	}
	// Implementa la lógica aquí.
}
func GetMaroto() core.Maroto {
	m := maroto.New()

	m.AddRow(20,
		//code.NewBarCol(4, "barcode"),
		code.NewMatrixCol(4, "matrixcode"),
		//code.NewQrCol(4, "qrcode"),
	)

	m.AddRow(10, col.New(12))
	return m
}
