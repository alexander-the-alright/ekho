// =============================================================================
// Auth: alex
// File: server.go
// Revn: 06-14-2024  4.0
// Func: host MOTD connection
//
// TODO: catch keyboard int. signal
//       remove ioutil import
//       add more /codes/ to handle
//       add flags, like specify port
//       write usage and errors to logfile
//       send logfile to client on request
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
//             wrote writeFile()
//*10-05-2023: commented
//             imported check() to ignore EOF
//             fixed bug in /r2/ where removing quote at final
//              position would break
//             removed deadcode from /r2/, unused "quit" bool?
//*10-17-2023: fixed bug where quotes are read from and saved to two
//              different files
//             removed implicit newline at end of qfile to stop ghost
//              line bug in /la/
//             removed dead debug print in fread()
// 05-26-2024: began writing PWFT update
//             added flag to specify transmission buffer size
// 06-12-2024: "successfully" wrote PWFT
// 06-13-2024: removed debug print statements
//             commented
//*06-14-2024: byte count-based PWFT rewrite
//
// =============================================================================

package main

import (
    "flag"      // IntVar, Parse
    "fmt"       // Println
    "io"        // io.EOF
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


// FIXME clean this two check() shit up
// print errors if they occur, quit
func checkN( err error, conn net.Conn ) {

    var errep string    // declare error report string
    switch err {        // switch case over error input
        case nil:       // if there was no error
            // then there was no error
            errep = "err is nil"
        case io.EOF:    // if there was an End Of File error
            errep = "Error: " + err.Error()     // document error
            conn.Close()    // close connection
        default:    // anything besides those are a big deal
            errep = "Fatal: " + err.Error()
    }

    var critical bool   // declare critical error flag
    // must be neither nil nor EOF
    critical = err != nil && err != io.EOF
    if critical {   // XXX critical print
        // TODO at some point, critical should be logging instead of
        // printing, only verbose should print
        fmt.Println( "check() -> ", errep )
    }
    
    if critical {       // if there was a critical error
        os.Exit( 1 )    // cut and run
    }
}


// read file and create quote list
func fread() string {
    // open file "quotes", read entire file as bytes into qfile
    qfile, err := ioutil.ReadFile( "list.q" )
    check( err )    // make sure read works
    // there is a trailing newline in the file, not removing it caused
    // bugs where a blank item was introduced into qlist and calling
    // /la/ would produce a blank line on the client side that didn't
    // really exist
    //qfile = qfile[:len( qfile ) - 1]
    // cast bytes to string, split string over newline into array
    qlist = strings.Split( string( qfile ), "\n" )
    // deal with the trailing
    return string( qfile )      // return file as string
}


// write qfile back to file
func writeFile( f string ) {
    // open file list.q as variable file
    file, err := os.Create( "list.q" )
    check( err )    // make sure open worked
    // write input string ( quote file ) back to file
    _, err = file.WriteString( f )
    check( err )    // make sure writeback worked
}


// handle connections ( new and improved )
func handle( conn net.Conn ) {

    // create buffer to hold read message
    buffer := make( []byte, tsize )
    //var buffer [256]byte
    // read n bytes from client
    n, err := conn.Read( buffer[:] )
    checkN( err, conn ) // make sure read worked

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
                resp = "ba@qr@!"    // bad answer
            } else {        // comm[1] was a number, but...
                switch {    // how do we know it's a valid number
                    case num > len( qlist ):
                        // bad answer in quote request
                        resp = "ba@qr@>"
                        // why bad? num is >
                    case num < 0:
                        // bad answer in quote request
                        resp = "ba@qr@-"
                        // why bad? num is -
                    default:
                        // good num, get quote, send off
                        resp = "qa@" + qlist[num]   // quote answer
                }
            }
        case "lr":      // /lr/, list request
            resp = "la@"    // begin crafting response
            // file too big for single transmission
            // initial header size is 5, la@1@, so counting that into
            // the transmission size, the file must be longer than 251
            // ( default ) bytes to reach this block of code
            if len( qfile ) > ( tsize - 5 ) {
                if command[1] == "init" {   // first request
                    // add message number, delim, and the first
                    // transmittable bytes
                    resp = resp + "0@" + qfile[:( tsize - 5 )]
                } else {    // subsequent requests
                    // convert byte count to int
                    bytes, err := strconv.Atoi( command[1] )
                    check( err )    // make sure convert worked
                    // replace tx'd bytes and delim to new header
                    resp = resp + command[1] + "@"
                    rlen := len( resp )     // capture header length
                    // slices are always of size "txsize - header
                    // length", so find end byte
                    end := bytes + ( tsize - rlen )
                    // if end goes beyond the bounds of the list
                    if end > len( qfile ) {
                        // replace byte count  with X to signal
                        // final message, concat the rest of the file,
                        // beginning with the start index
                        resp = "la@X@" + qfile[bytes:]
                    } else {
                        // get slice from start index to end, concat
                        // with response
                        resp = resp + qfile[num:end]
                    }
                }
            } else {    // file is smaller than the tx size
                // just concat entire list to send
                resp = resp + qfile    // list answer
            }
        case "ar":      // /ar/, add request
            // append new quote to list
            qlist = append( qlist, command[1] )
            // also append new quote to file variable
            qfile = qfile + "\n" + command[1]
            writeFile( qfile )
            resp = "aa@success"     // add answer
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
                resp = "ba@r2@!"    // bad answer
            } else {        // comm[1] was a number, but...
                num = num - 1
                switch {    // how do we know it's a valid number
                            // valid numbers are from 0 - len-1
                    case num > len( qlist ):
                        // bad answer in quote request
                        resp = "ba@r2@>"
                        // why bad? num is >
                    case num < 0:
                        // bad answer in quote request
                        resp = "ba@r2@-"
                        // why bad? num is -
                    // if num is not OUT OF RANGE, must be in range
                    default:
                        // replace the quote at that position with
                        // final quote
                        qlist[num] = qlist[len( qlist ) - 1]
                        // remove second copy of quote
                        qlist = qlist[:len( qlist ) - 1]
                        // O(1) remove from list at any position
                        // must entirely recreate quote string
                        var f string = ""   // init empty string
                        // iterate over list, disregard index
                        for _, v := range qlist {
                            // add next quote to file and separate
                            // quotes with newline
                            f = f + v + "\n"
                        }
                        // remove trailing newline, to prevent ghost
                        f = f[:len( f ) - 1]    // quote in lr
                        // reassign global string variable to
                        // recreated quote string
                        qfile = f
                        writeFile( qfile )  // write back to file
                        // let client know remove was successful
                        resp = "r2@success"
                }
            }
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
    checkN( err, conn ) // make sure write worked

    handle( conn )  // call handle() again to get response from client

}


// globals
var qlist []string      // list of quotes
var qfile string        // file object
var tsize int           // size of transmission buffer


// main, create port and wait for connection
func main() {

    // flag for default transmission size
    flag.IntVar( &tsize, "s", 256, "transmission size ( bytes )" )
    flag.Parse()

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

