package tests

import (
	"testing"

	"github.com/ayush-sr/score-keeper/backend/internal/dto"
)

func TestSuccess_ReturnsDataAndNoError(t *testing.T) {
	resp := dto.Success("hello")
	if resp.Data != "hello" {
		t.Errorf("expected 'hello', got %v", resp.Data)
	}
	if resp.Error != nil {
		t.Error("expected nil error")
	}
	if resp.Meta != nil {
		t.Error("expected nil meta")
	}
}

func TestSuccessWithMeta_ReturnsMeta(t *testing.T) {
	meta := &dto.Meta{Page: 2, PerPage: 10, Total: 50}
	resp := dto.SuccessWithMeta([]string{"a", "b"}, meta)
	if resp.Meta == nil {
		t.Fatal("expected non-nil meta")
	}
	if resp.Meta.Page != 2 || resp.Meta.Total != 50 || resp.Meta.PerPage != 10 {
		t.Errorf("unexpected meta: %+v", resp.Meta)
	}
}

func TestErrorResponse_ReturnsNilDataAndError(t *testing.T) {
	resp := dto.ErrorResponse("NOT_FOUND", "user not found")
	if resp.Data != nil {
		t.Error("expected nil data")
	}
	if resp.Error == nil {
		t.Fatal("expected non-nil error")
	}
	if resp.Error.Code != "NOT_FOUND" {
		t.Errorf("expected NOT_FOUND, got %s", resp.Error.Code)
	}
	if resp.Error.Message != "user not found" {
		t.Errorf("expected 'user not found', got %s", resp.Error.Message)
	}
}

func TestSuccess_NilData(t *testing.T) {
	resp := dto.Success(nil)
	if resp.Data != nil {
		t.Error("expected nil data")
	}
}

func TestSuccess_StructData(t *testing.T) {
	type payload struct {
		ID   int
		Name string
	}
	resp := dto.Success(payload{ID: 1, Name: "test"})
	p, ok := resp.Data.(payload)
	if !ok {
		t.Fatal("expected payload type")
	}
	if p.ID != 1 || p.Name != "test" {
		t.Errorf("unexpected: %+v", p)
	}
}

func TestSuccessWithMeta_NilMeta(t *testing.T) {
	resp := dto.SuccessWithMeta("data", nil)
	if resp.Meta != nil {
		t.Error("expected nil meta when nil passed")
	}
}
