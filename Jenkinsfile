pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh 'docker build -t openpitrix .'
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
        sh './scripts/deliver.sh'
      }
    }
  }
}
