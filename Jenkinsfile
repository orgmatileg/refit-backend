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
        stage('Build') {
            steps {
                sh '''
                    cd /var/jenkins_home/deployed/refit/refit-backend
                    docker-compose build
                    docker-compose up -d
                '''
            }
        }
    }
}