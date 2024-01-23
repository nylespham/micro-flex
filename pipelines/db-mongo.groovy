pipeline {
    environment {
        SERVICE="db-mongo"
    }
    agent {
        label "jenkins-02"
    }
    stage("Run Container"){ 
            steps {
                dir("./compose") {
                    sh "sudo docker compose up -d ${SERVICE}"
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