package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type AnalyticsRuleInterface interface {
	Delete(ctx context.Context) (*api.AnalyticsRule, error)
	Retrieve(ctx context.Context) (*api.AnalyticsRule, error)
	Update(ctx context.Context, ruleSchema *api.AnalyticsRuleUpdate) (*api.AnalyticsRule, error)
}

type analyticsRule struct {
	apiClient APIClientInterface
	ruleName  string
}

func (a *analyticsRule) Delete(ctx context.Context) (*api.AnalyticsRule, error) {
	response, err := a.apiClient.DeleteAnalyticsRuleWithResponse(ctx, a.ruleName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (a *analyticsRule) Retrieve(ctx context.Context) (*api.AnalyticsRule, error) {
	response, err := a.apiClient.RetrieveAnalyticsRuleWithResponse(ctx, a.ruleName)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (a *analyticsRule) Update(ctx context.Context, ruleSchema *api.AnalyticsRuleUpdate) (*api.AnalyticsRule, error) {
	response, err := a.apiClient.UpsertAnalyticsRuleWithResponse(ctx, a.ruleName, api.UpsertAnalyticsRuleJSONRequestBody(*ruleSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
