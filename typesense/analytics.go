package typesense

type AnalyticsInterface interface {
	// Events manages analytics events, such as clicks, conversions, and visits.
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	Events() AnalyticsEventsInterface
	// Rules manages analytics rules, which define how search queries and user events are aggregated.
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	Rules() AnalyticsRulesInterface
	// Rule returns a handle for the analytics rule with the given name.
	//
	// See: https://typesense.org/docs/latest/api/analytics-query-suggestions.html
	Rule(ruleName string) AnalyticsRuleInterface
}

type analytics struct {
	apiClient APIClientInterface
}

func (a *analytics) Events() AnalyticsEventsInterface {
	return &analyticsEvents{apiClient: a.apiClient}
}

func (a *analytics) Rules() AnalyticsRulesInterface {
	return &analyticsRules{apiClient: a.apiClient}
}

func (a *analytics) Rule(ruleName string) AnalyticsRuleInterface {
	return &analyticsRule{apiClient: a.apiClient, ruleName: ruleName}
}
