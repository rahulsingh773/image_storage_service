package model

type CreateAlbum struct {
	Name string `json:"name" validate:"nonzero"`
}
