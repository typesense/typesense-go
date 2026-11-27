package typesense

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type AnalyticsRulesInterface interface {
	// Create an analytics rule.
	//
	// Create one analytics rule. When the server rejects it, for example because a rule with that name exists, Create returns *HTTPError. Use CreateMany to create several rules in one request.
	//
	// HTTP: POST /analytics/rules
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	//
	//docgen:keep
	Create(ctx context.Context, rule *api.AnalyticsRuleCreate) (*api.AnalyticsRule, error)
	// Retrieve analytics rules.
	//
	// Retrieve all analytics rules, or only those tagged params.RuleTag when it is set. A nil params retrieves all rules.
	//
	// HTTP: GET /analytics/rules
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	//
	//docgen:keep
	Retrieve(ctx context.Context, params *api.RetrieveAnalyticsRulesParams) ([]*api.AnalyticsRule, error)
	// Create analytics rules in one request.
	//
	// The server answers 200 even when some rules fail. Responses come back in request order, and a rejected rule has Error set. An empty rules slice sends no request.
	//
	// HTTP: POST /analytics/rules
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	//
	//docgen:keep
	CreateMany(ctx context.Context, rules []*api.AnalyticsRuleCreate) ([]*api.AnalyticsRuleCreateResponse, error)
	// Create or update an analytics rule.
	//
	// Create the rule with the given name, or update it if it exists. An update merges into the stored rule and needs the same type, collection and event type. The name in the path wins over rule.Name.
	//
	// HTTP: PUT /analytics/rules/{ruleName}
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	//
	//docgen:keep
	Upsert(ctx context.Context, ruleName string, rule *api.AnalyticsRuleCreate) (*api.AnalyticsRule, error)
}

type analyticsRules struct {
	apiClient APIClientInterface
}

// Create an analytics rule.
//
// Create one analytics rule. When the server rejects it, for example because a rule with that name exists, Create returns *HTTPError. Use CreateMany to create several rules in one request.
//
// HTTP: POST /analytics/rules
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
//
//docgen:keep
func (a *analyticsRules) Create(ctx context.Context, rule *api.AnalyticsRuleCreate) (*api.AnalyticsRule, error) {
	jsonData, err := json.Marshal(rule)
	if err != nil {
		return nil, err
	}

	// The generated parser can't decode the response into its union type.
	httpResp, err := a.apiClient.CreateAnalyticsRuleWithBody(ctx, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, &HTTPError{Status: httpResp.StatusCode, Body: responseBody}
	}

	var created api.AnalyticsRule
	if err := json.Unmarshal(responseBody, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Retrieve analytics rules.
//
// Retrieve all analytics rules, or only those tagged params.RuleTag when it is set. A nil params retrieves all rules.
//
// HTTP: GET /analytics/rules
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
//
//docgen:keep
func (a *analyticsRules) Retrieve(ctx context.Context, params *api.RetrieveAnalyticsRulesParams) ([]*api.AnalyticsRule, error) {
	response, err := a.apiClient.RetrieveAnalyticsRulesWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}

	rules := make([]*api.AnalyticsRule, len(*response.JSON200))
	for i, rule := range *response.JSON200 {
		rules[i] = &rule
	}
	return rules, nil
}

// Create analytics rules in one request.
//
// The server answers 200 even when some rules fail. Responses come back in request order, and a rejected rule has Error set. An empty rules slice sends no request.
//
// HTTP: POST /analytics/rules
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
//
//docgen:keep
func (a *analyticsRules) CreateMany(ctx context.Context, rules []*api.AnalyticsRuleCreate) ([]*api.AnalyticsRuleCreateResponse, error) {
	if len(rules) == 0 {
		return []*api.AnalyticsRuleCreateResponse{}, nil
	}

	jsonData, err := json.Marshal(rules)
	if err != nil {
		return nil, err
	}

	// The generated parser can't decode an array body into its union type.
	httpResp, err := a.apiClient.CreateAnalyticsRuleWithBody(ctx, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, &HTTPError{Status: httpResp.StatusCode, Body: responseBody}
	}

	var responses []*api.AnalyticsRuleCreateResponse
	if err := json.Unmarshal(responseBody, &responses); err != nil {
		return nil, err
	}
	return responses, nil
}

// Create or update an analytics rule.
//
// Create the rule with the given name, or update it if it exists. An update merges into the stored rule and needs the same type, collection and event type. The name in the path wins over rule.Name.
//
// HTTP: PUT /analytics/rules/{ruleName}
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
//
//docgen:keep
func (a *analyticsRules) Upsert(ctx context.Context, ruleName string, rule *api.AnalyticsRuleCreate) (*api.AnalyticsRule, error) {
	jsonData, err := json.Marshal(rule)
	if err != nil {
		return nil, err
	}
	// The spec types this body as AnalyticsRuleUpdate, which can't create a rule.
	response, err := a.apiClient.UpsertAnalyticsRuleWithBodyWithResponse(ctx, ruleName, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}
	return response.JSON200, nil
}
