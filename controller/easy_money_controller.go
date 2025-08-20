package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mtripode101/easymoney-go/dto"
	"github.com/mtripode101/easymoney-go/model"
	"github.com/mtripode101/easymoney-go/service"
)

// Registrar rutas
func RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/easy-money")
	{
		group.POST("", create)
		group.PUT("/:id", update)
		group.DELETE("/:id", delete)
		group.GET("", getAll)
		group.GET("/difference", getDifference)
		group.GET("/by-date", getByDateRange)
		group.GET("/search", searchByDescription)
		group.GET("/differences", getConsecutiveDifferences)
		group.GET("/total-money-difference", getTotalMoneyDifference)
	}
}

// Crear nuevo registro
func create(c *gin.Context) {
	var em model.EasyMoney
	if err := c.ShouldBindJSON(&em); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.Save(&em); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, em)
}

// Actualizar registro
func update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var em model.EasyMoney
	if err := c.ShouldBindJSON(&em); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := service.Update(uint(id), &em)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Eliminar registro
func delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	if err := service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Obtener todos
func getAll(c *gin.Context) {
	results, err := service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// Diferencia total entre fechas
func getDifference(c *gin.Context) {
	start, end, err := parseDates(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total, err := service.CalculateDifference(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total.Text('f', 2)})
}

// Buscar por rango de fechas
func getByDateRange(c *gin.Context) {
	start, end, err := parseDates(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := service.FindByDateRange(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// Buscar por descripción
func searchByDescription(c *gin.Context) {
	keyword := c.Query("q")
	results, err := service.SearchByDescription(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// Diferencias consecutivas
func getConsecutiveDifferences(c *gin.Context) {
	start, end, err := parseDates(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sortBy := c.DefaultQuery("sortBy", "fromDate")
	differences, err := service.CalculateConsecutiveDifferences(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch sortBy {
	case "difference":
		dto.SortByDifferenceDesc(differences)
	case "toAmount":
		dto.SortByToAmountAsc(differences)
	default:
		dto.SortByFromDateAsc(differences)
	}

	c.JSON(http.StatusOK, differences)
}

// Diferencia total entre primera y última entrada
func getTotalMoneyDifference(c *gin.Context) {
	start, end, err := parseDates(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := service.CalculateTotalMoneyDifference(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Helper para parsear fechas
func parseDates(c *gin.Context) (time.Time, time.Time, error) {
	layout := "2006-01-02"
	startStr := c.Query("start")
	endStr := c.Query("end")

	start, err := time.Parse(layout, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err := time.Parse(layout, endStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}
