package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type PresetsInterface interface {
	Retrieve(ctx context.Context) ([]*api.PresetSchema, error)
	Upsert(ctx context.Context, presetName string, presetValue *api.PresetUpsertSchema) (*api.PresetSchema, error)
}

type presets struct {
	apiClient APIClientInterface
}

func (p *presets) Retrieve(ctx context.Context) ([]*api.PresetSchema, error) {
	response, err := p.apiClient.RetrieveAllPresetsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200.Presets, nil
}

func (p *presets) Upsert(ctx context.Context, presetName string, presetValue *api.PresetUpsertSchema) (*api.PresetSchema, error) {
	response, err := p.apiClient.UpsertPresetWithResponse(ctx, presetName, *presetValue)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
