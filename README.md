# GoHTTPSpeedTest

GoHTTPSpeedTest is a library to test the connection speed.

Currently, GoHTTPSpeedTest support Ookla's https://www.speedtest.net/ and Netflix's https://fast.com/ tests.

## Installation

To install GoHTTPSpeedTest, you can use go get command:

```
$ go get github.com/fernandomalmeida/gohttpspeedtest
```

## Usage

```
import "github.com/fernandomalmeida/gohttpspeedtest"

fastProvider, _ := gohttpspeedtest.FastProvider() // or
ooklaProvider, _ := gohttpspeedtest.OoklaProvider()

// downloadSpeed and uploadSpeed in Mbps
downloadSpeed, uploadSpeed, _ := MeasureDownloadAndUpload(provider)
```

## CLI installation

```
$ go get github.com/fernandomalmeida/gohttpspeedtest/cmd/gohttpspeedtest
```

## CLI usage

```
$ gohttpspeedtest fast # or
$ gohttpspeedtest ookla
```

## License

MIT
