#!/bin/bash

# Define image and container names
FRONTEND_IMAGE="social-frontend"
BACKEND_IMAGE="social-backend"
FRONTEND_CONTAINER="social-frontend"
BACKEND_CONTAINER="social-backend"

# Stop and remove any running containers
docker rm -f $FRONTEND_CONTAINER $BACKEND_CONTAINER 2>/dev/null || true

# Build the frontend image if it doesn't exist
if [[ "$(docker images -q $FRONTEND_IMAGE 2> /dev/null)" == "" ]]; then
    echo "🔧 Building frontend image..."
    docker build -f Dockerfile.frontend -t $FRONTEND_IMAGE .
else
    echo "✅ Frontend image already exists."
fi

# Build the backend image if it doesn't exist
if [[ "$(docker images -q $BACKEND_IMAGE 2> /dev/null)" == "" ]]; then
    echo "🔧 Building backend image..."
    docker build -f Dockerfile.backend -t $BACKEND_IMAGE .
else
    echo "✅ Backend image already exists."
fi

# Start backend container
echo "🚀 Starting backend container..."
docker run -d --name $BACKEND_CONTAINER -p 8080:8080 $BACKEND_IMAGE

# Start frontend container
echo "🚀 Starting frontend container..."
docker run -d --name $FRONTEND_CONTAINER -p 3000:3000 $FRONTEND_IMAGE

echo "✅ Both containers are running!"
echo "🌐 Frontend: http://localhost:3000"
echo "🔌 Backend:  http://localhost:8080"
