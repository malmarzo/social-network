# social-network

A Facebook like social media application made using Next.js and go.

### Prerequisites

- Install node if not installed
- Install pnpm:
  `npm install -g pnpm`
- Clone the repo: `git clone https://learn.reboot01.com/git/malmarzo/social-network.git`

### Running the frontend

- Navigate to the frontend dir: `cd frontend`
- Run: `pnpm dev` or `npm run dev`
- Application will start on port `3000`

### Running the backend

- Navigate to the backend dir: `cd backend`
- Run: `go run cmd/server.go`

### Running on docker
- In the root directory `/social-network`
- Run the command `docker compose up`
- Containers will start on port 3000 (frontend) and 8080 (backend)




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

