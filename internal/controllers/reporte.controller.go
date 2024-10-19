package controllers

import (
	"backend/internal/db"

	"log"
	"net/http"
	"strconv"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
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

// CriticaPDF

func (reporte) CriticaPDF(w http.ResponseWriter, r *http.Request) {
	m, err := MakePDFCritica()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}
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
func MakePDFCritica() (core.Maroto, error) {
	mrt := maroto.New()
	m := maroto.NewMetricsDecorator(mrt)

	// Añadir un título al documento
	m.AddRows(text.NewRow(20, "Reporte Criticas", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	var lista []critica
	query := "SELECT c.cod,c.tipo,c.estado, c.descripcion FROM critica c ORDER BY c.cod;"
	if err := db.GDB.Raw(query).Scan(&lista).Error; err != nil {
		return nil, err
	}
	rows, err := list.Build[critica](lista)
	if err != nil {
		log.Fatal(err.Error())
	}

	m.AddRows(rows...)
	return m, nil
}

type critica struct {
	COD         uint   `gorm:"primaryKey;AutoIncrement" json:"cod"`
	Tipo        string `json:"tipo"`
	Estado      string `json:"estado"`
	Descripcion string `json:"descripcion"`
}

func (o critica) GetHeader() core.Row {
	return row.New(10).Add(
		text.NewCol(2, "N°", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "COD", props.Text{Style: fontstyle.Bold}),
		text.NewCol(4, "DESCRIPCION", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "TIPO", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "ESTADO", props.Text{Style: fontstyle.Bold}),
	)
}

func (o critica) GetContent(i int) core.Row {
	r := row.New(6).Add(
		text.NewCol(2, strconv.Itoa(i+1)), // Agregar el índice (N°)
		text.NewCol(2, strconv.FormatUint(uint64(o.COD), 10)),
		text.NewCol(4, o.Descripcion),
		text.NewCol(2, o.Tipo),
		text.NewCol(2, o.Estado),
	)

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: background,
		})
	}

	return r
}

// Lecturadores

func (reporte) LecturadoresPDF(w http.ResponseWriter, r *http.Request) {
	m, err := MakePDFLecturadores()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}
	// Generar el PDF en el buffer
	doc, err := m.Generate()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	// Configurar los encabezados para la respuesta HTTP
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_lecturadores.pdf\"")

	// Escribir el contenido del buffer en la respuesta
	if _, err := w.Write(doc.GetBytes()); err != nil {
		log.Printf("Error escribiendo PDF en la respuesta: %v", err)
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
	}
}

func MakePDFLecturadores() (core.Maroto, error) {
	mrt := maroto.New()
	m := maroto.NewMetricsDecorator(mrt)

	// Añadir un título al documento
	m.AddRows(text.NewRow(20, "Reporte Lecturadores", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	var lista []lecturador
	query := `
		SELECT u.cod, u.usuario, p.nombre, p.apellido, p.ci
		FROM usuario u
		INNER JOIN persona p ON u.cod_persona = p.cod
		WHERE u.rol = 'lecturador'
		ORDER BY u.cod;`

	if err := db.GDB.Raw(query).Scan(&lista).Error; err != nil {
		return nil, err
	}
	rows, err := list.Build[lecturador](lista)
	if err != nil {
		log.Fatal(err.Error())
	}

	m.AddRows(rows...)
	return m, nil
}

// Estructura para lecturadores
type lecturador struct {
	COD      uint   `gorm:"primaryKey;AutoIncrement" json:"cod"`
	Usuario  string `json:"usuario"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	CI       string `json:"ci"`
}

func (o lecturador) GetHeader() core.Row {
	return row.New(10).Add(
		text.NewCol(1, "N°", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "COD", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "USUARIO", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "NOMBRE", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "APELLIDO", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "CEDULA DE IDENTIDAD", props.Text{Style: fontstyle.Bold}),
	)
}

func (o lecturador) GetContent(i int) core.Row {
	r := row.New(6).Add(
		text.NewCol(1, strconv.Itoa(i+1)), // Agregar el índice (N°)
		text.NewCol(2, strconv.FormatUint(uint64(o.COD), 10)),
		text.NewCol(2, o.Usuario),
		text.NewCol(2, o.Nombre),
		text.NewCol(2, o.Apellido),
		text.NewCol(2, o.CI),
	)

	// Alternar el color de fondo para filas pares
	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: background,
		})
	}

	return r
}

// Medidores PDF

func (reporte) MedidoresPDF(w http.ResponseWriter, r *http.Request) {
	m, err := MakePDFMedidores()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}
	// Generar el PDF en el buffer
	doc, err := m.Generate()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	// Configurar los encabezados para la respuesta HTTP
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_medidores.pdf\"")

	// Escribir el contenido del buffer en la respuesta
	if _, err := w.Write(doc.GetBytes()); err != nil {
		log.Printf("Error escribiendo PDF en la respuesta: %v", err)
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
	}
}

func MakePDFMedidores() (core.Maroto, error) {
	mrt := maroto.New()
	m := maroto.NewMetricsDecorator(mrt)

	// Añadir un título al documento
	m.AddRows(text.NewRow(20, "Reporte Medidores", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	var lista []medidorRuta
	query := `
		SELECT m.nombre AS medidor_nombre, m.propietario, r.nombre AS ruta_nombre, r.zona
		FROM medidor m
		INNER JOIN ruta r ON m.cod_ruta = r.cod
		ORDER BY m.cod;
	`
	if err := db.GDB.Raw(query).Scan(&lista).Error; err != nil {
		return nil, err
	}
	rows, err := list.Build[medidorRuta](lista)
	if err != nil {
		log.Fatal(err.Error())
	}

	m.AddRows(rows...)
	return m, nil
}

// Estructura para los datos combinados de Medidor y Ruta
type medidorRuta struct {
	MedidorNombre string `json:"medidor_nombre"`
	Propietario   string `json:"propietario"`
	RutaNombre    string `json:"ruta_nombre"`
	Zona          string `json:"zona"`
}

func (o medidorRuta) GetHeader() core.Row {
	return row.New(10).Add(
		text.NewCol(1, "N°", props.Text{Style: fontstyle.Bold}),
		text.NewCol(3, "Medidor", props.Text{Style: fontstyle.Bold}),
		text.NewCol(3, "Propietario", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Ruta", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Zona", props.Text{Style: fontstyle.Bold}),
	)
}

func (o medidorRuta) GetContent(i int) core.Row {
	r := row.New(6).Add(
		text.NewCol(1, strconv.Itoa(i+1)), // Agregar el índice (N°)
		text.NewCol(3, o.MedidorNombre),
		text.NewCol(3, o.Propietario),
		text.NewCol(2, o.RutaNombre),
		text.NewCol(2, o.Zona),
	)

	// Alternar el color de fondo para filas pares
	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: background,
		})
	}

	return r
}
