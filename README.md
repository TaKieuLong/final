# CI/CD Project with Jenkins, Docker, PostgreSQL, Golang, and React

This README provides step-by-step instructions to set up and deploy a CI/CD pipeline using Jenkins, Docker, PostgreSQL, Golang (backend), and React (frontend).

---

## **1. Project Overview**

This project demonstrates a CI/CD pipeline that:

1. Builds Docker images for the backend (Golang) and frontend (React).
2. Pushes these images to Docker Hub.
3. Deploys the application using Docker Compose with PostgreSQL as the database.
4. Sends build notifications via Telegram.

---

## **2. Prerequisites**

### **2.1 Tools and Technologies**
- **Jenkins**: CI/CD server
- **Docker**: For containerization
- **Docker Compose**: For managing multiple containers
- **Git**: Version control
- **PostgreSQL**: Database
- **Golang**: Backend API
- **React**: Frontend application

### **2.2 Accounts and Setup**
- **Docker Hub**: Create an account to push Docker images.
- **Telegram Bot**:
  1. Use [@BotFather](https://t.me/BotFather) on Telegram to create a bot and get your bot token.
  2. Find your chat ID using tools like [getIDBot](https://t.me/getidbot).

---

## **3. Repository Structure**

```plaintext
project/
├── backend/                 # Golang backend
│   ├── Dockerfile
│   ├── main.go
│   └── ...
├── frontend/                # React frontend
│   ├── Dockerfile
│   ├── src/
│   └── ...
├── docker-compose.yml       # Docker Compose file
├── Jenkinsfile              # Jenkins pipeline script
└── README.md                # This documentation
```

---

## **4. Step-by-Step Instructions**

### **4.1 Clone the Repository**
1. Create a GitHub repository (e.g., `ci-cd-demo`).
2. Clone the repository locally:
   ```bash
   git clone https://github.com/<your_username>/ci-cd-demo.git
   cd ci-cd-demo
   ```

### **4.2 Backend Setup**
1. Navigate to the `backend` directory.
2. Initialize the Go module:
   ```bash
   cd backend
   go mod init <module_name>
   ```
3. Add the `Gin Gonic` dependency:
   ```bash
   go get -u github.com/gin-gonic/gin
   ```
4. Add PostgreSQL dependency:
   ```bash
   go get github.com/lib/pq
   ```
5. Create a `Dockerfile` in the `backend` folder:
   ```dockerfile
   FROM golang:1.20 as builder
   WORKDIR /app
   COPY . .
   RUN go build -o main .

   FROM alpine:3.18
   WORKDIR /root/
   COPY --from=builder /app/main .
   EXPOSE 8080
   CMD ["./main"]
   ```

### **4.3 Frontend Setup**
1. Navigate to the `frontend` directory.
2. Create the React application:
   ```bash
   npx create-react-app frontend
   ```
3. Add a `Dockerfile` in the `frontend` folder:
   ```dockerfile
   FROM node:18 as builder
   WORKDIR /app
   COPY . .
   RUN npm install && npm run build

   FROM nginx:alpine
   COPY --from=builder /app/build /usr/share/nginx/html
   EXPOSE 80
   CMD ["nginx", "-g", "daemon off;"]
   ```

### **4.4 Add Docker Compose File**
Create a `docker-compose.yml` file in the root directory:
```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    networks:
      - app-network
  
  frontend:
    build:
      context: ./frontend
    ports:
      - "3000:80"
    networks:
      - app-network

default_network:
  app-network:
    driver: bridge
```

### **4.5 Set Up Jenkins**
1. Install Jenkins and required plugins:
   - Docker
   - Docker Pipeline
   - Git
   - Telegram Notifications
2. Add Docker Hub credentials to Jenkins.
3. Configure the GitHub repository webhook to trigger Jenkins builds.

### **4.6 Add Jenkinsfile**
1. Place the following `Jenkinsfile` in the root directory:
  
   pipeline {
       agent any

       environment {
           DOCKER_REGISTRY = 'your_dockerhub_username'
           BACKEND_IMAGE = "${DOCKER_REGISTRY}/backend-app"
           FRONTEND_IMAGE = "${DOCKER_REGISTRY}/frontend-app"
           TELEGRAM_BOT_TOKEN = 'YOUR_TELEGRAM_BOT_TOKEN'
           TELEGRAM_CHAT_ID = 'YOUR_TELEGRAM_CHAT_ID'
       }

       stages {
           stage('Clone Repository') {
               steps {
                   git branch: 'main', url: 'https://github.com/your_username/ci-cd-demo.git'
               }
           }

           stage('Build Backend') {
               steps {
                   dir('backend') {
                       script {
                           docker.build("${BACKEND_IMAGE}:latest")
                       }
                   }
               }
           }

           stage('Build Frontend') {
               steps {
                   dir('frontend') {
                       script {
                           docker.build("${FRONTEND_IMAGE}:latest")
                       }
                   }
               }
           }

           stage('Push to Docker Hub') {
               steps {
                   script {
                       docker.withRegistry('https://index.docker.io/v1/', 'docker-hub-credentials') {
                           docker.image("${BACKEND_IMAGE}:latest").push()
                           docker.image("${FRONTEND_IMAGE}:latest").push()
                       }
                   }
               }
           }

           stage('Deploy Services') {
               steps {
                   sh 'docker-compose up -d'
               }
           }
       }

       post {
           success {
               sendTelegramMessage("✅ Build #${BUILD_NUMBER} completed successfully!")
           }
           failure {
               sendTelegramMessage("❌ Build #${BUILD_NUMBER} failed.")
           }
       }
   }

   def sendTelegramMessage(String message) {
       sh """
       curl -s -X POST https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage \
       -d chat_id=${TELEGRAM_CHAT_ID} \
       -d text="${message}"
       """
   }
   ```

### **4.7 Run the Pipeline**
1. Push your code to GitHub.
2. Trigger a Jenkins build manually or via webhook.
3. Verify the application by accessing:
   - **Frontend**: `http://localhost:3000`
   - **Backend**: `http://localhost:8080`

---

## **5. Troubleshooting**
- Ensure Docker and Docker Compose are installed and running.
- Check Jenkins logs if builds fail.
- Verify PostgreSQL service connectivity from the backend.
- Test Telegram notifications using `curl` commands.

---

## **6. Future Enhancements**
- Add unit tests for backend and frontend.
- Deploy to a cloud platform (AWS, GCP, Azure).
- Use Kubernetes for scaling services.

---

