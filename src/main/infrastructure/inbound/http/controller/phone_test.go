package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"asset-control/src/main/application/ports"
	"asset-control/src/main/application/usecase"
	"asset-control/src/main/domain/dto"
)

type phoneStub struct {
	ports.PhonePorts
	input  dto.CreatePhoneDto
	called bool
}

func (s *phoneStub) CreatePhone(_ context.Context, input dto.CreatePhoneDto) (dto.PhoneDto, error) {
	s.input, s.called = input, true
	return dto.PhoneDto{ID: 17, Equipament: input.Equipament, Model: input.Model, IMEI: input.IMEI, SerialNumber: input.SerialNumber, HasChip: input.HasChip, Situation: input.Situation}, nil
}

func (s *phoneStub) UpdatePhone(ctx context.Context, _ uint64, input dto.CreatePhoneDto) (dto.PhoneDto, error) {
	return s.CreatePhone(ctx, input)
}

func TestPhoneCreateAndUpdateWithoutDeviceIdentification(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut} {
		for _, equipment := range []string{"Só chip", "Sem chip", "Celular", ""} {
			t.Run(method+"/"+equipment, func(t *testing.T) {
				stub := &phoneStub{}
				ctl := NewPhoneController(usecase.NewPhoneUseCase(stub))
				handler := ctl.CreatePhoneController()
				status := http.StatusCreated
				if method == http.MethodPut {
					handler, status = ctl.UpdatePhoneController(), http.StatusOK
				}
				payload := map[string]any{"situation": "Em uso"}
				if equipment != "" {
					payload["equipament"] = equipment
				}
				body, err := json.Marshal(payload)
				if err != nil {
					t.Fatal(err)
				}
				if equipment == "Celular" || equipment == "" {
					status = http.StatusBadRequest
				}
				result := request(handler, method, "/?id=17", string(body), context.Background())
				if result.Code != status {
					t.Fatalf("got %d: %s", result.Code, result.Body.String())
				}
				if status == http.StatusBadRequest {
					if stub.called {
						t.Fatal("invalid request reached repository")
					}
					return
				}
				if !stub.called || stub.input.Model != "" || stub.input.IMEI != "" || stub.input.SerialNumber != "" {
					t.Fatalf("unexpected saved input: %+v", stub.input)
				}
				var response struct {
					Data dto.PhoneDto `json:"data"`
				}
				if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				if response.Data.ID != 17 || response.Data.Equipament == nil || *response.Data.Equipament != equipment {
					t.Fatalf("unexpected response: %+v", response.Data)
				}
			})
		}
	}
}
