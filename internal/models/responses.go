package models

type AddPersonResponse struct {
	Content Person `json:"content"`
}

type UpdatePersonResponse struct {
	Success bool `json:"success"`
}

type DeletePersonResponse struct {
	Success bool `json:"success"`
}

type GetPersonResponse struct {
    Content Person `json:"content"`
}

type GetAllPersonsResponse struct {
    Content struct {
        Persons []Person `json:"persons"`
    } `json:"content"`
}
type SearchPersonsResponse struct {
    Content struct {
        Persons []Person `json:"persons"`
    } `json:"content"`
}