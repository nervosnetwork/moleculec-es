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

## Options

```
Usage of moleculec-es:
  -generateTypeScriptDefinition
        True to generate TypeScript definition
  -hasBigInt
        True to generate BigInt related functions
  -inputFile string
        Input file to use
  -outputFile string
        Output file to generate, use '-' to print to stdout (default "-")
```

## ESM to CommonJS

Since moleculec-es generates ECMAScript modules, if you are using TS (< 4.5) or NodeJS(without ESM support) then you can convert ESM to CommonJS first. For example

```
npm i @babel/core @babel/cli @babel/plugin-transform-modules-commonjs
./node_modules/.bin/babel --plugins @babel/plugin-transform-modules-commonjs generated.js
```
