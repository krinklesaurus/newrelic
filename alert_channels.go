package newrelic

import (
	"fmt"
	"strconv"
)

type AlertChannelLinks struct {
	PolicyIDs []int `json:"policy_ids,omitempty"`
}

type AlertChannel struct {
	ID            int                    `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	Links         AlertChannelLinks      `json:"links,omitempty"`
}

type AlertChannelOptions struct {
	Page      int
	policyIds []int
}

func (o *AlertChannelOptions) String() string {
	if o == nil {
		return ""
	}

	s := ""
	if o.Page != 0 {
		s += fmt.Sprintf("page=%s", o.Page)
	}
	for i, pID := range o.policyIds {
		if i == 0 {
			s += "policy_ids[]="
		}
		s += fmt.Sprintf("%d,", pID)
	}

	return s
}

func (c *Client) GetAlertChannels(options *AlertChannelOptions) ([]AlertChannel, error) {
	resp := &struct {
		AlertChannels []AlertChannel `json:"channels,omitempty"`
	}{}
	err := c.doGet("alerts_channels.json", options, resp)
	if err != nil {
		return nil, err
	}
	return resp.AlertChannels, nil
}

func (c *Client) CreateAlertChannel(channel *AlertChannel, policyIDs []int) ([]AlertChannel, error) {
	resp := &struct {
		AlertChannels []AlertChannel `json:"channels,omitempty"`
	}{}

	options := &AlertChannelOptions{
		policyIds: policyIDs,
	}

	body := struct {
		Channel *AlertChannel `json:"channel,omitempty"`
	}{Channel: channel}

	err := c.doPost("alerts_channels.json", options, body, resp)
	if err != nil {
		return nil, err
	}
	return resp.AlertChannels, nil
}

func (c *Client) DeleteAlertChannel(channel int) (*AlertChannel, error) {
	resp := &struct {
		AlertChannel *AlertChannel `json:"channel,omitempty"`
	}{}

	err := c.doDelete("alerts_channels/"+strconv.Itoa(channel)+".json", resp)
	if err != nil {
		return nil, err
	}
	return resp.AlertChannel, nil
}
