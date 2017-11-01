pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh './devops/scripts/build-images.sh'
      }
    }
    stage('Test') {
      steps {
        sh './devops/scripts/test.sh'
      }
    }
    stage('Push Images') {
      steps {
        sh './devops/scripts/push-images.sh'
      }
    }
    stage('Deploy') {
      steps {
        sh './devops/scripts/deploy-k8s.sh'
      }
    }
    stage ('Verify') {
      steps {
        sh './devops/scripts/verify.sh'
      }
    }
    stage ('Clean') {
      steps {
        sh './devops/scripts/clean.sh'
      }
    }
  }
}
