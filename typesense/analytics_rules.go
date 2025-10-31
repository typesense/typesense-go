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
	Create(ctx context.Context, rules []*api.AnalyticsRuleCreate) ([]*api.AnalyticsRule, error)
	Retrieve(ctx context.Context) ([]*api.AnalyticsRule, error)
}

type analyticsRules struct {
	apiClient APIClientInterface
}

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
