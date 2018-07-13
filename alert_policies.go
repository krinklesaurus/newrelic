package newrelic

import "strconv"

type AlertPolicy struct {
	ID                 int    `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	IncidentPreference string `json:"incident_preference,omitempty"`
	CreatedAt          int    `json:"created_at,omitempty"`
	UpdatedAt          int    `json:"updated_at,omitempty"`
}

type AlertPolicyFilter struct {
	Name string
}

type AlertPolicyOptions struct {
	Filter AlertPolicyFilter
	Page   int
}

func (o *AlertPolicyOptions) String() string {
	if o == nil {
		return ""
	}
	return encodeGetParams(map[string]interface{}{
		"filter[name]": o.Filter.Name,
		"page":         o.Page,
	})
}

func (c *Client) GetAlertPolicies(options *AlertPolicyOptions) ([]AlertPolicy, error) {
	resp := &struct {
		AlertPolicies []AlertPolicy `json:"policies,omitempty"`
	}{}
	err := c.doGet("alerts_policies.json", options, resp)
	if err != nil {
		return nil, err
	}
	return resp.AlertPolicies, nil
}

func (c *Client) CreateAlertPolicy(policy *AlertPolicy) (*AlertPolicy, error) {
	resp := &struct {
		Policy *AlertPolicy `json:"policy,omitempty"`
	}{}

	body := struct {
		Policy *AlertPolicy `json:"policy,omitempty"`
	}{policy}

	err := c.doPost("alerts_policies.json", nil, body, resp)
	if err != nil {
		return nil, err
	}
	return resp.Policy, nil
}

func (c *Client) UpdateAlertPolicy(policy *AlertPolicy) (*AlertPolicy, error) {
	resp := &struct {
		Policy *AlertPolicy `json:"policy,omitempty"`
	}{}

	body := struct {
		Policy *AlertPolicy `json:"policy,omitempty"`
	}{policy}

	err := c.doUpdate("alerts_policies/"+strconv.Itoa(policy.ID)+".json", body, resp)
	if err != nil {
		return nil, err
	}
	return resp.Policy, nil
}

func (c *Client) DeleteAlertPolicy(policy int) (*AlertPolicy, error) {
	resp := &struct {
		AlertPolicy *AlertPolicy `json:"policy,omitempty"`
	}{}

	err := c.doDelete("alerts_policies/"+strconv.Itoa(policy)+".json", resp)
	if err != nil {
		return nil, err
	}
	return resp.AlertPolicy, nil
}
