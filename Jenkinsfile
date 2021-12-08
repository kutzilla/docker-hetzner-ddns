pipeline {
    
    agent any

    stages {
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
