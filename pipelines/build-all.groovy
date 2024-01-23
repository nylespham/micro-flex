pipeline {
    agent {
        label "jenkins-02"
    }
    stages {
        stage("login-oauth"){
            steps {
                build "login-oauth-sit"
            }
        }
        stage("msg-broker"){
            steps {
                build "msg-broker-sit"
            }
        }
        stage("msg-logger"){
            steps {
                build "msg-logger-sit"
            }
        }
        stage("msg-listener"){
            steps {
                build "msg-listener-sit"
            }
        }
        stage("msg-mail"){
            steps {
                build "msg-mail-sit"
            }
        }
        stage("web-agent"){
            steps {
                build "web-agent-sit"
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