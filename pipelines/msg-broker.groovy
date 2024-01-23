pipeline {
    environment {
        SERVICE="msg-broker"
        FOLDER="./cmd/api"
        PORT="8300"
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
                sh "sudo docker run -dp ${PORT}:${PORT} --name ${SERVICE} ${SERVICE}:latest"
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