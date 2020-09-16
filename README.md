# rosetta-inspect

Rosetta server implementation inspector

## Installation

### Binary

Not available yet

### Go

If you have Go available on your machine you can install the inspector with:

```bash
go get github.com/figment-networks/rosetta-inspector
```

## Usage

```
Usage of ./rosetta-inspector:
  -listen string
    	Listen address (default "0.0.0.0:5555")
  -url string
    	Rosetta server URL
```

When you have a Rosetta server running on `http://localhost:8080` you can start
the inspector with:

```bash
rosetta-inspector -url=http://localhost:8080
```

Then you can open the inspector UI at `http://localhost:5555`

## License

Apache License v2.0
