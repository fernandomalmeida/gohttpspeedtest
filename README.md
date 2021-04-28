# GoHTTPSpeedTest

GoHTTPSpeedTest is a library to test the connection speed.

Currently, GoHTTPSpeedTest support Ookla's https://www.speedtest.net/ and Netflix's https://fast.com/ tests.

## Installation

To install GoHTTPSpeedTest, you can use go get command:

```shell
$ go get github.com/fernandomalmeida/gohttpspeedtest
```

## Usage

```go
import "github.com/fernandomalmeida/gohttpspeedtest"

fastProvider, _ := gohttpspeedtest.FastProvider() // or
ooklaProvider, _ := gohttpspeedtest.OoklaProvider()

// downloadSpeed and uploadSpeed in Mbps
downloadSpeed, uploadSpeed, _ := MeasureDownloadAndUpload(provider)
```

## CLI installation

```shell
$ go get github.com/fernandomalmeida/gohttpspeedtest/cmd/gohttpspeedtest
```

## CLI usage

```shell
$ gohttpspeedtest fast # or
$ gohttpspeedtest ookla
```

## Current test coverage

```shell
$ go test -cover .
ok      github.com/fernandomalmeida/gohttpspeedtest     0.802s  coverage: 93.4% of statements
```

## License

MIT
