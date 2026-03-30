package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes(r *gin.Engine, obraHandler *ObraHandler, entradaHandler *EntradaHandler, responsavelHandler *ResponsavelHandler, FotoHandler *FotoHandler) {
	api := r.Group("/api/v1")
	{
		api.GET("/metrics", gin.WrapH(promhttp.Handler()))
		obras := api.Group("/obras")
		{
			obras.POST("", obraHandler.Create)
			obras.GET("", obraHandler.List)
			obras.GET("/:id", obraHandler.FindByID)
			obras.PUT("/:id", obraHandler.Update)
			obras.DELETE("/:id", obraHandler.Delete)

			obras.POST("/:id/entradas", entradaHandler.Create)
			obras.GET("/:id/entradas", entradaHandler.ListByObra)
		}

		entradas := api.Group("/entradas")
		{
			entradas.GET("/:id", entradaHandler.FindByID)
			entradas.DELETE("/:id", entradaHandler.Delete)
			entradas.POST("/:id/fotos", FotoHandler.Create)
			entradas.GET("/:id/fotos", FotoHandler.ListByEntrada)
			entradas.DELETE("/:id/fotos/:fotoId", FotoHandler.Delete)
		}

		responsaveis := api.Group("/responsaveis")
		{
			responsaveis.POST("", responsavelHandler.Create)
			responsaveis.GET("", responsavelHandler.List)
			responsaveis.GET("/:id", responsavelHandler.FindByID)
			responsaveis.PUT("/:id", responsavelHandler.Update)
			responsaveis.DELETE("/:id", responsavelHandler.Delete)
		}
	}
}
