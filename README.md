# moleculec-es
ECMAScript plugin for the molecule serialization system

# Usage

First, download precompiled binary from [releases](https://github.com/xxuejie/moleculec-es/releases) page, and put the binary in our PATH. You could also clone and build the project from source. The only dependency is Golang here.

```
$ cargo install moleculec
$ moleculec --language - --schema-file "your schema file" --format json > /tmp/schema.json
$ moleculec-es -inputFile /tmp/schema.json -outputFile "your JS file"
```

Generated file from this project follows latest ECMAScript standard, it is also orgnaized as an ECMAScript module to allow for tree shaking.
