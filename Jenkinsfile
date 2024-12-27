pipeline {
    agent any

    environment {
        DOCKER_REGISTRY = 'takieulong'
        BACKEND_IMAGE = "takieulong/backend-app"
        FRONTEND_IMAGE = "takieulong/frontend-app"
        POSTGRES_IMAGE = "postgres:15"
        DOCKER_TAG = 'latest'
        TELEGRAM_BOT_TOKEN = '7710797141:AAG07MhFy_x5XP_7XlUvSjaHwxA4B6DR8SI'
        TELEGRAM_CHAT_ID = '-1002194149690'
    }

    stages {
        stage('Clone Repository') {
            steps {
                git branch: 'master', url: 'https://github.com/TaKieuLong/final.git'
            }
        }

        stage('Build Backend') {
            steps {
                dir('backend') {
                    script {
                        docker.build("${BACKEND_IMAGE}:${DOCKER_TAG}")
                    }
                }
            }
        }

        stage('Build Frontend') {
            steps {
                dir('frontend') {
                    script {
                        docker.build("${FRONTEND_IMAGE}:${DOCKER_TAG}")
                    }
                }
            }
        }

        stage('Push to Docker Hub') {
            steps {
                script {
                    docker.withRegistry('https://index.docker.io/v1/', 'docker-hub-credentials') {
                        docker.image("${BACKEND_IMAGE}:${DOCKER_TAG}").push()
                        docker.image("${FRONTEND_IMAGE}:${DOCKER_TAG}").push()
                    }
                }
            }
        }

        stage('Start PostgreSQL') {
            steps {
                script {
                    sh """
                    docker network create dev || echo "Network 'dev' already exists"
                    docker run -d --rm --name postgres --network dev -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=productdb -p 5432:5432 ${POSTGRES_IMAGE}
                    sleep 10
                    """
                }
            }
        }

        stage('Deploy Backend and Frontend') {
            steps {
                script {
                    sh """
                    # Dừng và xóa container cũ nếu tồn tại
                    docker container stop backend-app || echo "No existing Backend container"
                    docker container stop frontend-app || echo "No existing Frontend container"

                    # Chạy container Backend
                    docker run -d --rm --name backend-app --network dev -p 8084:8080 ${BACKEND_IMAGE}:${DOCKER_TAG}

                    # Chạy container Frontend
                    docker run -d --rm --name frontend-app --network dev -p 3000:80 ${FRONTEND_IMAGE}:${DOCKER_TAG}
                    """
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }

        success {
            sendTelegramMessage("✅ CI/CD Pipeline Build #${BUILD_NUMBER} completed successfully! ✅")
        }

        failure {
            sendTelegramMessage("❌ CI/CD Pipeline Build #${BUILD_NUMBER} failed. ❌")
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
