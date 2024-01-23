pipeline {
    environment {
        SERVICE="msg-listener"
        FOLDER="."
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
        // stage("Discard old Container"){
        //     steps {
        //         dir("./compose") {
        //             sh "sudo docker compose down ${SERVICE}"
        //         }
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