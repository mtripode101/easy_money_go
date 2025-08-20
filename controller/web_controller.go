package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mtripode101/easymoney-go/model"
	"github.com/mtripode101/easymoney-go/service"
)

// Registrar rutas web
func RegisterWebRoutes(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*.html")

	group := r.Group("/web/easy-money")
	{
		group.GET("/", list)
		group.GET("/new", createForm)
		group.POST("/save", save)
		group.GET("/edit/:id", editForm)
		group.GET("/delete/:id", deleteWeb)
		group.GET("/differences", differences)
		group.GET("/total", total)
	}
}

// Mostrar lista
func list(c *gin.Context) {
	entries, _ := service.FindAll()
	c.HTML(http.StatusOK, "index.html", gin.H{"entries": entries})
}

// Formulario de creación
func createForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", gin.H{
		"easyMoney": model.EasyMoney{
			Date: time.Now(), // ← valor válido para el input
		},
	})
}

func save(c *gin.Context) {
	var em model.EasyMoney
	if err := c.ShouldBind(&em); err == nil {
		fmt.Printf("Recibido: %+v\n", em)
		err := service.Save(&em)
		if err != nil {
			fmt.Println("Error al guardar:", err)
		} else {
			fmt.Println("Guardado OK")
		}
	} else {
		fmt.Println("Error de bind:", err)
	}
	c.Redirect(http.StatusSeeOther, "/web/easy-money/")
}

// Formulario de edición
func editForm(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	em, _ := service.FindByID(uint(id))
	c.HTML(http.StatusOK, "form.html", gin.H{"easyMoney": em})
}

// Eliminar entrada (renombrada para evitar conflicto)
func deleteWeb(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	service.Delete(uint(id))
	c.Redirect(http.StatusSeeOther, "/web/easy-money/")
}

// Mostrar diferencias consecutivas
func differences(c *gin.Context) {
	startStr := c.Query("start")
	endStr := c.Query("end")

	start, err1 := time.Parse("2006-01-02", startStr)
	end, err2 := time.Parse("2006-01-02", endStr)

	if err1 != nil || err2 != nil {
		c.HTML(http.StatusOK, "differences.html", gin.H{"error": "Fechas inválidas. Usa el formato YYYY-MM-DD."})
		return
	}

	diffs, _ := service.CalculateConsecutiveDifferences(start, end)
	c.HTML(http.StatusOK, "differences.html", gin.H{"differences": diffs})
}

// Mostrar diferencia total
func total(c *gin.Context) {
	start, _ := time.Parse("2006-01-02", c.Query("start"))
	end, _ := time.Parse("2006-01-02", c.Query("end"))

	result, _ := service.CalculateTotalMoneyDifference(start, end)
	c.HTML(http.StatusOK, "total.html", gin.H{"total": result})
}
