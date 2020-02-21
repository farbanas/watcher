pipeline {
    agent { docker { image 'golang' } }
    environment {
        GOCACHE = '/tmp/.cache'	
    }
    stages {
    	stage('test') {
		    steps {
		    	sh 'go get github.com/tebeka/go2xunit'
				sh 'go test -v | $GOPATH/bin/go2xunit > test_output.xml'
		    }
    	}
        stage('build') {
            steps {
                sh 'go build -o watcher utils.go watcher.go'
            }
        }
    }
    post {
		always {
		    archiveArtifacts artifacts: 'watcher', fingerprint: true	
		    junit 'test_output.xml'
		} 	
    }
}
