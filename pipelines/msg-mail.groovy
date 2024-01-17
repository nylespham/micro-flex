pipeline {
    environment {
        SERVICE="msg-mail"
        SERVICE_BINARY="/usr/local/bin/msg-mail"
    }
    agent gcp
    stages {
        stage("Build Image"){
            steps {
                sh "docker build -t msg-mail:latest -f dockerfiles/Dockerfile . --build-arg SERVICE_NAME=${SERVICE} --build-arg BINARY=${SERVICE_BINARY}"
            }
        }
        stage("Run Container"){
            steps {
                sh "docker run -dp 4100:4100 ${SERVICE}:latest"
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