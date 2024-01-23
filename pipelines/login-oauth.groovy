pipeline {
    environment {
        SERVICE="login-oauth"
        FOLDER="./cmd/api"
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
        // stage("Discard Old Container"){
        //     steps {
        //         sh "sudo docker rm -f \$(sudo docker ps -qaf 'name=${SERVICE}')"
        //     }
        // }
        stage("Run Container"){
            steps {
                dir("./compose") {
                    sh "sudo docker compose up -d ${SERVICE}"
                }
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