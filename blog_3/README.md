# Container Restarts Watcher

## Building and running the restart watcher:

### Build
Run the following command in the current directory:
```
go build -o /bin/restart-watcher ./main.go
```

Or if you prefer, you can build docker image using:
```
docker build .
```

### Running
Running the watcher is straightforward, just run the following after building:
```
# this will use kubeconfig default path: "~/.kube/config"
./artifacts/restart-watcher

# if using docker
docker run -v ~/.kube:/root/.kube <image_name>

# if your kubeconfig resides in other path you can specify it by running :
./artifacts/restart-watcher --kubeconfig /custom/path/config
```