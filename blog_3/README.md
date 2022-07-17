# Container Restarts Watcher

## Building and running the restart watcher:

### Build
If you have earthly build system installed you can run from the current directory:
```
earthly +build-restart-watcher
```
Then the watcher binary will be located under `artifacts/watcher`

If you don't have earthly, run the following command in the current directory:
```
go build -o /bin/restart-watcher ./main.go
```

### Running
Running the watcher is straightforward, just run the following after building:
```
# this will use kubeconfig default path: "~/.kube/config"
./artifacts/restart-watcher

# if your kubeconfig resides in other path you can specify it by running :
./artifacts/restart-watcher --kubeconfig /custom/path/config
```