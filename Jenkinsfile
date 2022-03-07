pipeline {
    
    agent any

    tools {
        go '1.17.4'
    }

    environment {
        BUILD_VERSION = 'latest'
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
                anyOf {
                    branch 'develop';
                    branch 'release/*'
                }
            } 
            steps {
                script {
                    def buildVersion = env.BRANCH_NAME.split('\\/')
                    if (buildVersion.length) > 1 {
                        env.BUILD_VERSION = buildVersion[1]
                    }

                    def image = docker.build("kutzilla/hetzner-ddns:${BUILD_VERSION}")
                }
            }      
        }
        stage('Publish') {
            when {
                branch 'release/*'
            }
            steps {
                script {
                    image.push()
                }

            }
        }
    }
}