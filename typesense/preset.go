package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type PresetInterface interface {
	Retrieve(ctx context.Context) (*api.PresetSchema, error)
	Delete(ctx context.Context) (*api.PresetDeleteSchema, error)
}

type preset struct {
	apiClient  APIClientInterface
	presetName string
}

func (p *preset) Retrieve(ctx context.Context) (*api.PresetSchema, error) {
	response, err := p.apiClient.RetrievePresetWithResponse(ctx, p.presetName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (p *preset) Delete(ctx context.Context) (*api.PresetDeleteSchema, error) {
	response, err := p.apiClient.DeletePresetWithResponse(ctx, p.presetName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
