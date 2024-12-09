package main

import (
	"WST_lab1_client/internal/handlers"
	"WST_lab1_client/internal/logger"
	"flag"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	logConfig := logger.NewLoggerConfig()

	log, err := logger.NewLogger(logConfig)
	if err != nil {
		fmt.Println("Failed to initialize logger")
		log.Fatal("Failed to initialize logger")
	}
	defer func(log *zap.Logger) {
		err := log.Sync()
		if err != nil {
			log.Error("Error syncing logger", zap.Error(err))
		}
	}(log)

	url := flag.String("url", "http://localhost:8095", "REST server URL")
	method := flag.String("method", "", "Method to call (addperson|getperson|getallpersons|updateperson|deleteperson|searchperson)")
	name := flag.String("name", "", "Name of the person (required for addperson and updateperson)")
	surname := flag.String("surname", "", "Surname of the person (required for addperson and updateperson)")
	id := flag.Int("id", 0, "ID of the person (required for getperson, updateperson and deleteperson)")
	query := flag.String("query", "", "Query for searching person (required for searchperson)")
	age := flag.Int("age", 0, "Age of the person (required for addperson and updateperson)")
	email := flag.String("email", "", "Email of the person (required for addperson)")
	telephone := flag.String("telephone", "", "Telephone of the person (required for addperson)")

	flag.Parse()

	switch *method {
	case "addperson":
		if *name == "" || *surname == "" || *age <= 0 || *email == "" || *telephone == "" {
			fmt.Println("Name, surname, age (greater than 0), email, and telephone are required for addperson.")
			log.Fatal("Name, surname, age (greater than 0), email, and telephone are required for addperson.")
		}
		handlers.AddPersonHandler(*url, *name, *surname, *age, *email, *telephone, log)
	case "getperson":
		if *id <= 0 {
			fmt.Println("ID must be greater than 0 for getperson.")
			log.Fatal("ID must be greater than 0 for getperson.")
		}
		handlers.GetPersonHandler(*url, *id, log)
	case "getallpersons":
		handlers.GetAllPersonsHandler(*url, log)
	case "updateperson":
		if *id <= 0 || *name == "" || *surname == "" || *age <= 0 {
			fmt.Println("ID must be greater than 0 and both name and surname are required for updateperson. Age must be greater than 0.")
			log.Fatal("ID must be greater than 0 and both name and surname are required for updateperson. Age must be greater than 0.")
		}
		handlers.UpdatePersonHandler(*url, *id, *name, *surname, *age, *email, *telephone, log)
	case "deleteperson":
		if *id <= 0 {
			fmt.Println("ID must be greater than 0 for deleteperson.")
			log.Fatal("ID must be greater than 0 for deleteperson.")
		}
		handlers.DeletePersonHandler(*url, *id, log)
	case "searchperson":
		if *query == "" {
			fmt.Println("Query is required for searchperson.")
			log.Fatal("Query is required for searchperson.")
		}
		handlers.SearchPersonsHandler(*url, *query, log)
	default:
		fmt.Println("Unknown method. Use one of addperson|getperson|getallpersons|updateperson|deleteperson|searchperson.")
		log.Fatal("Unknown method. Use one of addperson|getperson|getallpersons|updateperson|deleteperson|searchperson.")
	}
}
