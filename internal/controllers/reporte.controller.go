package controllers

import (
	"backend/internal/db"
	"backend/internal/models"
	"fmt"
	"strings"
	"time"

	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"gorm.io/datatypes"
)

type reporte struct{}

var Reporte = reporte{}

var background = &props.Color{
	Red:   200,
	Green: 200,
	Blue:  200,
}

func (reporte) CriticaPDF(w http.ResponseWriter, r *http.Request) {
	m, err := MakePDFCritica()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	doc, err := m.Generate()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_critica.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		log.Printf("Error escribiendo PDF en la respuesta: %v", err)
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
	}
}

func MakePDFCritica() (core.Maroto, error) {
	mrt := maroto.New()
	m := maroto.NewMetricsDecorator(mrt)

	err := m.RegisterHeader(getPageHeader())

	if err != nil {
		log.Fatal(err.Error())
	}

	m.RegisterFooter(getPageFooter())

	if err != nil {
		log.Fatal(err.Error())
	}

	m.AddRows(text.NewRow(22, "Reporte Criticas", props.Text{
		Top:   3,
		Size:  20,
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
	return row.New(8).Add(
		// text.NewCol(1, "N°", props.Text{Style: fontstyle.Bold}),
		text.NewCol(1, "COD", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(6, "DESCRIPCION", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "TIPO", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(1, "ESTADO", props.Text{Style: fontstyle.Bold, Align: align.Center}),
	)
}

func (o critica) GetContent(i int) core.Row {
	r := row.New(8).Add(
		// text.NewCol(1, strconv.Itoa(i+1)), // Agregar el índice (N°)
		text.NewCol(1, strconv.FormatUint(uint64(o.COD), 10), props.Text{Align: align.Center}),
		text.NewCol(6, o.Descripcion),
		text.NewCol(3, o.Tipo, props.Text{Align: align.Center}),
		text.NewCol(1, o.Estado, props.Text{Align: align.Center}),
	)

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: background,
		})
	}

	return r
}

func (reporte) LecturadoresPDF(w http.ResponseWriter, r *http.Request) {
	m, err := MakePDFLecturadores()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	doc, err := m.Generate()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_lecturadores.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		log.Printf("Error escribiendo PDF en la respuesta: %v", err)
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
	}
}

func MakePDFLecturadores() (core.Maroto, error) {
	mrt := maroto.New()
	m := maroto.NewMetricsDecorator(mrt)

	err := m.RegisterHeader(getPageHeader())

	if err != nil {
		log.Fatal(err.Error())
	}

	m.RegisterFooter(getPageFooter())

	if err != nil {
		log.Fatal(err.Error())
	}

	m.AddRows(text.NewRow(22, "Reporte Lecturadores", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Size:  20,
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

type lecturador struct {
	COD      uint   `gorm:"primaryKey;AutoIncrement" json:"cod"`
	Usuario  string `json:"usuario"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	CI       string `json:"ci"`
}

func (o lecturador) GetHeader() core.Row {
	return row.New(8).Add(
		text.NewCol(1, "N°", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		// text.NewCol(2, "COD", props.Text{Style: fontstyle.Bold}),
		text.NewCol(3, "USUARIO", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "NOMBRE", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "APELLIDO", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(2, "CEDULA DE IDENTIDAD", props.Text{Style: fontstyle.Bold, Align: align.Center}),
	)
}

func (o lecturador) GetContent(i int) core.Row {
	r := row.New(8).Add(
		text.NewCol(1, strconv.Itoa(i+1), props.Text{Align: align.Center}),
		// text.NewCol(2, strconv.FormatUint(uint64(o.COD), 10)),
		text.NewCol(3, o.Usuario, props.Text{Align: align.Center}),
		text.NewCol(3, o.Nombre, props.Text{Align: align.Center}),
		text.NewCol(3, o.Apellido, props.Text{Align: align.Center}),
		text.NewCol(2, o.CI, props.Text{Align: align.Center}),
	)

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: background,
		})
	}

	return r
}

func (reporte) MedidoresPDF(w http.ResponseWriter, r *http.Request) {
	m, err := MakePDFMedidores()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	doc, err := m.Generate()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_medidores.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		log.Printf("Error escribiendo PDF en la respuesta: %v", err)
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
	}
}

func MakePDFMedidores() (core.Maroto, error) {
	mrt := maroto.New()
	m := maroto.NewMetricsDecorator(mrt)

	err := m.RegisterHeader(getPageHeader())

	if err != nil {
		log.Fatal(err.Error())
	}

	m.RegisterFooter(getPageFooter())

	if err != nil {
		log.Fatal(err.Error())
	}

	m.AddRows(text.NewRow(22, "Reporte Medidores", props.Text{
		Top:   3,
		Size:  20,
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

type medidorRuta struct {
	MedidorNombre string `json:"medidor_nombre"`
	Propietario   string `json:"propietario"`
	RutaNombre    string `json:"ruta_nombre"`
	Zona          string `json:"zona"`
}

func (o medidorRuta) GetHeader() core.Row {
	return row.New(8).Add(
		text.NewCol(1, "N°", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "Medidor", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "Propietario", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(2, "Ruta", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(2, "Zona", props.Text{Style: fontstyle.Bold, Align: align.Center}),
	)
}

func (o medidorRuta) GetContent(i int) core.Row {
	r := row.New(8).Add(
		text.NewCol(1, strconv.Itoa(i+1), props.Text{Align: align.Center}),
		text.NewCol(3, o.MedidorNombre, props.Text{Align: align.Center}),
		text.NewCol(3, o.Propietario, props.Text{Align: align.Center}),
		text.NewCol(2, o.RutaNombre, props.Text{Align: align.Center}),
		text.NewCol(2, o.Zona, props.Text{Align: align.Center}),
	)

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: background,
		})
	}

	return r
}

func (reporte) LecturacionPDF(w http.ResponseWriter, r *http.Request) {
	codUsuario := mux.Vars(r)["codUsuario"]

	m, err := MakePDFLecturacion(codUsuario)
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	doc, err := m.Generate()
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_lecturacion.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		log.Printf("Error escribiendo PDF en la respuesta: %v", err)
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
	}
}

func MakePDFLecturacion(codUsuario string) (core.Maroto, error) {
	mrt := maroto.New()
	m := maroto.NewMetricsDecorator(mrt)

	err := m.RegisterHeader(getPageHeader())
	if err != nil {
		log.Fatal(err.Error())
	}

	m.RegisterFooter(getPageFooter())

	var nombreUsuario string
	err = db.GDB.Raw("SELECT usuario FROM usuario WHERE cod = ? LIMIT 1", codUsuario).Scan(&nombreUsuario).Error
	if err != nil {
		return nil, err
	}

	titulo := "Reporte de lecturaciones del usuario: " + nombreUsuario
	m.AddRows(text.NewRow(22, titulo, props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Size:  18,
		Align: align.Center,
	}))

	// Obtener medidores a través de los grupos y rutas
	var medidores []models.Medidor
	queryMedidores := `
		SELECT m.*
		FROM medidor m
		INNER JOIN ruta r ON m.cod_ruta = r.cod
		INNER JOIN grupo g ON g.cod_ruta = r.cod
		WHERE g.cod_usuario = ?
	`
	if err := db.GDB.Raw(queryMedidores, codUsuario).Scan(&medidores).Error; err != nil {
		return nil, err
	}

	if len(medidores) == 0 {
		return nil, fmt.Errorf("no se encontraron medidores para el usuario con código %s", codUsuario)
	}

	// Obtener lecturaciones realizadas
	var lecturacionesRealizadas []Lecturacion
	queryLecturaciones := `
		SELECT l.cod_medidor, m.nombre AS nombre_medidor, l.medicion, l.hora, l.fecha
		FROM lecturacion l
		INNER JOIN medidor m ON l.cod_medidor = m.cod
		INNER JOIN grupo g ON g.cod_usuario = ?
		WHERE l.cod_medidor IN (
			SELECT cod 
			FROM medidor 
			WHERE cod_ruta IN (
				SELECT cod_ruta 
				FROM grupo 
				WHERE cod_usuario = ?
			)
		)
		ORDER BY l.fecha, l.hora
	`
	if err := db.GDB.Raw(queryLecturaciones, codUsuario, codUsuario).Scan(&lecturacionesRealizadas).Error; err != nil {
		return nil, err
	}

	// Filtrar medidores no lecturados
	medidoresLecturados := make(map[uint]bool)
	for _, l := range lecturacionesRealizadas {
		medidoresLecturados[l.CodMedidor] = true
	}

	var medidoresNoLecturados []models.Medidor
	for _, medidor := range medidores {
		if !medidoresLecturados[medidor.COD] {
			medidoresNoLecturados = append(medidoresNoLecturados, medidor)
		}
	}

	// Construir tabla de lecturaciones realizadas
	m.AddRows(text.NewRow(10, "Lecturaciones Realizadas", props.Text{
		Style: fontstyle.Bold,
		Size:  14,
		Align: align.Left,
	}))

	if len(lecturacionesRealizadas) > 0 {
		rows, err := list.Build[Lecturacion](lecturacionesRealizadas)
		if err != nil {
			log.Fatal(err.Error())
		}
		m.AddRows(rows...)
	} else {
		m.AddRows(text.NewRow(10, "No se encontraron lecturaciones realizadas.", props.Text{
			Align: align.Left,
			Size:  10,
		}))
	}

	// Construir tabla de medidores no lecturados
	if len(medidoresNoLecturados) > 0 {

		m.AddRows(text.NewRow(10, "Medidores No Lecturados", props.Text{
			Top:   2,
			Style: fontstyle.Bold,
			Size:  14,
			Align: align.Left,
		}))

		m.AddRows(row.New(10).Add(
			text.NewCol(2, "Código", props.Text{Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(4, "Nombre", props.Text{Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(4, "Propietario", props.Text{Style: fontstyle.Bold, Align: align.Center}),
			text.NewCol(2, "Estado", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		))
		for _, medidor := range medidoresNoLecturados {
			m.AddRows(row.New(10).Add(
				text.NewCol(2, strconv.Itoa(int(medidor.COD)), props.Text{Align: align.Center}),
				text.NewCol(4, medidor.Nombre, props.Text{Align: align.Center}),
				text.NewCol(4, medidor.Propietario, props.Text{Align: align.Center}),
				text.NewCol(2, medidor.Estado, props.Text{Align: align.Center}),
			))
		}
	} else {
		m.AddRows(text.NewRow(10, "Todos los medidores tienen lecturaciones realizadas.", props.Text{
			Align: align.Left,
			Size:  10,
		}))
	}

	return m, nil
}

func MakePDFLecturacion2(codUsuario string) (core.Maroto, error) {
	mrt := maroto.New()
	m := maroto.NewMetricsDecorator(mrt)

	err := m.RegisterHeader(getPageHeader())
	if err != nil {
		log.Fatal(err.Error())
	}

	m.RegisterFooter(getPageFooter())
	if err != nil {
		log.Fatal(err.Error())
	}

	var nombreUsuario string
	err = db.GDB.Raw("SELECT usuario FROM usuario WHERE cod = ? limit 1", codUsuario).Scan(&nombreUsuario).Error
	if err != nil {
		return nil, err
	}

	titulo := "Reporte de lecturaciones del usuario: " + nombreUsuario

	m.AddRows(text.NewRow(22, titulo, props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Size:  18,
		Align: align.Center,
	}))

	var lista []Lecturacion
	query := `
		SELECT l.cod_medidor, m.nombre AS nombre_medidor, l.medicion, l.hora, l.fecha
		FROM lecturacion l
		INNER JOIN medidor m ON l.cod_medidor = m.cod
		WHERE l.cod_lecturador = ?
		ORDER BY l.fecha, l.hora;
	`

	if err := db.GDB.Raw(query, codUsuario).Scan(&lista).Error; err != nil {
		return nil, err
	}

	if len(lista) == 0 {
		return nil, fmt.Errorf("no se encontraron lecturaciones para el usuario con código %s", codUsuario)
	}

	rows, err := list.Build[Lecturacion](lista)
	if err != nil {
		log.Fatal(err.Error())
	}

	m.AddRows(rows...)
	return m, nil
}

type Lecturacion struct {
	CodMedidor    uint           `json:"codMedidor"`
	NombreMedidor string         `json:"nombreMedidor"`
	Medicion      *uint          `json:"medicion"`
	Hora          datatypes.Time `json:"hora"`
	Fecha         string         `json:"fecha"`
}

func (l Lecturacion) GetHeader() core.Row {
	return row.New(10).Add(
		text.NewCol(2, "Código Medidor", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(4, "Nombre Medidor", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(2, "Medición", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(2, "Hora", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(2, "Fecha", props.Text{Style: fontstyle.Bold, Align: align.Center}),
	)
}

func (l Lecturacion) GetContent(i int) core.Row {
	r := row.New(10).Add(
		text.NewCol(2, strconv.Itoa(int(l.CodMedidor)), props.Text{Align: align.Center}),
		text.NewCol(4, l.NombreMedidor, props.Text{Align: align.Center}),
		text.NewCol(2, strconv.Itoa(int(*l.Medicion)), props.Text{Align: align.Center}),
		text.NewCol(2, l.Hora.String(), props.Text{Align: align.Center}),
		text.NewCol(2, strings.Split(l.Fecha, "T")[0], props.Text{Align: align.Center}),
	)

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: background,
		})
	}
	return r
}

func getPageHeader() core.Row {
	return row.New(20).Add(
		image.NewFromFileCol(2, "internal/images/cosaalt_logo.png", props.Rect{
			Center:  true,
			Percent: 80,
		}),
		col.New(6),
		col.New(4).Add(
			text.New("Cooperativa de Servicios Públicos de Agua Potable y Alcantarillado Sanitario Tarija", props.Text{
				Top:   3,
				Size:  8,
				Style: fontstyle.Italic,
				Align: align.Right,
				Color: &props.BlackColor,
			}),
			text.New("TARIJA - BOLIVIA", props.Text{
				Top:   12,
				Size:  8,
				Align: align.Right,
				Color: &props.BlackColor,
			}),
		),
	)
}

func getPageFooter() core.Row {
	now := time.Now()

	formated := now.Format("02/01/2006 15:04:05")

	return row.New(20).Add(
		col.New(12).Add(
			text.New("Generada por Encargado Ricardo Campos", props.Text{
				Top:   12,
				Size:  8,
				Align: align.Left,
				Color: &props.BlackColor,
			}),
			text.New("Fecha y hora: "+formated, props.Text{
				Top:   16,
				Size:  8,
				Align: align.Left,
				Color: &props.BlackColor,
			}),
		),
	)
}
