package typesense

type AnalyticsInterface interface {
	Events() AnalyticsEventsInterface
	Rules() AnalyticsRulesInterface
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
