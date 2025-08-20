package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mtripode101/easymoney-go/config"
	"github.com/mtripode101/easymoney-go/controller"
	"github.com/mtripode101/easymoney-go/model"
)

func main() {
	// Inicializar conexión a la base de datos
	config.InitDB()

	// Migrar el modelo EasyMoney si no existe
	config.DB.AutoMigrate(&model.EasyMoney{})

	// Crear router Gin
	r := gin.Default()

	// Cargar vistas HTML desde carpeta templates/
	r.LoadHTMLGlob("templates/*.html")

	// Registrar rutas REST
	controller.RegisterRoutes(r)

	// Registrar rutas web
	controller.RegisterWebRoutes(r)

	// Iniciar servidor en puerto 8080
	r.Run(":8080")
}
