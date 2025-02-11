# Todo CLI

A simple CLI for managing your todo list.

## Installation

```bash
go install github.com/dt-az/todo-cli@latest

```
or clone the repository.

## Build

```bash
go build -o todo main.go
```

## Usage

```bash
./todo add -text "Buy groceries"
./todo list
./todo complete -number 1
./todo clear
./todo help
```

## Testing

```bash
go test
```

## Coverage

```bash
go test -cover
```
