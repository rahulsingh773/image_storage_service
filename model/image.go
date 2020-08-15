package model

type KafkaEvent struct {
	Event     string `json:"event" validate:"nonzero"`
	Album     string `json:"album" validate:"nonzero"`
	Name      string `json:"name" validate:"nonzero"`
	Operation string `json:"operation" validate:"nonzero"`
}
