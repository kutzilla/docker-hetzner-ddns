pipeline {
    
    agent any

    tools {
        go '1.17.4'
    }
    
    environment {
        GO111MODULE = 'on'
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'develop', url: 'https://gitea.matthias-kutz.com/matthias/docker-hetzner-ddns.git'

            }
        }
        
        stage('Install') {
            steps {
                sh 'go mod download'
            }
        }
        stage('Build') {
            steps {
                sh 'go build -o hetzner-ddns ./cmd/hetzner-ddns'
            }
        }
        stage('Test') {
            steps {
                sh 'go test ./pkg/*'
            }
        }
    }
}