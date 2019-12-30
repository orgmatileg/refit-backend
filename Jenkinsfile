// pipeline {
//     agent { docker { image 'maven:3.3.3' } }
//     stages {
//         stage('build') {
//             steps {
//                 sh 'ls'
//             }
//         }
//     }
// }

pipeline {
    agent any
    stages {
        stage('SCM') {
            steps {
                sh '''
                    cd /var/jenkins_home/deployed/refit/refit-backend
                    git pull origin master
                '''
            }
        }
        stage('Build') {
            steps {
                sh '''
                    cd /var/jenkins_home/deployed/refit/refit-backend
                    docker-compose build
                '''
            }
        }
        stage('Deploy') {
            steps {
                sh '''
                    cd /var/jenkins_home/deployed/refit/refit-backend
                    docker-compose up -d
                '''
            }
        }
    }
}