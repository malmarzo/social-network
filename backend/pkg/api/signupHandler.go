package api

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	datamodels "social-network/pkg/dataModels"
	"social-network/pkg/db/queries"
	"social-network/pkg/utils"
	"strings"
)

// Handles the signup request
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	//Response to be sent back
	response := datamodels.Response{}

	//User to be created
	user := datamodels.User{}

	//Allow only POST method
	if r.Method == http.MethodPost {
		//Recieve and parse the form data
		err := r.ParseMultipartForm(10 << 20) // Limit file size to 10MB
		if err != nil {
			log.Println(err)
			utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Failed to parse form"})
			return
		}

		//Create a new uuid for the user
		id := utils.GenerateUUID()
		if id == "" {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
			return
		}

		// Extract form values
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		dob := r.FormValue("dob")
		nickname := r.FormValue("nickname")
		aboutMe := r.FormValue("about_me")

		//Check if the user is already registered
		// Check if the user is already registered
		errMsg, exists, err := queries.CheckUserExists(email, strings.ToLower(nickname))
		if err != nil {
			log.Println(err)
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
			return
		}
		if exists {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: errMsg})
			return
		}

		// Validate input
		errMsg, valid := utils.SignupValidator(firstName, lastName, email, password, dob, nickname)
		if !valid {
			utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "Invalid input: " + errMsg})
			return
		}

		// Hash the password
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			log.Println(err)
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
			return
		}

		user.ID = id
		user.FirstName = strings.ToUpper(strings.ToLower(firstName))[0:1] + strings.ToLower(firstName)[1:]
		user.LastName = strings.ToUpper(strings.ToLower(lastName))[0:1] + strings.ToLower(lastName)[1:]
		user.Email = email
		user.Password = hashedPassword
		user.DOB = dob
		user.Nickname = strings.ToLower(nickname)
		user.About = aboutMe
		user.Avatar = ""
		user.Private = false

		// Handle file upload
		file, handler, err := r.FormFile("avatar")
		var filePath string
		if err == nil { // Only process if an image is uploaded
			defer file.Close()

			// Validate file type
			allowedExtensions := map[string]bool{
				".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
			}
			ext := filepath.Ext(handler.Filename)
			if !allowedExtensions[ext] {
				utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
				return
			}

			// Create uploads directory if not exists
			uploadDir := "./pkg/db/uploads"
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				err = os.MkdirAll(uploadDir, os.ModePerm)
				if err != nil {
					log.Println(err)
					utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
					return
				}
			}

			// Generate unique filename with UUID
			//Files name will be (avatar_UUID.png) or (avatar_UUID.jpg) or (avatar_UUID.jpeg) or (avatar_UUID.gif)
			filename := "avatar_" + id + ext
			filePath = filepath.Join(uploadDir, filename)
			outFile, err := os.Create(filePath)
			if err != nil {
				log.Println(err)
				utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
				return
			}
			defer outFile.Close()

			// Copy file contents
			_, err = io.Copy(outFile, file)
			if err != nil {
				log.Println(err)
				utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
				return
			}

			// Store the relative path in user.Avatar
			user.Avatar = "uploads/" + filename
		}

		err = queries.AddUser(user)
		if err != nil {
			log.Println(err)
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
			return
		}

		response.Code = 200
		response.Status = "OK"

		utils.SendResponse(w, response) //send the response
	} else {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "Method Not Allowed"})
		return
	}
}
