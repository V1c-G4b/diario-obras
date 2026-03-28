package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/application"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
)

const maxFileSize = 10 << 20 // 10 MB

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/heic": true,
	"image/heif": true,
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".heic": true,
	".heif": true,
}

type FotoHandler struct {
	service *application.FotoService
}

func NewFotoHandler(service *application.FotoService) *FotoHandler {
	return &FotoHandler{service: service}
}

// Create godoc
// @Summary      Upload de foto
// @Description  Faz upload de uma foto para uma entrada. Aceita jpg, png, webp, heic (máx 10MB)
// @Tags         Fotos
// @Accept       multipart/form-data
// @Produce      json
// @Param        id         path      string  true  "ID da entrada"  format(uuid)
// @Param        file       formData  file    true  "Arquivo de imagem"
// @Param        descricao  formData  string  false "Descrição da foto"
// @Success      201        {object}  entity.Foto
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /entradas/{id}/fotos [post]
func (h *FotoHandler) Create(c *gin.Context) {
	entradaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da entrada inválido"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "arquivo obrigatório"})
		return
	}

	if fileHeader.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "arquivo excede o tamanho máximo de 10MB"})
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo de arquivo não permitido, envie uma imagem (jpg, png, webp, heic)"})
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "" && !allowedMimeTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo MIME não permitido, envie uma imagem"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao ler arquivo"})
		return
	}
	defer file.Close()

	descricao := c.PostForm("descricao")
	fileName := fmt.Sprintf("%s/%s%s", entradaID, uuid.New(), ext)

	foto := &entity.Foto{
		EntradaID: entradaID,
		Descricao: descricao,
	}

	if err := h.service.Create(c.Request.Context(), foto, fileName, file, fileHeader.Size); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, foto)
}

// ListByEntrada godoc
// @Summary      Listar fotos de uma entrada
// @Description  Retorna todas as fotos de uma entrada específica
// @Tags         Fotos
// @Produce      json
// @Param        id   path      string  true  "ID da entrada"  format(uuid)
// @Success      200  {array}   entity.Foto
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /entradas/{id}/fotos [get]
func (h *FotoHandler) ListByEntrada(c *gin.Context) {
	entradaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da entrada inválido"})
		return
	}

	fotos, err := h.service.ListByEntrada(c.Request.Context(), entradaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fotos)
}

// Delete godoc
// @Summary      Deletar foto
// @Description  Remove uma foto do banco e do storage
// @Tags         Fotos
// @Param        id      path  string  true  "ID da entrada"  format(uuid)
// @Param        fotoId  path  string  true  "ID da foto"     format(uuid)
// @Success      204     "Sem conteúdo"
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /entradas/{id}/fotos/{fotoId} [delete]
func (h *FotoHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("fotoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da foto inválido"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
