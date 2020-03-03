podTemplate(
	containers: [
		containerTemplate(name: "golang", image: "golang:latest", ttyEnabled: true, command: "cat")
	]
) {
    node(POD_LABEL) {
    	git "https://github.com/farbanas/watcher.git"
		stage('test') {
			container("golang") {
				sh 'go get github.com/tebeka/go2xunit'
				sh 'go test -v | $GOPATH/bin/go2xunit > test_output.xml'
			}
		}
		stage('build') {
			container("golang") {
				sh 'go build -o watcher utils.go watcher.go'
			}
		}
    }
}
//GOCACHE = '/tmp/.cache'	
//archiveArtifacts artifacts: 'watcher', fingerprint: true	
//junit 'test_output.xml'
