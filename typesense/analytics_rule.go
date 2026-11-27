package typesense

import (
	"context"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type AnalyticsRuleInterface interface {
	// Delete an analytics rule.
	//
	// Permanently deletes an analytics rule, given it's name.
	//
	// HTTP: DELETE /analytics/rules/{ruleName}
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	Delete(ctx context.Context) (*api.AnalyticsRule, error)
	// Retrieves an analytics rule.
	//
	// Retrieve the details of an analytics rule, given it's name.
	//
	// HTTP: GET /analytics/rules/{ruleName}
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	Retrieve(ctx context.Context) (*api.AnalyticsRule, error)
	// Update an analytics rule.
	//
	// Merge ruleSchema into the stored rule with this name and return the merged rule. Fields that ruleSchema leaves nil keep their stored values. When no rule has this name, the server rejects the update with status 400 and Update returns *HTTPError. Use Rules().Upsert to create a rule.
	//
	// HTTP: PUT /analytics/rules/{ruleName}
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	//
	//docgen:keep
	Update(ctx context.Context, ruleSchema *api.AnalyticsRuleUpdate) (*api.AnalyticsRule, error)
}

type analyticsRule struct {
	apiClient APIClientInterface
	ruleName  string
}

// Delete an analytics rule.
//
// Permanently deletes an analytics rule, given it's name.
//
// HTTP: DELETE /analytics/rules/{ruleName}
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
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

// Retrieves an analytics rule.
//
// Retrieve the details of an analytics rule, given it's name.
//
// HTTP: GET /analytics/rules/{ruleName}
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
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

// Update an analytics rule.
//
// Merge ruleSchema into the stored rule with this name and return the merged rule. Fields that ruleSchema leaves nil keep their stored values. When no rule has this name, the server rejects the update with status 400 and Update returns *HTTPError. Use Rules().Upsert to create a rule.
//
// HTTP: PUT /analytics/rules/{ruleName}
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
//
//docgen:keep
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
