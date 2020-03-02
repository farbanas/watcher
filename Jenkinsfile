import jenkins.model.*

podTemplate(containers: [
	containerTemplate(name: 'golang', image: 'golang', command: 'echo "this is test"')
]) {
    node(POD_LABEL) {
		git url: 'https://github.com/farbanas/watcher.gi://github.com/farbanas/watcher.git'
    	container('golang') {
			stage('test') {
				/* scm checkout */
				sh 'go get github.com/tebeka/go2xunit'
				sh 'go test -v | $GOPATH/bin/go2xunit > test_output.xml'
			}
			stage('build') {
				/* scm checkout */
				sh 'go build -o watcher utils.go watcher.go'
			}
		}
    }
}
//GOCACHE = '/tmp/.cache'	
//archiveArtifacts artifacts: 'watcher', fingerprint: true	
//junit 'test_output.xml'
