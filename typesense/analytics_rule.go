package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v3/typesense/api"
)

type AnalyticsRuleInterface interface {
	Delete(ctx context.Context) (*api.AnalyticsRuleDeleteResponse, error)
	Retrieve(ctx context.Context) (*api.AnalyticsRuleSchema, error)
}

type analyticsRule struct {
	apiClient APIClientInterface
	ruleName  string
}

func (a *analyticsRule) Delete(ctx context.Context) (*api.AnalyticsRuleDeleteResponse, error) {
	response, err := a.apiClient.DeleteAnalyticsRuleWithResponse(ctx, a.ruleName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (a *analyticsRule) Retrieve(ctx context.Context) (*api.AnalyticsRuleSchema, error) {
	response, err := a.apiClient.RetrieveAnalyticsRuleWithResponse(ctx, a.ruleName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
