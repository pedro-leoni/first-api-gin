package models

type Character struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name" validate:"required"`
	Birthday  string `json:"birthday,omitempty"`
	Dead      bool   `json:"dead" validate:"required"`
	Relevance string `json:"relevance" validate:"required"`
	Seasons   int    `json:"seasons" validate:"required"`
}
