## hsnipe

Command-line tool written in Go for comparing product prices across NZ tech retailers. 

## Using GitHub Releases
Pre-built binaries: [Releases page](https://github.com/izzymg/hsnipe/releases).  
Download the appropriate binary for your platform, extract it, and run:

```sh
cp config.example.json config.json
./hsnipe
```

## Development / Source

### Prerequisites
- Go 1.20 or newer

### Installation
Clone the repository:
```sh
git clone https://github.com/izzymg/hsnipe.git
cd hsnipe
```

Install dependencies:
```sh
go mod tidy
```

### Configuration
Copy the example config and edit as needed:
```sh
cp config.example.json config.json
```
Edit `config.json` to set your search term:
```json
{
    "searchTerm": "some hardware"
}
```

### Usage
Run the tool:
```sh
go run main.go
```
Or build and run:
```sh
go build -o hsnipe
./hsnipe
```

