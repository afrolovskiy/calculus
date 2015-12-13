# Calculator

Simple client-server calculator with thrift as rpc. Supports addition, subtraction, multiplication and division.

[Godep][godep] used for dependency management. External code is vendored and imports are rewritten. So you don't need godep to install.

## Installation

Was tested on OSX. Should also work on Unix systems.

Requires Go 1.3+. Make sure you have Go properly installed, including setting up your GOPATH. 

Create directory in your GOPATH and move folder with project there:

    $ mkdir -p $GOPATH/src/github.com/afrolovskiy
    $ mv calculus $GOPATH/src/github.com/afrolovskiy

Now you can install it:

    $ go install github.com/afrolovskiy/calculus/cmd/calculus
    $ go install github.com/afrolovskiy/calculus/cmd/calculusd

## Usage

First you must run the server. Server will start to listen on `localhost:9090`. It can be changed with `--addr` flag for both server and client.

    $ calculusd
    2015/01/17 23:22:23 starting server on localhost:9090

Now you can use `calculus` client to perform operations by passing expression from command line or via stdin.

    $ calculus
    (1 + 5) /2
    3
    $ echo "(3 + 5) * (7 - 2)" | calculus
    40
