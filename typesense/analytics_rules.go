package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type AnalyticsRulesInterface interface {
	Upsert(ctx context.Context, ruleName string, ruleSchema *api.AnalyticsRuleUpsertSchema) (*api.AnalyticsRuleSchema, error)
	Retrieve(ctx context.Context) ([]*api.AnalyticsRuleSchema, error)
}

type analyticsRules struct {
	apiClient APIClientInterface
}

func (a *analyticsRules) Upsert(ctx context.Context, ruleName string, ruleSchema *api.AnalyticsRuleUpsertSchema) (*api.AnalyticsRuleSchema, error) {
	response, err := a.apiClient.UpsertAnalyticsRuleWithResponse(ctx,
		ruleName, api.UpsertAnalyticsRuleJSONRequestBody(*ruleSchema))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}

func (a *analyticsRules) Retrieve(ctx context.Context) ([]*api.AnalyticsRuleSchema, error) {
	response, err := a.apiClient.RetrieveAnalyticsRulesWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return *response.JSON200.Rules, nil
}
