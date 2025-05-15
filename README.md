## hsnipe

Command-line tool written in Go for searching and comparing product prices (such as graphics cards) across New Zealand tech retailers. It is designed to help you quickly find and compare products from multiple providers.

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

