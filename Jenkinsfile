node {
  docker.image('golang').inside {

		stage ('Setup env vars') {
			sh 'echo "GOROOT=/usr/local/go"'
			sh 'echo "GOPATH=/home/go"'
			sh 'echo "PATH=$GOPATH/bin:$GOROOT/bin:$PATH"'
			sh 'export GOCACHE="/tmp/.cache"'
			sh 'export XDG_CACHE_HOME="/tmp/.cache"'
		}
		stage('Print go version') {
    	sh 'go version'
  	}

		stage('Get openx package') {
			sh 'go get -v github.com/YaleOpenLab/openx'
		}
  }
}
