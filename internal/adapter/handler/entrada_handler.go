package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/v1c-g4b/diario-obras/internal/application"
	"github.com/v1c-g4b/diario-obras/internal/domain/entity"
)

type EntradaHandler struct {
	service *application.EntradaService
}

// CreateEntradaRequest representa o body para criar uma entrada.
type CreateEntradaRequest struct {
	ResponsavelID           uuid.UUID                `json:"responsavel_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Observacoes             string                   `json:"observacoes" example:"Fundação concluída sem intercorrências"`
	Etapa                   string                   `json:"etapa" binding:"required" example:"Fundação"`
	ProgressoEtapa          float64                  `json:"progresso_etapa" example:"75.5"`
	CustoDia                float64                  `json:"custo_dia" example:"15000.00"`
	QuantidadeTrabalhadores int64                    `json:"quantidade_trabalhadores" binding:"required" example:"12"`
	CondicoesClimaticas     entity.CondicaoClimatica `json:"condicoes_climaticas" example:"0"`
	Paralisacao             bool                     `json:"paralisacao" example:"false"`
}

func NewEntradaHandler(service *application.EntradaService) *EntradaHandler {
	return &EntradaHandler{service: service}
}

// Create godoc
// @Summary      Criar entrada
// @Description  Cria uma nova entrada (registro diário) para uma obra
// @Tags         Entradas
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "ID da obra"  format(uuid)
// @Param        body  body      CreateEntradaRequest  true  "Dados da entrada"
// @Success      201   {object}  entity.Entrada
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /obras/{id}/entradas [post]
func (h *EntradaHandler) Create(c *gin.Context) {
	obraID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da obra inválido"})
		return
	}

	var req CreateEntradaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entrada := &entity.Entrada{
		ResponsavelID:           req.ResponsavelID,
		Observacoes:             req.Observacoes,
		Etapa:                   req.Etapa,
		ProgressoEtapa:          req.ProgressoEtapa,
		CustoDia:                req.CustoDia,
		QuantidadeTrabalhadores: req.QuantidadeTrabalhadores,
		CondicoesClimaticas:     req.CondicoesClimaticas,
		Paralisacao:             req.Paralisacao,
	}

	if err := h.service.Create(c.Request.Context(), entrada, obraID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entrada)
}

// FindByID godoc
// @Summary      Buscar entrada por ID
// @Description  Retorna uma entrada pelo seu ID, incluindo fotos e responsável
// @Tags         Entradas
// @Produce      json
// @Param        id   path      string  true  "ID da entrada"  format(uuid)
// @Success      200  {object}  entity.Entrada
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /entradas/{id} [get]
func (h *EntradaHandler) FindByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	entrada, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if entrada == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "não encontrado"})
		return
	}

	c.JSON(http.StatusOK, entrada)
}

// ListByObra godoc
// @Summary      Listar entradas de uma obra
// @Description  Retorna todas as entradas de uma obra específica
// @Tags         Entradas
// @Produce      json
// @Param        id   path      string  true  "ID da obra"  format(uuid)
// @Success      200  {array}   entity.Entrada
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /obras/{id}/entradas [get]
func (h *EntradaHandler) ListByObra(c *gin.Context) {
	obraID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da obra inválido"})
		return
	}

	entradas, err := h.service.ListByObra(c.Request.Context(), obraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entradas)
}

// Delete godoc
// @Summary      Deletar entrada
// @Description  Remove uma entrada e todas as suas fotos associadas
// @Tags         Entradas
// @Param        id   path  string  true  "ID da entrada"  format(uuid)
// @Success      204  "Sem conteúdo"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /entradas/{id} [delete]
func (h *EntradaHandler) Delete(c *gin.Context) {
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
