pipeline {
    environment {
        SERVICE="web-agent"
        FOLDER="./cmd/web"
    }
    agent {
        label "jenkins-02"
    }
    stages {
        stage("Build Image"){
            steps {
                sh "sudo docker build -t ${SERVICE}:latest -f dockerfiles/Dockerfile . --build-arg SERVICE_NAME=${SERVICE} --build-arg FOLDER=${FOLDER}"
            }
        }
        stage("Run Container"){
            steps {
                sh "service_name='$(docker ps -qaf 'name=quirky_jepsen')'"
            }
        }
        stage("Run Container"){
            steps {
                sh "sudo docker run -dp 8100:8100 ${SERVICE}:latest"
            }
        }
    }
    post {
        always {
            cleanWs(cleanWhenNotBuilt: false,
                    deleteDirs: true,
                    disableDeferredWipeout: true,
                    notFailBuild: true,
                    patterns: [[pattern: '.gitignore', type: 'INCLUDE'],
                               [pattern: '.propsfile', type: 'EXCLUDE']])
        }
    }
}