etcd_watcher_experiment
=======================

An experiment with etcd v2 client watcher.

## Set up

Install [Glide | Package Management For Go](https://glide.sh/) and put it in your path.

Get and build the `etcd_watcher_experiment` command.

```
go get -d github.com/hnakamur/etcd_watcher_experiment
cd "$GOPATH/src/github.com/hnakamur/etcd_watcher_experiment"
glide install
go build
```

Run a local etcd cluster.

```
(cd vendor/github.com/coreos/etcd && go build)
mkdir bin
cp vendor/github.com/coreos/etcd/etcd bin/
cp vendor/github.com/coreos/etcd/Procfile .
go get -u github.com/mattn/goreman/...
goreman start
```

In another terminal, start a `etcd_watcher_experiment` and watch the key `key1`.

```
./etcd_watcher_experiment watch
```

In another terminal, set, modify or delete the key `key1.

```
./etcd_watcher_experiment set -value value1
./etcd_watcher_experiment set -value value1,value2
./etcd_watcher_experiment set -value ""
./etcd_watcher_experiment set -value value3
./etcd_watcher_experiment delete
./etcd_watcher_experiment set -value value4
```

## License
MIT
