package models

type GetPayload struct {
	Filter map[string]interface{} `json:"filter" bson:"filter" msgpack:"filter"`
	Offset int                    `json:"offset" bson:"offset" msgpack:"offset"`
	Limit  int                    `json:"limit" bson:"limit" msgpack:"limit"`
}

func (m *GetPayload) String() string {
	return "payload"
}
