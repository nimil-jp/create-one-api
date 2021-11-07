package request

import "time"

type WebhookPaypalConnected struct {
	Id           string    `json:"id"`
	EventVersion string    `json:"event_version"`
	CreateTime   time.Time `json:"create_time"`
	ResourceType string    `json:"resource_type"`
	EventType    string    `json:"event_type"`
	Summary      string    `json:"summary"`
	Resource     struct {
		PartnerClientId string `json:"partner_client_id"`
		Links           []struct {
			Method      string `json:"method"`
			Rel         string `json:"rel"`
			Description string `json:"description"`
			Href        string `json:"href"`
		} `json:"links"`
		MerchantId string `json:"merchant_id"`
		TrackingId string `json:"tracking_id"`
	} `json:"resource"`
	Links []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}
