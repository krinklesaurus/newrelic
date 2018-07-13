package newrelic

import (
	"encoding/json"
	"strconv"
)

// AlertCondition describes what triggers an alert for a specific policy.
type AlertCondition struct {
	ID                  int                  `json:"id,omitempty"`
	Type                string               `json:"type,omitempty"`
	Name                string               `json:"name,omitempty"`
	Enabled             bool                 `json:"enabled,omitempty"`
	Entities            []string             `json:"entities,omitempty"`
	Metric              string               `json:"metric,omitempty"`
	GCMetric            string               `json:"gc_metric"`
	ConditionScope      string               `json:"condition_scope"`
	ViolationCloserTime int                  `json:"violation_close_timer,omitempty"`
	RunbookURL          string               `json:"runbook_url,omitempty"`
	Terms               []AlertConditionTerm `json:"terms,omitempty"`
	UserDefined         AlertUserDefined     `json:"user_defined,omitempty"`
}

func (c *AlertCondition) String() string {
	if c == nil {
		return ""
	}
	asJson, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(asJson)
}

// AlertConditionTerm defines thresholds that trigger an AlertCondition.
type AlertConditionTerm struct {
	Duration     string `json:"duration,omitempty"`
	Operator     string `json:"operator,omitempty"`
	Priority     string `json:"priority,omitempty"`
	Threshold    string `json:"threshold,omitempty"`
	TimeFunction string `json:"time_function,omitempty"`
}

// AlertUserDefined describes user-defined behavior for an AlertCondition.
type AlertUserDefined struct {
	Metric        string `json:"metric,omitempty"`
	ValueFunction string `json:"value_function,omitempty"`
}

// AlertConditionOptions define filters for GetAlertConditions.
type AlertConditionOptions struct {
	policyID int
	Page     int
}

func (o *AlertConditionOptions) String() string {
	if o == nil {
		return ""
	}
	return encodeGetParams(map[string]interface{}{
		"policy_id": o.policyID,
		"page":      o.Page,
	})
}

// GetAlertConditions will return any AlertCondition defined for a given
// policy, optionally filtered by AlertConditionOptions.
func (c *Client) GetAlertConditions(policy int, options *AlertConditionOptions) ([]AlertCondition, error) {
	resp := &struct {
		Conditions []AlertCondition `json:"conditions,omitempty"`
	}{}
	options.policyID = policy
	err := c.doGet("alerts_conditions.json", options, resp)
	if err != nil {
		return nil, err
	}
	return resp.Conditions, nil
}

func (c *Client) CreateAlertCondition(policy int, condition *AlertCondition) (*AlertCondition, error) {
	resp := &struct {
		Condition *AlertCondition `json:"condition,omitempty"`
	}{}

	body := struct {
		Condition *AlertCondition `json:"condition,omitempty"`
	}{condition}

	err := c.doPost("alerts_conditions/policies/"+strconv.Itoa(policy)+".json", nil, body, resp)
	if err != nil {
		return nil, err
	}
	return resp.Condition, nil
}

func (c *Client) UpdateAlertCondition(condition *AlertCondition) (*AlertCondition, error) {
	resp := &struct {
		Condition *AlertCondition `json:"condition,omitempty"`
	}{}

	body := struct {
		Condition *AlertCondition `json:"condition,omitempty"`
	}{condition}

	err := c.doUpdate("alerts_conditions/"+strconv.Itoa(condition.ID)+".json", body, resp)
	if err != nil {
		return nil, err
	}
	return resp.Condition, nil
}

func (c *Client) DeleteAlertCondition(condition int) (*AlertCondition, error) {
	resp := &struct {
		AlertCondition *AlertCondition `json:"condition,omitempty"`
	}{}

	err := c.doDelete("alerts_conditions/"+strconv.Itoa(condition)+".json", resp)
	if err != nil {
		return nil, err
	}
	return resp.AlertCondition, nil
}
