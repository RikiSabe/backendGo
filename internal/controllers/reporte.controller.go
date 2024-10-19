package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type reporte struct{}

var Reporte = reporte{}

var background = &props.Color{
	Red:   200,
	Green: 200,
	Blue:  200,
}

// CriticaPDF genera un reporte PDF y lo retorna en la respuesta HTTP.
func (reporte) CriticaPDF(w http.ResponseWriter, r *http.Request) {
	m := MakePDFCritica()

	// Generar el PDF en el buffer
	doc, err := m.Generate()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	// Configurar los encabezados para la respuesta HTTP
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_critica.pdf\"")

	// Escribir el contenido del buffer en la respuesta
	if _, err := w.Write(doc.GetBytes()); err != nil {
		log.Printf("Error escribiendo PDF en la respuesta: %v", err)
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
	}
}

// MakePDFCritica crea un PDF con una lista dinámica.
func MakePDFCritica() core.Maroto {
	mrt := maroto.New()
	m := maroto.NewMetricsDecorator(mrt)

	// Añadir un título al documento

	objects := getObjects(10)
	rows, err := list.Build[Object](objects)
	if err != nil {
		log.Fatal(err.Error())
	}

	m.AddRows(rows...)
	return m
}

type Object struct {
	Key   string
	Value string
}

func (o Object) GetHeader() core.Row {
	return row.New(10).Add(
		text.NewCol(4, "ID", props.Text{Style: fontstyle.Bold}),
		text.NewCol(8, "Valor", props.Text{Style: fontstyle.Bold}),
	)
}

func (o Object) GetContent(i int) core.Row {
	r := row.New(5).Add(
		text.NewCol(4, o.Key),
		text.NewCol(8, o.Value),
	)

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: background,
		})
	}

	return r
}

func getObjects(max int) []Object {
	var objects []Object
	for i := 0; i < max; i++ {
		objects = append(objects, Object{
			Key:   fmt.Sprintf("ID: %d", i),
			Value: fmt.Sprintf("Valor: %d", i*10),
		})
	}
	return objects
}
