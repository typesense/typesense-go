package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type AnalyticsEventsInterface interface {
	Create(ctx context.Context, eventSchema *api.AnalyticsEvent) (*api.AnalyticsEventCreateResponse, error)
	Retrieve(ctx context.Context, params *api.GetAnalyticsEventsParams) (*api.AnalyticsEventsResponse, error)
}

type analyticsEvents struct {
	apiClient APIClientInterface
}

func (a *analyticsEvents) Create(ctx context.Context, eventSchema *api.AnalyticsEvent) (*api.AnalyticsEventCreateResponse, error) {
	response, err := a.apiClient.CreateAnalyticsEventWithResponse(ctx, api.CreateAnalyticsEventJSONRequestBody(*eventSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (a *analyticsEvents) Retrieve(ctx context.Context, params *api.GetAnalyticsEventsParams) (*api.AnalyticsEventsResponse, error) {
	response, err := a.apiClient.GetAnalyticsEventsWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
