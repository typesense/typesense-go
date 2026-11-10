package typesense

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/typesense/typesense-go/v4/typesense/api"
)

type AnalyticsRulesInterface interface {
	// Create analytics rule(s).
	//
	// Create one or more analytics rules. You can send a single rule object or an array of rule objects.
	//
	// HTTP: POST /analytics/rules
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	Create(ctx context.Context, rules []*api.AnalyticsRuleCreate) ([]*api.AnalyticsRule, error)
	// Retrieve analytics rules.
	//
	// Retrieve all analytics rules. Use the optional rule_tag filter to narrow down results.
	//
	// HTTP: GET /analytics/rules
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	Retrieve(ctx context.Context) ([]*api.AnalyticsRule, error)
	// Retrieve analytics rules.
	//
	// Retrieve all analytics rules. Use the optional rule_tag filter to narrow down results.
	//
	// HTTP: GET /analytics/rules
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	RetrieveWithParams(ctx context.Context, params *api.RetrieveAnalyticsRulesParams) ([]*api.AnalyticsRule, error)
	// Create analytics rules in one request.
	//
	// The server answers 200 even when some rules fail. Responses come back in request order, and a rejected rule has Error set.
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

// Create analytics rule(s).
//
// Create one or more analytics rules. You can send a single rule object or an array of rule objects.
//
// HTTP: POST /analytics/rules
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
func (a *analyticsRules) Create(ctx context.Context, rules []*api.AnalyticsRuleCreate) ([]*api.AnalyticsRule, error) {
	// Convert []*AnalyticsRuleCreate to []AnalyticsRuleCreate for the API call
	ruleCreates := make([]api.AnalyticsRuleCreate, len(rules))
	for i, rule := range rules {
		ruleCreates[i] = *rule
	}

	// Encode the rules as JSON
	jsonData, err := json.Marshal(ruleCreates)
	if err != nil {
		return nil, err
	}

	// Use the lower-level API to get the raw response
	httpResp, err := a.apiClient.CreateAnalyticsRuleWithBody(ctx, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// Read the raw response body
	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, &HTTPError{Status: httpResp.StatusCode, Body: responseBody}
	}

	// Parse the response manually since the generated union type has issues
	// First try as an array (most common case)
	var rulesArray []api.AnalyticsRule
	if err := json.Unmarshal(responseBody, &rulesArray); err == nil {
		// Successfully parsed as array
		result := make([]*api.AnalyticsRule, len(rulesArray))
		for i, rule := range rulesArray {
			ruleCopy := rule
			result[i] = &ruleCopy
		}
		return result, nil
	}

	// Try as single rule
	var singleRule api.AnalyticsRule
	if err := json.Unmarshal(responseBody, &singleRule); err == nil {
		// Successfully parsed as single rule
		return []*api.AnalyticsRule{&singleRule}, nil
	}

	// If we can't parse either way, return error
	return nil, fmt.Errorf("failed to parse response: %s", string(responseBody))
}

// Retrieve analytics rules.
//
// Retrieve all analytics rules. Use the optional rule_tag filter to narrow down results.
//
// HTTP: GET /analytics/rules
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
func (a *analyticsRules) Retrieve(ctx context.Context) ([]*api.AnalyticsRule, error) {
	response, err := a.apiClient.RetrieveAnalyticsRulesWithResponse(ctx, &api.RetrieveAnalyticsRulesParams{})
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, &HTTPError{Status: response.StatusCode(), Body: response.Body}
	}

	// Convert []AnalyticsRule to []*AnalyticsRule
	rules := make([]*api.AnalyticsRule, len(*response.JSON200))
	for i, rule := range *response.JSON200 {
		ruleCopy := rule
		rules[i] = &ruleCopy
	}
	return rules, nil
}

// Retrieve analytics rules.
//
// Retrieve all analytics rules. Use the optional rule_tag filter to narrow down results.
//
// HTTP: GET /analytics/rules
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
func (a *analyticsRules) RetrieveWithParams(ctx context.Context, params *api.RetrieveAnalyticsRulesParams) ([]*api.AnalyticsRule, error) {
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
// The server answers 200 even when some rules fail. Responses come back in request order, and a rejected rule has Error set.
//
// HTTP: POST /analytics/rules
//
// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
//
//docgen:keep
func (a *analyticsRules) CreateMany(ctx context.Context, rules []*api.AnalyticsRuleCreate) ([]*api.AnalyticsRuleCreateResponse, error) {
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
