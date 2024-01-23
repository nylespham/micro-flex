pipeline {
    environment {
        SERVICE="db-postgres"
    }
    agent {
        label "jenkins-02"
    }
    stages {
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