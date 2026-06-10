package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type PresetInterface interface {
	// Retrieves a preset.
	//
	// Retrieve the details of a preset, given it's name.
	//
	// HTTP: GET /presets/{presetId}
	//
	// See: https://typesense.org/docs/latest/api/search.html#presets
	Retrieve(ctx context.Context) (*api.PresetSchema, error)
	// Delete a preset.
	//
	// Permanently deletes a preset, given it's name.
	//
	// HTTP: DELETE /presets/{presetId}
	//
	// See: https://typesense.org/docs/latest/api/search.html#presets
	Delete(ctx context.Context) (*api.PresetDeleteSchema, error)
}

type preset struct {
	apiClient  APIClientInterface
	presetName string
}

// Retrieves a preset.
//
// Retrieve the details of a preset, given it's name.
//
// HTTP: GET /presets/{presetId}
//
// See: https://typesense.org/docs/latest/api/search.html#presets
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

// Delete a preset.
//
// Permanently deletes a preset, given it's name.
//
// HTTP: DELETE /presets/{presetId}
//
// See: https://typesense.org/docs/latest/api/search.html#presets
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
