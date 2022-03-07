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
            steps {
                script {
                    docker.build("kutzilla/hetzner-ddns:latest")
                }
            }      
        }
        stage('Publish') {
            when {
                branch 'develop'
            }
            steps {
                script {
                    def latestTag = sh(returnStdout: true, script: "git tag --sort version:refname | tail -1").trim()
                    echo "Publishing latest ${latestTag}"
                }

            }
        }
    }
}