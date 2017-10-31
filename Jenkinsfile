pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh './devops/scripts/build-images.sh'
      }
    }
    stage('Push Images') {
      steps {
        sh './devops/scripts/push-images.sh'
      }
    }
    stage('Test') {
      steps {
        echo 'Hello, OpenPitrix'
        sh 'docker images'
      }
    }
    stage('Deploy') {
      steps {
        sh './devops/scripts/deploy-k8s.sh'
      }
    }
    stage ('Verify') {
      steps {
        echo 'Verifying...'
      }
    }
    stage ('Clean') {
      steps {
        echo 'Cleaning...'
        sh 'docker image rm openpitrix/openpitrix:dev'
        sh 'docker image rm openpitrix'
      }
    }
  }
}
