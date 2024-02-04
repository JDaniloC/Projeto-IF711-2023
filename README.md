# UDP/TCP Links verifier
A Go service to verify valid links of a website. It will verify recursively the links of the links and serve using TCP/UDP protocols.

## How to install
We recommend to run this project in the Github [Codespaces](https://github.com/features/codespaces), but you can run in your computer, installing the [golang](https://go.dev/doc/install) and running the following command to install the dependencies:
```bash
go mod download
```

## How to run th application
To execute the TCP/UDP server/client you need enter in the directory and run the `go run .` as showed in the next lines, where the _application_ can be `tcp_server`, `tcp_client`, `udp_server` or `udp_client`:
```bash
cd cmd/{application}
go run .
```

## How to run the tests
You can run the tests just entering in the module that you want to test and running the `go test` command. You can specify the test using `go test -v -run {test_name}`.

> Important: The [runner module](./pkg/runner/) has no a normal test, but has a benchmark test, so you need to run as `go test -bench=. -count=100`, where the count can be omitted being just `go test -bench .`.

To execute the main benchmark tests, we recommend to enter in the main folder and use the following commands:
```
go test -bench=. -benchtime=1x -benchmem -cpuprofile=cpu.out -memprofile=mem.out -timeout=0 -count=10000 | tee bench.txt

# See the graphs of CPU/MEM usage
go tool pprof -http :8081 (mem|cpu).out
```

## How it works?
The servers receive a request with the `link` url that will starts the crawl, and the number of `depth`, that specifies how far the recursion will go. The response it will be a object with `validLinks` and `invalidLinks` string arrays.

The crawler doesn't uses a crawler library, but is implemented using the `net/http` and `soup` libraries to access and parses the links, a recursive strategic and the paralellism of goroutines waited by a _wait group_.

There's a [internal](./internal/utils/) module to store the following structures:
- [Controller](./internal/utils/control.go) of the crawler, that was implemented with a RMutex to avoid race condition and other paralellism problems
- [StringSet](./internal/utils/stringset.go) object to ensures that the links added are unique (it's a string array with a hashmap), like a python set.
- [Request](./internal/utils/request.go) object to allow the Json.Unmarshall return the int value of `depth` parameter.