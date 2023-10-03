// =============================================================================
// Auth: Alex Celani
// File: server.go
// Revn: 10-02-2023  1.0
// Func: host connection, reply to speak.go
//
// TODO: fix remove
//       write back to file, re-read file
//       add more /codes/ to handle
//       add flags
// =============================================================================
// CHANGE LOG
// -----------------------------------------------------------------------------
// 05-13-2022: init
// 08-29-2023: gutted and rewrote
//             commented
// 09-12-2023: imported io/ioutil, strings, strconv
//             added size conversation
//             added file read function, fread()
//             commented
// 09-21-2023: began rewrite of handle()
// 09-26-2023: got /list/ command working
//             got /add/ command (mostly) working
//*10-02-2023: commented handle()
//             added checks to /qr/ and /r2/ to make sure recv'd
//              length is valid
// =============================================================================

package main

import (
    "fmt"       // Println
    "io/ioutil" // ReadFile
    "net"       // ResovleTCPAddr, ListenTCP, listener.Accept
                // conn.Read,Write
    "os"        // Exit
    "strconv"   // Itoa
    "strings"   // Split
)


// handle errors catastrophically
func check( err error ) {
    if err != nil {    // error is not nil on error
        // print error
        fmt.Println( "Fatal: ", err.Error() )
        os.Exit( 1 )    // bail
    }
}


// read file and create quote list
func fread() string {
    // open file "quotes", read entire file as bytes into qfile
    qfile, err := ioutil.ReadFile( "quotes" )
    check( err )    // make sure read works
    // cast bytes to string, split string over newline into array
    qlist = strings.Split( string( qfile ), "\n" )
    // remove last item in list ( trailing newline leaves empty item )
    qlist = qlist[:len( qlist ) - 1]
    if false {      // XXX debug print ( next )
        for _, q := range qlist {   // iterate over all quotes in list
            fmt.Println( q )        // print quote
        }
    }
    return string( qfile )      // return file as string
}


// write qfile back to file
func writeFile( f string ) {
    file, err := os.Create( "quotes" )
    check( err )
    _, err = file.WriteString( f )
    check( err )
}


// handle connections ( new and improved )
func handle( conn net.Conn ) {

    // create buffer to hold read message
    var buffer [256]byte
    // read n bytes from client
    n, err := conn.Read( buffer[:] )
    check( err )    // make sure read worked

    // cast message to string, take first n bytes, split over @
    command := strings.Split( string( buffer[:n] ), "@" )

    // declare response string
    var resp string
    switch command[0] {     // switch/case over opcode
        case "sr":      // /sr/, size request
            // get length of list, convert to string
            // XXX probably no reason to have this declaration scheme
            var num string = strconv.Itoa( len( qlist ) )
            // concat and send
            resp = "sa@" + num      // size answer
        case "qr":      // /qr/, quote request
            // convert message argument to string to make sure it's a
            // valid number
            num, nerr := strconv.Atoi( command[1] )
            // nerr will not be nil if comm[1] is not a number
            if nerr != nil {
                // formulate an error message
                // TODO flesh this out more
                resp = "ba@Number not understood"   // bad answer
            } else {        // comm[1] was a number, but...
                switch {    // how do we know it's a valid number
                    case num > len( qlist ):
                        // bad answer in quote request
                        resp = "ba@qr@>"
                        // why bad? num is >
                    case num < 0:
                        // bad answer in quote request
                        resp = "bad@qr@-"
                        // why bad? num is -
                    default:
                        // good num, get quote, send off
                        resp = "qa@" + qlist[num]   // quote answer
                }
            }
        case "lr":      // /lr/, list request
            // just concat entire list to send
            resp = "la@" + qfile[:len( qfile ) - 1]     // list answer
            // TODO what happens if this is longer than 256 char?
        case "ar":      // /ar/, add request
            // append new quote to list
            qlist = append( qlist, command[1] )
            // also append new quote to file variable
            qfile = qfile + command[1] + "\n"
            writeFile( qfile )
            resp = "aa@success"     // add answer
            // TODO write to file
        case "r1":      // /r1/, remove request 1
            // same as list, just send whole file
            resp = "r1@" + qfile    // remove request 1
        case "r2":      // /r2/, remove request 2
            // convert message argument to string to make sure it's a
            // valid number
            num, nerr := strconv.Atoi( command[1] )
            // nerr will not be nil if comm[1] is not a number
            if nerr != nil {
                // formulate an error message
                // TODO flesh this out more
                resp = "ba@Number not understood"   // bad answer
            } else {        // comm[1] was a number, but...
                switch {    // how do we know it's a valid number
                    case num > len( qlist ):
                        // bad answer in quote request
                        resp = "ba@r2@>"
                        // why bad? num is >
                    case num < 0:
                        // bad answer in quote request
                        resp = "ba@r2@-"
                        // why bad? num is -
                    default:
                        qlist[num + 1] = qlist[len( qlist ) - 1]
                        qlist = qlist[:len( qlist ) - 1]
                        var f string = ""
                        for _, v := range qlist {
                            f = f + v + "\n"
                        }
                        f = f[:len( f ) - 1]
                        qfile = f
                        writeFile( qfile )
                        resp = "r2@" + qfile
                }
            }
            /*
            resp = "r2@"    // remove request 2
            // 
            index, nerr := strconv.Atoi( command[1] )
            if nerr == nil {
                if index < len( qlist ) {
                    qlist[index] = qlist[len( qlist ) - 1]
                    qlist = qlist[:len( qlist ) - 1]
                    //writeFile() //TODO
                    resp = resp + "success"
                } else {
                    resp = resp + "error"
                }
            } else {
                resp = resp + "error"
            }
            */
        case "te":      // /te/, transaction end
            // TODO should I be sending something?
            // TODO is it better to close now or outside of handle()
            conn.Close()    // close connection and leave
            return          // return from handle(), finish connection
        default:        // opcode not recognized
            // TODO there should be more being down, here, right?
            return          // return from handle(), finish connection
    }

    // write response to client
    _, err = conn.Write( []byte( resp ) )
    check( err )    // make sure write worked

    handle( conn )  // call handle() again to get response from client

}


// globals
var qlist []string      // list of quotes
var qfile string        // file object
var quit bool


// main, create port and wait for connection
func main() {

    qfile = fread() // read quote file, init ( global ) list of quotes

    service := ":1300"  // create service on ip and port

    // resolve ip address and port
    addr, err := net.ResolveTCPAddr( "tcp", service )
    check( err )    // make sure ip resolves

    // create listener object from ip:port
    listener, err := net.ListenTCP( "tcp", addr )

    for {
        // wait, create connection when found
        conn, err := listener.Accept()
        check( err )    // make sure connection works
    
        go handle( conn )  // handle connection
    }

    os.Exit( 0 )    // exeunt

}

