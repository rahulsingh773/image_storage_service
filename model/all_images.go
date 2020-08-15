package model

type AllImagesList struct {
	Image string `json:"image" validate:"nonzero"`
	URL   string `json:"url" validate:"nonzero"`
}
