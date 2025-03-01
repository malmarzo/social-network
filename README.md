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
- The name of the page for the route has to be named as `page.js or page.jsx`

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

- Includes the server.go file to be ran

#### backend/pkg

- Includes the different go packages created (api, db, middleware,...)

#### database & migration

- install the go-migrate tool from the terminal
- to create a new migration file run the command "migrate create -ext sql -dir pkg/db/migrations/sqlite -seq create_notifications_table" by changing the create_notifications_table you change the name of the file
- after creating a migration file there will be one file called down which has the drop table code and another file
  that is called up which has the table itself
- by running the application the migrations will be applied automatically and if the migrations are succesful the
  function will print a message in the terminal

#### CORS Middleware

- Has to be used as a wrapper over the routes handlers to allow requests from the frontend

```go
// Usage in server.go
	http.HandleFunc("/signup", middleware.CorsMiddleware(api.SignupHandler))
```

#### Auth Middleware

- Has to be used over the routes handlers to protect the routes

```go
//Usage in server.go
http.HandleFunc("/logout", middleware.CorsMiddleware(middleware.AuthMiddleware(api.LogoutHandler)))
```

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

### Frontend Route Protection

- Routes are protected using an auth middeware function that validates if the user is logged in or not and redirects him to the correct page found in `frontend/middleware.js`
- Login/Signup are only for non-auth users

### Accessing user Auth state in the frontend in other components

Use the `useAuth` hook to access authentication state in any component:

```javascript
import { useAuth } from "@/context/AuthContext";

const MyComponent = () => {
  const { isLoggedIn, loading, setIsLoggedIn } = useAuth();

  return (
    <div>
      {!isLoggedIn && !loading && <LoginButton />}
      {isLoggedIn && !loading && <UserDashboard />}
    </div>
  );
};
```

### Accessing the websocket connection in other components

Use the `useWebSocket` hook to access WebSocket functionality:

```javascript
import { useWebSocket } from "@/context/Websocket";

const ChatComponent = () => {
  const { addMessageHandler, sendMessage } = useWebSocket();
};
```

### Receiving and handling different types of websocket msgs in other components

- Example in the `UserNotifier` component

```javascript
import { useEffect } from "react";
import { useWebSocket } from "@/context/Websocket";

//Used this component in the layout.js file to notify users
const UserNotifier = () => {
  const { addMessageHandler } = useWebSocket();

  useEffect(() => {
    //Adding msg Handlers (set the msg type and the function to handle it)

    addMessageHandler("newUser", (msg) => {
      alert("New user joined");
    });

    addMessageHandler("removeUser", (msg) => {
      alert("User left");
    });

    addMessageHandler("hello", (msg) => {
      alert(msg.content);
    });
  }, [addMessageHandler]);

  return null; 
};

export default UserNotifier;
```

### Sending websocket msgs in other components

- Example in `HelloSender` component

```javascript
import React from "react";
import { useWebSocket } from "@/context/Websocket";

const HelloSender = () => {
  const { sendMessage } = useWebSocket();

  const sendHello = () => {
    const msg = {};
    msg.type = "hello";
    msg.content = "Hello World!";
    sendMessage(msg);
  };

  return (
    <div>
      <button
        onClick={sendHello}
      >
        Send Hello
      </button>
    </div>
  );
};

export default HelloSender;
```
