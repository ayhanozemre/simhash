# SimHash

Simple and efficient SimHash implementation written in Go. SimHash is a locality-sensitive hashing algorithm used to measure text similarity.

## Installation

```bash
go get github.com/ayhanozemre/simhash
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/ayhanozemre/simhash"
)

func main() {
    text := "The cat sat on the mat"
    hash := simhash.NewSimHash(text)
    fmt.Printf("Hash: %d\n", hash)
}
```

### With Stop Words

```go
hash := simhash.NewSimHash(text, simhash.WithStopWords("the", "on", "a"))
```

### With N-gram

```go
hash := simhash.NewSimHash(text, simhash.WithNgram(4))
```

### Minimum Token Length

```go
hash := simhash.NewSimHash(text, simhash.WithMinTokenLength(3))
```

### Combined Options

```go
hash := simhash.NewSimHash(text,
    simhash.WithStopWords("the", "on"),
    simhash.WithMinTokenLength(3),
)
```

### Similarity Measurement

```go
text1 := "The cat sat on the mat"
text2 := "The cat sat on the rug"

hash1 := simhash.NewSimHash(text1)
hash2 := simhash.NewSimHash(text2)

distance := simhash.HammingDistance(hash1, hash2)
fmt.Printf("Hamming distance: %d\n", distance)
```

Smaller Hamming distance values indicate that texts are more similar.

## Testing

```bash
go test -v
```

## License

MIT License - See LICENSE file for details.

## References

- [Detecting Near-Duplicates for Web Crawling](https://static.googleusercontent.com/media/research.google.com/en//pubs/archive/33026.pdf) - Gurmeet Singh Manku, Arvind Jain, Anish Das Sarma (WWW 2007)

