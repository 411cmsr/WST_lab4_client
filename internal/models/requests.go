package models

import (
	
)

type AddPersonRequest struct {
	Person  Person   `json:"person"`
	
}

type UpdatePersonRequest struct {
    Person Person `json:"person"`
}
type SearchPersonRequest struct {
    Query string `json:"query"`
}
