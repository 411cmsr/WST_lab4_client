# WST_lab4_client
#### To create a request, use the following command line arguments:

-url string 

######     _REST server URL (default "http://localhost:8095")_ 

-method string 

######     _Method to call (addperson|getperson|getallpersons|updateperson|deleteperson|searchperson)_ 

-id int 

######     _ID of the person (required for getperson, updateperson and deleteperson)_ 

-name string 

######     _Name of the person (required for addperson and updateperson)_ 

-age int 

######     _Age of the person (required for addperson and updateperson)_ 

-surname string 

######     _Surname of the person (required for addperson and updateperson)_

-telephone string

######     _Telephone of the person (required for addperson)_

-email string

######     _Email of the person (required for addperson)_

-query string 

###### _Query for searching person (required for searchperson)_


### **example request:** 

#### go run main.go getperson -id  9518 -url http://127.0.0.1:8095

### **example response:**
_Retrieved person:_
_ID: 9518_
_Name: Ольга_
_Surname: Беригальц_
_Age: 34_
_Email: olga@mail.com_
_Telephone: +70011234576_

