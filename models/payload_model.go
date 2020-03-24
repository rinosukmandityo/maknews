package models

type GetPayload struct {
	Filter map[string]interface{} `json:"filter" bson:"filter" msgpack:"filter"`
	Offset int                    `json:"offset" bson:"offset" msgpack:"offset"`
	Limit  int                    `json:"limit" bson:"limit" msgpack:"limit"`
	Order  map[string]bool        `json:"order" bson:"order" msgpack:"order"`
}

func (m *GetPayload) String() string {
	return "payload"
}
