# az-lang

An English-like programming language that reads like natural language. Write code that speaks human.

```
set greeting to "Hello"
say greeting plus ", World!"
```

## Features

- **Natural Language Syntax** - Code reads like English sentences
- **HTTP Client** - Make web requests with `fetch`, `send`, `put`, `delete`
- **Web Server** - Host APIs with natural language route definitions
- **JSON Support** - Parse and encode JSON data
- **English Numbers** - Use words like `forty two` instead of `42`

## Installation

### Prerequisites

- Go 1.19 or higher

### Build from Source

```bash
git clone https://github.com/henrydennis/az-lang.git
cd az-lang
go build -o abc .
```

### Run

```bash
# Run a script
./abc examples/hello.abc

# Interactive REPL
./abc
```

## Quick Start

Create a file `hello.abc`:

```
say "Hello, World!"
```

Run it:

```bash
./abc hello.abc
```

## Language Reference

### Variables

```
set name to "Alice"
set age to 25
set count to zero
```

### Arithmetic

Use English words for operators:

```
set result to 10 plus 5        # 15
set result to 10 minus 3       # 7
set result to 4 times 3        # 12
set result to 20 divided by 4  # 5
```

Modify variables in place:

```
set count to 0
increase count by 1
decrease count by 1
```

### String Concatenation

```
set greeting to "Hello, "
set message to greeting plus "World!"
say message    # Hello, World!
```

### Comparisons

```
if x equals y then
    say "equal"
done

if x is greater than y then
    say "x is bigger"
done

if x is less than y then
    say "x is smaller"
done
```

### Logical Operators

```
if x equals 1 and y equals 2 then
    say "both true"
done

if x equals 1 or y equals 2 then
    say "at least one true"
done

if not x equals 1 then
    say "x is not 1"
done
```

### Conditionals

```
if age is greater than 18 then
    say "adult"
otherwise
    say "minor"
done
```

### Loops

**While loops:**

```
set count to 5
while count is greater than 0 do
    say count
    decrease count by 1
done
```

**For each loops:**

```
set items to a list of "apple" and "banana" and "cherry"
for each item in items do
    say item
done
```

### Lists

```
# Create a list
set fruits to a list of "apple" and "banana" and "cherry"

# Get length
say length of fruits    # 3

# Access by index (1-indexed)
say item 1 from fruits  # apple

# Append to list
append "orange" to fruits
```

### Functions

```
to greet with name
    say "Hello, " plus name plus "!"
done

greet with "Alice"    # Hello, Alice!
```

**With return values:**

```
to add with a and b
    return a plus b
done

set result to add with 5 and 3
say result    # 8
```

### Input/Output

```
say "What is your name?"
ask into name
say "Hello, " plus name plus "!"
```

### English Numbers

Use words for numbers 0-999,999:

```
set x to forty two
set y to one hundred twenty three
set big to one million
say x plus y    # 165
```

## HTTP Client

### GET Request

```
fetch from "https://api.example.com/users" into response
say body of response
say status of response
```

### POST Request

```
set data to "{\"name\": \"Alice\"}"
send data to "https://api.example.com/users" into response
say body of response
```

### PUT Request

```
set data to "{\"name\": \"Bob\"}"
put data to "https://api.example.com/users/1" into response
```

### DELETE Request

```
delete from "https://api.example.com/users/1" into response
```

### With Headers

```
set headers to a list of "Authorization: Bearer token123" and "Content-Type: application/json"
fetch from "https://api.example.com/data" with headers into response
```

### Response Properties

```
say body of response              # Response body as string
say status of response            # HTTP status code (200, 404, etc.)
say header "Content-Type" from response
```

## JSON

### Parse JSON

```
set jsonString to "{\"name\": \"Alice\", \"age\": 30}"
parse jsonString as json into data
say field "name" from data    # Alice
say field "age" from data     # 30
```

### Nested Fields

```
set json to "{\"user\": {\"name\": \"Bob\"}}"
parse json as json into data
say field "user.name" from data    # Bob
```

### Encode to JSON

```
set items to a list of "apple" and "banana"
encode items as json into result
say result    # ["apple","banana"]
```

## Web Server

### Basic Server

```
when request at "/" do
    reply with "Hello, World!"
done

serve on 8080
```

### HTTP Method-Specific Routes

```
when fetch at "/users" do
    reply with "List of users"
done

when send at "/users" using req do
    reply with "User created" with status 201
done

when put at "/users" using req do
    reply with "User updated"
done

when delete at "/users" using req do
    reply with "User deleted"
done
```

Note: `fetch` = GET, `send` = POST

### Accessing Request Data

```
when send at "/api/data" using req do
    set requestMethod to method of req
    set requestPath to path of req
    set requestBody to body of req
    set name to query "name" from req
    set auth to header "Authorization" from req

    reply with "Received: " plus requestBody
done
```

### JSON Responses

```
when fetch at "/api/users" do
    set users to a list of "Alice" and "Bob" and "Charlie"
    reply with users as json
done
```

The `as json` modifier automatically:
- Encodes the data as JSON
- Sets `Content-Type: application/json` header

### Custom Status Codes

```
when send at "/api/users" using req do
    reply with "Created" with status 201
done

when fetch at "/api/error" do
    reply with "Not Found" with status 404
done
```

### Custom Headers

```
when fetch at "/api/data" do
    reply with "data" with header "X-Custom" as "value"
done
```

### Function Reference Handlers

```
to handleHome with req
    reply with "Welcome!"
done

to handleAbout with req
    reply with "About us"
done

route "/" to handleHome
route "/about" to handleAbout

serve on 8080
```

### Background Server

```
when fetch at "/health" do
    reply with "OK"
done

serve on 8080 in background

say "Server running..."
# Script continues executing
```

### Stop Server

```
stop server              # Stop all servers
stop server on 8080      # Stop specific server
```

### Complete REST API Example

```
say "Starting API..."

set users to a list of "Alice" and "Bob"

when fetch at "/users" do
    reply with users as json
done

when fetch at "/users/greet" using req do
    set name to query "name" from req
    if name equals null then
        set name to "Guest"
    done
    reply with "Hello, " plus name plus "!"
done

when send at "/users" using req do
    parse body of req as json into data
    set name to field "name" from data
    append name to users
    reply with users as json with status 201
done

serve on 8080
```

Test with:

```bash
curl http://localhost:8080/users
curl http://localhost:8080/users/greet?name=Developer
curl -X POST -d '{"name":"Charlie"}' http://localhost:8080/users
```

## Examples

The `examples/` directory contains sample programs:

| File | Description |
|------|-------------|
| `hello.abc` | Hello World |
| `countdown.abc` | While loop countdown |
| `factorial.abc` | Recursive factorial function |
| `fizzbuzz.abc` | Classic FizzBuzz |
| `lists.abc` | List operations |
| `english_numbers.abc` | Using word numbers |
| `http_example.abc` | HTTP GET request |
| `http_post.abc` | HTTP POST request |
| `server.abc` | Full REST API server |
| `server_background.abc` | Background server |
| `api_aggregator.abc` | API aggregator fetching from external services |
| `notes_api.abc` | Simple notes CRUD API |

## Architecture

az-lang is a tree-walking interpreter written in Go:

```
Source Code (.abc)
       ↓
    [Lexer]     → Tokenizes into token stream
       ↓
    [Parser]    → Builds Abstract Syntax Tree
       ↓
  [Interpreter] → Evaluates AST and produces output
```

### Project Structure

```
az-lang/
├── main.go           # Entry point, REPL
├── token/
│   └── token.go      # Token definitions
├── lexer/
│   └── lexer.go      # Tokenizer
├── ast/
│   └── ast.go        # AST node definitions
├── parser/
│   └── parser.go     # Recursive descent parser
├── object/
│   └── object.go     # Runtime value types
├── interpreter/
│   └── interpreter.go # Tree-walking evaluator
└── examples/         # Example programs
```

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
