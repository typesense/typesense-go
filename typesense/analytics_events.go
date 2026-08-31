package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type AnalyticsEventsInterface interface {
	// Create an analytics event.
	//
	// Submit a single analytics event. The event must correspond to an existing analytics rule by name.
	//
	// HTTP: POST /analytics/events
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	Create(ctx context.Context, eventSchema *api.AnalyticsEvent) (*api.AnalyticsEventCreateResponse, error)
	// Retrieve analytics events.
	//
	// Retrieve the most recent events for a user and rule.
	//
	// HTTP: GET /analytics/events
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	Retrieve(ctx context.Context, params *api.GetAnalyticsEventsParams) (*api.AnalyticsEventsResponse, error)
}

type analyticsEvents struct {
	apiClient APIClientInterface
}

// Create an analytics event.
//
// Submit a single analytics event. The event must correspond to an existing analytics rule by name.
//
// HTTP: POST /analytics/events
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
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

// Retrieve analytics events.
//
// Retrieve the most recent events for a user and rule.
//
// HTTP: GET /analytics/events
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
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
