package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/v1c-g4b/diario-obras/internal/adapter/handler"
	"github.com/v1c-g4b/diario-obras/internal/adapter/repository"
	"github.com/v1c-g4b/diario-obras/internal/adapter/storage"
	"github.com/v1c-g4b/diario-obras/internal/application"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/v1c-g4b/diario-obras/docs"
)

// @title           Diário de Obras API
// @version         1.0
// @description     API para gerenciamento de diário de obras, entradas diárias, fotos e responsáveis.

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado, utilizando variáveis de ambiente de sistema")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados: ", err)
	}

	if err := db.AutoMigrate(&entity.Obra{}, &entity.Entrada{}, entity.Foto{}, entity.Responsavel{}); err != nil {
		log.Fatal("Falha na migração do banco: ", err)
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_USER")
	secretKey := os.Getenv("MINIO_PASSWORD")
	bucketName := os.Getenv("MINIO_BUCKET")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	storage := storage.NewStorage(client, bucketName)

	// Repositories
	obraRepo := repository.NewObraGormRepository(db)
	entradaRepo := repository.NewEntradaGormRepository(db)
	responsavelRepo := repository.NewResponsavelGormRepository(db)
	fotoRepo := repository.NewFotoGormRepository(db)

	// Services
	obraService := application.NewObraService(obraRepo)
	fotoService := application.NewFotoService(fotoRepo, storage)
	entradaService := application.NewEntradaService(entradaRepo, obraRepo, fotoService, responsavelRepo)
	responsavelService := application.NewResponsavelService(responsavelRepo)

	// Handlers
	obraHandler := handler.NewObraHandler(obraService)
	entradaHandler := handler.NewEntradaHandler(entradaService)
	responsavelHandler := handler.NewResponsavelHandler(responsavelService)
	fotoHandler := handler.NewFotoHandler(fotoService)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	handler.SetupRoutes(router, obraHandler, entradaHandler, responsavelHandler, fotoHandler)

	if err := router.Run(); err != nil {
		log.Fatal("Falha ao iniciar o servidor: ", err)
	}
}
