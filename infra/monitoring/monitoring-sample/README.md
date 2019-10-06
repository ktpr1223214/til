# monitoring-sample

## setup
* install
``` bash
$ brew install jsonnet jq
$ go get -u github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb
$ jb init
```

* setup
``` bash
$ jb install https://github.com/grafana/grafonnet-lib
$ ./create_grafana_folder.sh playground-for-testing "playground for testing"
```