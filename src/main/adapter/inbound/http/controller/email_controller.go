package controller

import (
	"math"
	"net/http"
	"strings"

	"asset-control/src/main/domain/dto"
	"asset-control/src/main/usecase"

	"github.com/gin-gonic/gin"
)

type EmailController struct {
	emailUseCase *usecase.EmailUseCase
}

func NewEmailController(emailUseCase *usecase.EmailUseCase) *EmailController {
	return &EmailController{emailUseCase: emailUseCase}
}

// GetAllEmails godoc
// @Summary      Lista todos os emails
// @Description  Retorna todos os registros de usableEmails em ordem crescente de ID, incluindo password_hash em texto puro, exatamente como armazenado. Retorna data como uma lista vazia quando não houver registros.
// @Tags         emails
// @Produce      json
// @Success      200  {object}  dto.EmailListResponseDto
// @Failure      500  {object}  map[string]string  "Erro interno do servidor"
// @Failure      504  {object}  map[string]string  "Tempo limite da operação no banco de dados"
// @Router       /emails [get]
func (ctl *EmailController) GetAllEmailsController() gin.HandlerFunc {
	return func(c *gin.Context) {
		emails, err := ctl.emailUseCase.GetAllEmails(c.Request.Context())
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto.EmailListResponseDto{Status: "emails retrieved successfully", Data: emails})
	}
}

// CreateEmail godoc
// @Summary      Cria um email
// @Description  Cria um registro em usableEmails com ID int32 e created_at gerados pelo servidor. Os campos user, email e password_hash são obrigatórios. A senha é armazenada e retornada em texto puro, sem alterações.
// @Tags         emails
// @Accept       json
// @Produce      json
// @Param        email  body  dto.CreateEmailDto  true  "Dados do email"
// @Success      201  {object}  dto.EmailResponseDto  "Email criado; status: created"
// @Failure      400  {object}  map[string]string  "Dados inválidos"
// @Failure      409  {object}  map[string]string  "Conflito de registro"
// @Failure      500  {object}  map[string]string  "Erro interno do servidor ou IDs esgotados"
// @Failure      504  {object}  map[string]string  "Tempo limite da operação no banco de dados"
// @Router       /emails [post]
func (ctl *EmailController) CreateEmailController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input dto.CreateEmailDto
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		created, err := ctl.emailUseCase.CreateEmail(c.Request.Context(), input)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusCreated, dto.EmailResponseDto{Status: "created", Data: created})
	}
}

// UpdateEmail godoc
// @Summary      Atualiza um email
// @Description  Substitui user, email e password_hash do registro indicado pelo ID; os três campos são obrigatórios. Preserva id e created_at. A senha é armazenada e retornada em texto puro, sem alterações. Não cria registros quando o ID não existe.
// @Tags         emails
// @Accept       json
// @Produce      json
// @Param        id     query  int  true  "ID do email (inteiro positivo de 32 bits)" Format(int32) minimum(1) maximum(2147483647)
// @Param        email  body  dto.CreateEmailDto  true  "Dados atualizados do email"
// @Success      200  {object}  dto.EmailResponseDto  "Email atualizado; status: email updated successfully"
// @Failure      400  {object}  map[string]string  "ID ou dados inválidos"
// @Failure      404  {object}  map[string]string  "Email não encontrado"
// @Failure      409  {object}  map[string]string  "Conflito de registro"
// @Failure      500  {object}  map[string]string  "Erro interno do servidor"
// @Failure      504  {object}  map[string]string  "Tempo limite da operação no banco de dados"
// @Router       /emails [put]
func (ctl *EmailController) UpdateEmailController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")
		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'id' is required"})
			return
		}
		id, err := parseID(idParam)
		if err != nil || id > math.MaxInt32 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'id' value"})
			return
		}
		var input dto.CreateEmailDto
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updated, err := ctl.emailUseCase.UpdateEmail(c.Request.Context(), int32(id), input)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto.EmailResponseDto{Status: "email updated successfully", Data: updated})
	}
}

// FindEmailByID godoc
// @Summary      Busca email por ID
// @Description  Retorna um email de usableEmails pelo ID. O campo password_hash retorna a senha em texto puro, exatamente como armazenada.
// @Tags         emails
// @Produce      json
// @Param        id   query     int  true  "ID do email (inteiro positivo de 32 bits)" Format(int32) minimum(1) maximum(2147483647)
// @Success      200  {object}  dto.EmailResponseDto
// @Failure      400  {object}  map[string]string  "ID inválido ou não informado"
// @Failure      404  {object}  map[string]string  "Email não encontrado"
// @Failure      500  {object}  map[string]string  "Erro interno do servidor"
// @Failure      504  {object}  map[string]string  "Tempo limite da operação no banco de dados"
// @Router       /emails/id [get]
func (ctl *EmailController) FindByIDController() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Query("id")
		if idParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'id' is required"})
			return
		}
		id, err := parseID(idParam)
		if err != nil || id > math.MaxInt32 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'id' value"})
			return
		}
		email, err := ctl.emailUseCase.FindByID(c.Request.Context(), int32(id))
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto.EmailResponseDto{Status: "email found successfully", Data: email})
	}
}

// FindEmailByUser godoc
// @Summary      Busca email por usuário
// @Description  Busca parcial pelo campo user, sem diferenciar maiúsculas de minúsculas; caracteres especiais são tratados literalmente. Espaços nas extremidades da busca são removidos. Retorna o email com menor ID quando houver múltiplos resultados. O campo password_hash retorna a senha em texto puro, exatamente como armazenada.
// @Tags         emails
// @Produce      json
// @Param        user  query     string  true  "Usuário do email (não pode ser vazio)"
// @Success      200   {object}  dto.EmailResponseDto
// @Failure      400   {object}  map[string]string  "Usuário não informado ou vazio"
// @Failure      404   {object}  map[string]string  "Email não encontrado"
// @Failure      500   {object}  map[string]string  "Erro interno do servidor"
// @Failure      504   {object}  map[string]string  "Tempo limite da operação no banco de dados"
// @Router       /emails/user [get]
func (ctl *EmailController) FindByUserController() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := strings.TrimSpace(c.Query("user"))
		if user == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is required"})
			return
		}
		email, err := ctl.emailUseCase.FindByUser(c.Request.Context(), user)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto.EmailResponseDto{Status: "email found successfully", Data: email})
	}
}
