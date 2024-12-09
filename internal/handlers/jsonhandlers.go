package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"

	"WST_lab1_client/internal/models"
)

func AddPersonHandler(url string, name string, surname string, age int, email string, telephone string, logger *zap.Logger) {
	person := models.Person{Name: name, Surname: surname, Age: age, Email: email, Telephone: telephone}
	requestJSON, err := json.Marshal(person)
	if err != nil {
		logger.Fatal("Error marshaling request", zap.Error(err))
	}
	body, err := sendRESTRequest(http.MethodPost, url+"/api/v1/persons", requestJSON, logger)
	if err != nil {
		logger.Warn("Error calling AddPerson", zap.Error(err))
		return
	}
	var response models.Person
	if err := json.Unmarshal(body, &response); err != nil {
		handleErrorResponse(body, logger)
		return
	}
	fmt.Printf("Added person with ID: %d\n", response.ID)
}

func GetPersonHandler(url string, id int, logger *zap.Logger) {
	body, err := sendRESTRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/person/%d", url, id), nil, logger)
	if err != nil {
		logger.Warn("Error calling GetPerson", zap.Error(err))
		return
	}
	if body == nil {
		fmt.Println("Not Found: The requested resource could not be found.")
		return
	}
	var response models.Person
	if err := json.Unmarshal(body, &response); err != nil {
		handleErrorResponse(body, logger)
		return
	}
	values := reflect.ValueOf(response)
	types := values.Type()
	fmt.Println("Retrieved person:")
	for i := 0; i < values.NumField(); i++ {
		fmt.Println(types.Field(i).Name+":", values.Field(i))
	}
}

func GetAllPersonsHandler(url string, logger *zap.Logger) {
	body, err := sendRESTRequest(http.MethodGet, url+"/api/v1/persons/list", nil, logger)
	if err != nil {
		logger.Warn("Error calling GetAllPersons", zap.Error(err))
		return
	}
	var persons []models.Person
	if err := json.Unmarshal(body, &persons); err != nil {
		handleErrorResponse(body, logger)
		return
	}
	if len(persons) == 0 {
		fmt.Println("No persons found.")
		return
	}
	for _, person := range persons {
		fmt.Printf("- ID: %d, Name: %s, Surname: %s, Age: %d\n",
			person.ID,
			person.Name,
			person.Surname,
			person.Age,
		)
	}
}
func UpdatePersonHandler(url string, id int, name string, surname string, age int, email string, telephone string, logger *zap.Logger) {
	person := models.Person{ID: id, Name: name, Surname: surname, Age: age, Email: email, Telephone: telephone}
	requestJSON, err := json.Marshal(person)
	if err != nil {
		logger.Fatal("Error marshaling request", zap.Error(err))
		return
	}
	body, err := sendRESTRequest(http.MethodPut, fmt.Sprintf("%s/api/v1/person/%d", url, id), requestJSON, logger)
	if err != nil {
		logger.Warn("Error calling UpdatePerson", zap.Error(err))
		return
	}
	var response struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		handleErrorResponse(body, logger)
		return
	}
	if response.Message == "Person updated successfully!" {
		fmt.Println("Updated person successfully.")
	} else {
		fmt.Printf("Failed to update person: %s\n", response.Message)
	}
}

func DeletePersonHandler(url string, id int, logger *zap.Logger) {
	body, err := sendRESTRequest(http.MethodDelete,
		fmt.Sprintf("%s/api/v1/person/%d", url, id), nil,
		logger)

	if err != nil {
		logger.Warn("Error calling DeletePerson", zap.Error(err))
		return
	}
	var response struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		handleErrorResponse(body, logger)
		return
	}
	if response.Message == "Deleted Successfully" {
		fmt.Println("Deleted person successfully.")
	} else {
		fmt.Printf("Failed to delete person: %s\n", response.Message)
	}
}
func SearchPersonsHandler(url string, query string, logger *zap.Logger) {
	searchURL := fmt.Sprintf("%s/api/v1/persons?query=%s", url, query)
	body, err := sendRESTRequest(http.MethodGet, searchURL, nil, logger)
	if err != nil {
		logger.Warn("Error calling SearchPersons", zap.Error(err))
		return
	}
	var errorResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if json.Unmarshal(body, &errorResponse) == nil {
		if errorResponse.Code == "not_found" {
			fmt.Println("Not Found: The requested resource could not be found.")
			return
		}
		fmt.Printf("Error: %s\n", errorResponse.Message)
		return
	}
	var persons []models.Person

	if err := json.Unmarshal(body, &persons); err != nil {
		logger.Error("Failed to unmarshal response body", zap.Error(err))
		handleErrorResponse(body, logger)
		return
	}
	if len(persons) == 0 {
		fmt.Println("No persons found.")
		return
	}
	fmt.Println("Search Results:")
	for _, person := range persons {
		fmt.Printf("- ID: %d, Name: %s, Surname: %s, Age: %d\n",
			person.ID,
			person.Name,
			person.Surname,
			person.Age,
		)
	}
}

func sendRESTRequest(method string, url string, body []byte, logger *zap.Logger) ([]byte, error) {
	reqBody := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		logger.Error("Error creating request", zap.Error(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error sending request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response body", zap.Error(err))
		return nil, err
	}
	if resp.StatusCode >= 400 {
		handleHTTPError(resp.StatusCode, responseBody)
		return responseBody, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	return responseBody, nil
}

func handleHTTPError(statusCode int, responseBody []byte) {
	switch statusCode {
	case http.StatusBadRequest:
		fmt.Println("Bad Request: The server could not understand the request due to invalid syntax.")
	case http.StatusUnauthorized:
		fmt.Println("Unauthorized: Access is denied due to invalid credentials.")
	case http.StatusForbidden:
		fmt.Println("Forbidden: You do not have permission to access this resource.")
	case http.StatusNotFound:
		fmt.Println("Not Found: The requested resource could not be found.")
	case http.StatusConflict:
		var errorResponse struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(responseBody, &errorResponse) == nil {
			fmt.Printf("Status Code: %d\n", statusCode)
			fmt.Printf("Conflict: %s\n", errorResponse.Message)
		} else {
			fmt.Println("Conflict occurred but no message was returned.")
		}
	case http.StatusInternalServerError:
		fmt.Println("Internal Server Error: The server encountered an unexpected condition.")
	default:
		fmt.Printf("Unexpected error occurred. Status Code: %d\n", statusCode)
	}
}

func handleErrorResponse(body []byte, logger *zap.Logger) {
	var errorResponse struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Details string `json:"details,omitempty"`
	}
	if json.Unmarshal(body, &errorResponse) == nil {
		fmt.Printf("Server Error (%d): %s\nDetails: %s\n", errorResponse.Code, errorResponse.Message, errorResponse.Details)
	} else {
		fmt.Println("An unknown error occurred while processing the server's response.")
	}
}
