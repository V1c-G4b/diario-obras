package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/application"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
)

type ResponsavelHandler struct {
	service *application.ResponsavelService
}

// CreateResponsavelRequest representa o body para criar um responsável.
type CreateResponsavelRequest struct {
	Nome  string `json:"nome" binding:"required" example:"João Silva"`
	Email string `json:"email" binding:"required" example:"joao@exemplo.com"`
	Cargo string `json:"cargo" binding:"required" example:"Engenheiro Civil"`
}

// UpdateResponsavelRequest representa o body para atualizar um responsável.
type UpdateResponsavelRequest struct {
	Nome  string `json:"nome" example:"João da Silva"`
	Email string `json:"email" example:"joao.silva@exemplo.com"`
	Cargo string `json:"cargo" example:"Engenheiro Chefe"`
}

func NewResponsavelHandler(service *application.ResponsavelService) *ResponsavelHandler {
	return &ResponsavelHandler{service: service}
}

// Create godoc
// @Summary      Criar responsável
// @Description  Cria um novo responsável
// @Tags         Responsáveis
// @Accept       json
// @Produce      json
// @Param        body  body      CreateResponsavelRequest  true  "Dados do responsável"
// @Success      201   {object}  entity.Responsavel
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /responsaveis [post]
func (h *ResponsavelHandler) Create(c *gin.Context) {
	var req CreateResponsavelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responsavel := &entity.Responsavel{
		Nome:  req.Nome,
		Email: req.Email,
		Cargo: req.Cargo,
	}

	if err := h.service.Create(c.Request.Context(), responsavel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, responsavel)
}

// FindByID godoc
// @Summary      Buscar responsável por ID
// @Description  Retorna um responsável pelo seu ID
// @Tags         Responsáveis
// @Produce      json
// @Param        id   path      string  true  "ID do responsável"  format(uuid)
// @Success      200  {object}  entity.Responsavel
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /responsaveis/{id} [get]
func (h *ResponsavelHandler) FindByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	responsavel, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if responsavel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "não encontrado"})
		return
	}

	c.JSON(http.StatusOK, responsavel)
}

// List godoc
// @Summary      Listar responsáveis
// @Description  Retorna todos os responsáveis
// @Tags         Responsáveis
// @Produce      json
// @Success      200  {array}   entity.Responsavel
// @Failure      500  {object}  map[string]string
// @Router       /responsaveis [get]
func (h *ResponsavelHandler) List(c *gin.Context) {
	responsaveis, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responsaveis)
}

// Update godoc
// @Summary      Atualizar responsável
// @Description  Atualiza um responsável existente
// @Tags         Responsáveis
// @Accept       json
// @Produce      json
// @Param        id    path      string                    true  "ID do responsável"  format(uuid)
// @Param        body  body      UpdateResponsavelRequest  true  "Dados para atualização"
// @Success      200   {object}  entity.Responsavel
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /responsaveis/{id} [put]
func (h *ResponsavelHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	responsavel, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Responsável não encontrado"})
		return
	}

	var req UpdateResponsavelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Nome != "" {
		responsavel.Nome = req.Nome
	}
	if req.Email != "" {
		responsavel.Email = req.Email
	}
	if req.Cargo != "" {
		responsavel.Cargo = req.Cargo
	}

	if err := h.service.Update(c.Request.Context(), responsavel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responsavel)
}

// Delete godoc
// @Summary      Deletar responsável
// @Description  Remove um responsável pelo ID
// @Tags         Responsáveis
// @Param        id   path  string  true  "ID do responsável"  format(uuid)
// @Success      204  "Sem conteúdo"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /responsaveis/{id} [delete]
func (h *ResponsavelHandler) Delete(c *gin.Context) {
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
