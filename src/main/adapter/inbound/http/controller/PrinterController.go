package controller

import (
	"net/http"

	"asset-control/src/main/domain/dto"
	"asset-control/src/main/usecase"

	"github.com/gin-gonic/gin"
)

type PrinterController struct {
	printerUseCase *usecase.PrinterUseCase
}

func NewPrinterController(printerUseCase *usecase.PrinterUseCase) *PrinterController {
	return &PrinterController{
		printerUseCase: printerUseCase,
	}
}

// CreatePrinter godoc
// @Summary      Cria uma nova impressora
// @Description  Cria uma nova impressora a partir dos dados enviados no corpo da requisição.
// @Tags         printers
// @Accept       json
// @Produce      json
// @Param        printer  body      dto.CreatePrinterDto  true  "Dados da impressora"
// @Success      201       {object}  map[string]interface{}  "Impressora criada com sucesso"
// @Failure      400       {object}  map[string]string  "Dados inválidos"
// @Failure      500       {object}  map[string]string  "Erro interno do servidor"
// @Router       /printers [post]
func (ctl *PrinterController) CreatePrinterController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var printer dto.CreatePrinterDto

		if err := c.ShouldBindJSON(&printer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		created, err := ctl.printerUseCase.CreatePrinter(c.Request.Context(), printer)
		if err != nil {
			writeError(c, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "printer created successfully",
			"data":   created,
		})
	}
}

// UpdatePrinter godoc
// @Summary      Atualiza uma impressora
// @Description  Atualiza os dados de uma impressora existente através do ID.
// @Tags         printers
// @Accept       json
// @Produce      json
// @Param        id       query     uint64                true  "ID da impressora"
// @Param        printer  body      dto.CreatePrinterDto  true  "Dados atualizados da impressora"
// @Success      200      {object}  map[string]interface{}  "Impressora atualizada com sucesso"
// @Failure      400      {object}  map[string]string  "ID inválido ou dados inválidos"
// @Failure      404      {object}  map[string]string  "Impressora não encontrada"
// @Router       /printers [put]
func (ctl *PrinterController) UpdatePrinterController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")

		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "query param 'id' is required",
			})
			return
		}

		id, err := parseID(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid 'id' value",
			})
			return
		}

		var printer dto.CreatePrinterDto

		if err := c.ShouldBindJSON(&printer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		updated, err := ctl.printerUseCase.UpdatePrinter(c.Request.Context(), id, printer)

		if err != nil {
			writeError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "printer updated successfully",
			"data":   updated,
		})
	}
}

// GetAllPrinters godoc
// @Summary      Lista todas as impressoras
// @Description  Retorna todas as impressoras cadastradas.
// @Tags         printers
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Impressoras encontradas com sucesso"
// @Failure      500  {object}  map[string]string  "Erro interno do servidor"
// @Router       /printers [get]
func (ctl *PrinterController) GetAllPrintersController() gin.HandlerFunc {
	return func(c *gin.Context) {
		printers, err := ctl.printerUseCase.GetAllPrinters(c.Request.Context())

		if err != nil {
			writeError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "printers retrieved successfully",
			"data":   printers,
		})
	}
}

// FindByID godoc
// @Summary      Busca impressora por ID
// @Description  Busca uma impressora pelo ID informado como query parameter.
// @Tags         printers
// @Produce      json
// @Param        id  query  uint64  true  "ID da impressora"
// @Success      200  {object}  map[string]interface{}  "Impressora encontrada com sucesso"
// @Failure      400  {object}  map[string]string  "ID inválido ou não informado"
// @Failure      404  {object}  map[string]string  "Impressora não encontrada"
// @Router       /printers/id [get]
func (ctl *PrinterController) FindByIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")

		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "query param 'id' is required",
			})
			return
		}

		id, err := parseID(idParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid 'id' value",
			})
			return
		}

		printer, err := ctl.printerUseCase.FindByID(c.Request.Context(), id)

		if err != nil {
			writeError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "printer found successfully",
			"data":   printer,
		})
	}
}
