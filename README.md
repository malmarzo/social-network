# social-network

A Facebook like social media application made using Next.js and go.

### Prerequisites

- Install node if not installed
- Install pnpm:
  `npm install -g pnpm`
- Clone the repo: `git clone https://learn.reboot01.com/git/malmarzo/social-network.git`

### Running the frontend

- Navigate to the frontend dir: `cd frontend`
- Run: `pnpm dev`
- Application will start on port `3000`

### Running the backend

- Navigate to the backend dir: `cd backend`
- Run: `go run cmd/server.go`

### Project Structure

```
social-network/
│── frontend/
│   ├── app/
│   ├── public/
│   ├── styles/
│   ├── utils/
│
└── backend/
    ├── cmd/
    │   └── main.go
    │
    ├── pkg/
        ├── api/
        ├── dataModels/
        ├── db/
        │   ├── migrations/
        │   ├── queries/
        │   └── sqlite/
        ├── middleware/
        ├── utils/
        └── websocket/
```

### frontend

- Will contain everything related to the frontend application (Next.JS APP)

#### frontend/app

- Will hold the different routes as directories that include the different pages of the app
- The test directory inside it is an example of how to create a new route
- Navigating to localhost:3000/test would show the newly created page
- The name of the page for the route has to be named as `page.js or page.jsx`

#### frontend/public

- Will include all the public assests like images, etc..

#### frontend/styles

- CSS styles for the pages

#### frontend/utils

- Utility code

### Frontend Utilities

#### invokeAPI Function

```javascript
// Usage: invokeAPI(route, body, method, contentType)
// Example:
const response = await invokeAPI(
  "signup",
  formData,
  "POST",
  "multipart/form-data"
);
```

- `route`: API endpoint path (string)
- `body`: Request body (object or FormData)
- `method`: HTTP method ("GET", "POST", etc.)
- `contentType`: Content type header (optional, defaults to "application/json")
- Returns: Promise with API response

#### Alert Components

```jsx
// Success Alert
<SuccessAlert
  msg="Operation successful!"
  link="/home"
  linkText="Go to Home"
/>

// Fail Alert
<FailAlert msg="Operation failed!" />
```

- `msg`: Message to display

### backend

- Will hold all the backend components: Server, App Logic, Database

#### backend/cmd

- Includes the main.go file to be ran

#### backend/pkg

- Includes the different go packages created (api, db, middleware,...)

#### pkg/api

- Will include all the different API handlers

#### pkg/dataModels

- Will include all the different structs to be used and constant variables

#### pkg/middleware

- Will include all the middlewares (auth, cors,...)

#### pkg/utils

- Utility code

#### pkg/websocket

- Will include all the websocket handling

#### pkg/db

- Everything related to db interation
- Migrations in the `/migrations` dir
- DB connection in the `/sqlite` dir in the `sqlite.go` file
- Functions and queries in the `queries` dir

#### database & migration

- to create a new migration file run the command "migrate create -ext sql -dir pkg/db/migrations/sqlite -seq create_notifications_table" by changing the create_notifications_table you change the name of the file
- after creating a migration file there will be one file called down which has the drop table code and another file
  that is called up which has the table itself
- by running the application the migrations will be applied automatically and if the migrations are succesful the
  function will print a message in the terminal

### Backend Utilities

#### CORS Middleware

```go
// Usage in server.go
	http.HandleFunc("/signup", middleware.CorsMiddleware(api.SignupHandler))
```

Enables Cross-Origin Resource Sharing with configured options:

- Allowed Origins: localhost:3000
- Allowed Methods: GET, POST, PUT, DELETE
- Allowed Headers: Content-Type
- Credentials: Allowed

#### SendResponse Function

```go
// Usage in handlers
utils.SendResponse(w, datamodels.Response{
    Code: http.StatusOK,
    Status: "Success",
    ErrorMsg: "",
})
```

Sends a standardized JSON response with:

- `Code`: HTTP status code
- `Status`: "Success" or "Failed"
- `ErrorMsg`: Error message (optional)

#### Password Utilities

```go
// Hash a password
hashedPassword, err := utils.HashPassword(password)

// Verify a password
match := utils.CheckPasswordHash(password, hashedPassword)
```

- `HashPassword`: Generates a secure hash using bcrypt
- `CheckPasswordHash`: Verifies a password against its hash
- Uses bcrypt with default cost factor
