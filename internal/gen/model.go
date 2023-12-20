package gen

import "encoding/json"

type AsyncAPIModel struct {
	Version    string                 `json:"asyncapi"`
	Servers    map[string]Server      `json:"servers,omitempty"`
	Channels   map[string]ChannelItem `json:"channels"`
	Components *Components            `json:"components,omitempty"`
}

type Server struct {
	Url      string `json:"url"`
	Protocol string `json:"protocol"`
}

type ChannelItem struct {
	Publish   *Operation `json:"publish,omitempty"`
	Subscribe *Operation `json:"subscribe,omitempty"`
}

type Operation struct {
	Security []map[string][]string `json:"security,omitempty"`
	ID       string                `json:"operationId,omitempty"`
	Message  *Message              `json:"message,omitempty"`
}

type Message struct {
	Reference     *Reference     `json:"-"`
	MessageEntity *MessageEntity `json:"-,omitempty"`
}

type Reference struct {
	Ref string `json:"$ref"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &m.MessageEntity)
	if err == nil && m.MessageEntity.Payload != nil {
		return nil
	}
	m.MessageEntity = nil

	if err := json.Unmarshal(data, &m.Reference); err != nil {
		return err
	}
	return nil
}

type MessageEntity struct {
	SchemaFormat string                 `json:"schemaFormat,omitempty"`
	ContentType  string                 `json:"contentType,omitempty"`
	MessageID    string                 `json:"messageId,omitempty"`
	Payload      map[string]interface{} `json:"payload,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Title        string                 `json:"title,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Deprecated   bool                   `json:"deprecated,omitempty"`
}

type Components struct {
	Schemas  map[string]map[string]interface{} `json:"schemas,omitempty"`
	Messages map[string]Message                `json:"messages,omitempty"`
}
