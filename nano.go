/*
    nano.go
        Main entrypoint for the go-nano client

    Author: pat@dcrypt.io
*/
package main

import (
    "github.com/pdepip/go-nano/node"
)

func main() {

    // Connect to default Nano Node
    err := node.FindDefaultNode()
    if err != nil {
        panic(err)
    }

    node.ListenUDP()
}
