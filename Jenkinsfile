pipeline {
    
    agent any

    tools {
        go '1.17.4'
    }

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
        stage('Package') {
            when {
                branch 'develop'
            } 
            
            docker.build("kutzilla/hetzner-ddns:${env.BRANCH_NAME}")
     
        }
        stage('Publish') {
            when {
                branch 'master'
            }
            steps {
                echo 'Add publish'
            }
        }
    }
}