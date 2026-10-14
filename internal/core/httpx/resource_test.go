package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"asset-control/internal/core/dto"

	"github.com/gin-gonic/gin"
)

type listRequestDTO struct {
	dto.PaginatedRequestDTO
	Name        *string  `form:"name"`
	EmployeeIDs []uint64 `form:"employee_id" binding:"dive,gt=0"`
}

type itemDTO struct {
	Name string `json:"name"`
}

func serveResource(method, target, body string) (*httptest.ResponseRecorder, listRequestDTO) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	var received listRequestDTO
	router.GET("/items", func(c *gin.Context) {
		List(c, []string{"id", "name"}, func(_ context.Context, request listRequestDTO) (dto.PaginatedResponseDTO[itemDTO], error) {
			received = request
			return dto.NewPaginatedResponseDTO[itemDTO](nil, request.Page, request.EffectivePageSize(20), 0), nil
		})
	})
	router.GET("/items/:id", func(c *gin.Context) {
		Get(c, func(_ context.Context, id uint64) (itemDTO, error) { return itemDTO{Name: "item"}, nil })
	})
	router.POST("/items", func(c *gin.Context) {
		Create(c, func(_ context.Context, input itemDTO) (itemDTO, error) { return input, nil })
	})
	router.PUT("/items/:id", func(c *gin.Context) {
		Update(c, func(_ context.Context, id uint64, input itemDTO) (itemDTO, error) { return input, nil })
	})
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder, received
}

func TestListBindsPaginationAndFilters(t *testing.T) {
	recorder, request := serveResource(http.MethodGet, "/items?page=2&page_size=5&sort_by=name&direction=desc&name=an&employee_id=3&employee_id=4&ignored=1", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("%d %s", recorder.Code, recorder.Body.String())
	}
	if request.Page != 2 || request.PageSize != 5 || request.SortBy != "name" || request.Direction != dto.DirectionDesc ||
		*request.Name != "an" || !reflect.DeepEqual(request.EmployeeIDs, []uint64{3, 4}) {
		t.Fatalf("unexpected request %+v", request)
	}
	want := `{"content":[],"paging":{"total_elements":0,"size":5,"page":2,"total_pages":0,"has_content":false,"has_next":false,"has_previous":true}}`
	if recorder.Body.String() != want {
		t.Fatal(recorder.Body.String())
	}
	_, request = serveResource(http.MethodGet, "/items", "")
	if request.Name != nil || request.EmployeeIDs != nil {
		t.Fatalf("absent filters must stay nil: %+v", request)
	}
}

func TestListRejectsInvalidQueries(t *testing.T) {
	for target, wantMessage := range map[string]string{
		"/items?page=-1":        "page must satisfy gte=0",
		"/items?page_size=101":  "page_size must satisfy lte=100",
		"/items?direction=up":   "direction must be one of the following values: ASC, DESC",
		"/items?sort_by=secret": "sort_by must be one of the following values: id, name",
		"/items?employee_id=0":  "employee_id[0] must satisfy gt=0",
		"/items?page=abc":       "query parameters are invalid",
	} {
		recorder, _ := serveResource(http.MethodGet, target, "")
		if recorder.Code != http.StatusBadRequest || !strings.Contains(recorder.Body.String(), `"code":"CORE_001"`) ||
			!strings.Contains(recorder.Body.String(), wantMessage) {
			t.Errorf("%s: %d %s", target, recorder.Code, recorder.Body.String())
		}
	}
}

func TestSingleRecordHandlers(t *testing.T) {
	for _, tc := range []struct {
		method, target, body string
		status               int
		response             string
	}{
		{http.MethodGet, "/items/7", "", http.StatusOK, `{"name":"item"}`},
		{http.MethodGet, "/items/0", "", http.StatusBadRequest, `{"message":"invalid id","code":"CORE_002","logs":{}}`},
		{http.MethodPost, "/items", `{"name":"new"}`, http.StatusCreated, `{"name":"new"}`},
		{http.MethodPost, "/items", `{broken`, http.StatusBadRequest, `{"message":"Malformed JSON body","code":"CORE_002","logs":{}}`},
		{http.MethodPut, "/items/7", `{"name":"edited"}`, http.StatusOK, `{"name":"edited"}`},
		{http.MethodPut, "/items/x", `{"name":"edited"}`, http.StatusBadRequest, `{"message":"invalid id","code":"CORE_002","logs":{}}`},
	} {
		recorder, _ := serveResource(tc.method, tc.target, tc.body)
		if recorder.Code != tc.status || recorder.Body.String() != tc.response {
			t.Errorf("%s %s: %d %s", tc.method, tc.target, recorder.Code, recorder.Body.String())
		}
	}
}
