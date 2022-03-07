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
                branch 'release/*'
            }
            steps {
                script {
                    def prefix, releaseVersion = env.BRANCH_NAME.split('\\/')
                    echo 'Publishing ' + releaseVersion
                }

            }
        }
    }
}