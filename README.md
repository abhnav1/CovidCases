# CovidCasesInYourState
 An API which gives the total number of covid cases in your state , taking your GPS co ordinates (latitude and longitude) as Input. Only valid for Indian locations. 

 ## How to use the API?

 1. Run the Final.go file using the command prompt (command: go run Final.go)
 2. Enter your longitude and latitude in the url, considering the following format  

           localhost:8080/getCases?lat=<Your_latitude_value>&lon=<Your_longitude_value>  
           
    Example:
    	   localhost:8080/getCases?lat=26.8467&lon=80.9462

 The webpage will show you the state which contains these co ordinates and the total number of covid-19 cases in your state


 ## Tech Used 

 * Golang using visual studio
 * Echo framework 
 * MongoDB database using Atlas
