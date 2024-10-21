package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v2/typesense/api"
)

type AnalyticsEventsInterface interface {
	Create(ctx context.Context, eventSchema *api.AnalyticsEventCreateSchema) (*api.AnalyticsEventCreateResponse, error)
}

type analyticsEvents struct {
	apiClient APIClientInterface
}

func (a *analyticsEvents) Create(ctx context.Context, eventSchema *api.AnalyticsEventCreateSchema) (*api.AnalyticsEventCreateResponse, error) {
	response, err := a.apiClient.CreateAnalyticsEventWithResponse(ctx, api.CreateAnalyticsEventJSONRequestBody(*eventSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON201 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON201, nil
}
