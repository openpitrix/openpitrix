pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh './deploy/devops/scripts/build-images.sh'
      }
    }
    stage('Test') {
      steps {
        sh './deploy/devops/scripts/test.sh'
      }
    }
    stage('Push Images') {
      steps {
        sh './deploy/devops/scripts/push-images.sh'
      }
    }
    stage('Deploy') {
      steps {
        sh './deploy/kubernetes/scripts/deploy-k8s.sh -v dev -b -d -m'
      }
    }
    stage ('Verify') {
      steps {
        sh './deploy/devops/scripts/verify.sh'
      }
    }
    stage ('Clean') {
      steps {
        sh './deploy/kubernetes/scripts/clean.sh'
      }
    }
  }
}
