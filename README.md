# rosetta-inspector

Rosetta server implementation inspector

## Installation

### Binary

Download a binary release from [Github](https://github.com/figment-networks/rosetta-inspector/releases)

### Go

If you have Go available on your machine you can install the inspector with:

```bash
go get github.com/figment-networks/rosetta-inspector
```

### Docker

Pull the image:

```bash
docker pull figmentnetworks/rosetta-inspector
```

Start the inspector container:

```bash
docker run -p 5555:5555 figmentnetworks/rosetta-inspector -url=http://rosetta-server:port
```

You should be able to view the UI by visiting http://localhost:5555

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
