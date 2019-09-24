node {
  docker.image('golang').inside {

		stage('Print go version') {
    	sh 'go version'
  	}

		stage('Get openx package') {
			sh 'echo "GOROOT=/usr/local/go"'
			sh 'echo "GOPATH=/home/go"'
			sh 'echo "PATH=$GOPATH/bin:$GOROOT/bin:$PATH"'
			sh 'echo "$GOPATH"'
			sh 'export GOCACHE="/tmp/.cache"'
			sh 'export XDG_CACHE_HOME="/tmp/.cache"'
			sh 'go get -v github.com/YaleOpenLab/openx'
		}
  }
}
