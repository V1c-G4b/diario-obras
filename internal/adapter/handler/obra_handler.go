package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/application"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
)

type ObraHandler struct {
	service *application.ObraService
}

// CreateObraRequest representa o body para criar uma obra.
type CreateObraRequest struct {
	Nome         string    `json:"nome" binding:"required" example:"Residencial Aurora"`
	Endereco     string    `json:"endereco" binding:"required" example:"Rua das Flores, 123"`
	DataInicio   time.Time `json:"data_inicio" binding:"required" example:"2026-04-01T00:00:00Z"`
	DataEstimada time.Time `json:"data_estimada" binding:"required" example:"2026-12-31T00:00:00Z"`
}

// UpdateObraRequest representa o body para atualizar uma obra.
type UpdateObraRequest struct {
	Nome         string    `json:"nome" example:"Residencial Aurora II"`
	Endereco     string    `json:"endereco" example:"Rua das Flores, 456"`
	DataInicio   time.Time `json:"data_inicio" example:"2026-04-01T00:00:00Z"`
	DataEstimada time.Time `json:"data_estimada" example:"2027-06-30T00:00:00Z"`
}

func NewObraHandler(service *application.ObraService) *ObraHandler {
	return &ObraHandler{service: service}
}

// Create godoc
// @Summary      Criar obra
// @Description  Cria uma nova obra
// @Tags         Obras
// @Accept       json
// @Produce      json
// @Param        body  body      CreateObraRequest  true  "Dados da obra"
// @Success      201   {object}  entity.Obra
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /obras [post]
func (h *ObraHandler) Create(c *gin.Context) {
	var req CreateObraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	obra := &entity.Obra{
		Nome:         req.Nome,
		Endereco:     req.Endereco,
		DataInicio:   req.DataInicio,
		DataEstimada: req.DataEstimada,
	}

	if err := h.service.Create(c.Request.Context(), obra); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, obra)
}

// FindByID godoc
// @Summary      Buscar obra por ID
// @Description  Retorna uma obra pelo seu ID
// @Tags         Obras
// @Produce      json
// @Param        id   path      string  true  "ID da obra"  format(uuid)
// @Success      200  {object}  entity.Obra
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /obras/{id} [get]
func (h *ObraHandler) FindByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	obra, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if obra == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "não encontrado"})
		return
	}

	c.JSON(http.StatusOK, obra)
}

// List godoc
// @Summary      Listar obras
// @Description  Retorna todas as obras
// @Tags         Obras
// @Produce      json
// @Success      200  {array}   entity.Obra
// @Failure      500  {object}  map[string]string
// @Router       /obras [get]
func (h *ObraHandler) List(c *gin.Context) {
	obras, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, obras)
}

// Update godoc
// @Summary      Atualizar obra
// @Description  Atualiza uma obra existente
// @Tags         Obras
// @Accept       json
// @Produce      json
// @Param        id    path      string            true  "ID da obra"  format(uuid)
// @Param        body  body      UpdateObraRequest  true  "Dados para atualização"
// @Success      200   {object}  entity.Obra
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /obras/{id} [put]
func (h *ObraHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	obra, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Obra não encontrada"})
		return
	}

	var req UpdateObraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Nome != "" {
		obra.Nome = req.Nome
	}
	if req.Endereco != "" {
		obra.Endereco = req.Endereco
	}
	if !req.DataInicio.IsZero() {
		obra.DataInicio = req.DataInicio
	}
	if !req.DataEstimada.IsZero() {
		obra.DataEstimada = req.DataEstimada
	}

	if err := h.service.Update(c.Request.Context(), obra); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, obra)
}

// Delete godoc
// @Summary      Deletar obra
// @Description  Remove uma obra pelo ID
// @Tags         Obras
// @Param        id   path  string  true  "ID da obra"  format(uuid)
// @Success      204  "Sem conteúdo"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /obras/{id} [delete]
func (h *ObraHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
