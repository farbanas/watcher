podTemplate(
	containers: [
		containerTemplate(name: 'golang', image: 'golang')
	]
) {
    node(POD_LABEL) {
		container("jnlp") {
			stage('test') {
				scm checkout
				sh 'go get github.com/tebeka/go2xunit'
				sh 'go test -v | $GOPATH/bin/go2xunit > test_output.xml'
			}
		}
		container("jnlp") {
			stage('build') {
				scm checkout
				sh 'go build -o watcher utils.go watcher.go'
			}
		}
    }
}
//GOCACHE = '/tmp/.cache'	
//archiveArtifacts artifacts: 'watcher', fingerprint: true	
//junit 'test_output.xml'
