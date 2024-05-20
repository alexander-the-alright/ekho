// =============================================================================
// Auth: alex
// File: client.go
// Revn: 05-19-2024  0.2
// Func: ask server for a message of the day quote
//
// TODO: gut
//       write
// =============================================================================
// CHANGE LOG
// -----------------------------------------------------------------------------
// 05-18-2024:  gutted
// 05-19-2024:  fixed
//              commented
//
// =============================================================================

package main

import ( 
    //"io"        // io.EOF
    "flag"      // BoolVar, StringVar, Parse
    "fmt"       // Println
    "net"       // ResolveTCPAddr, DialTCP, conn.Write,Read
    "os"        // Exit
    //"strconv"   // Atoi, Itoa
    //"strings"   // Split
)


/*
// print errors if they occur, quit
func check( err error ) {

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


// get server response and print
func handle( conn net.Conn ) {

    var buffer [tsize]byte      // create buffer to hold response
    // read n bytes from server into buffer ( byte slice )
    n, err := conn.Read( buffer[:] )
    check( err )    // check for errors

    // split recv'd message over @, separate header from message
    msg := strings.Split( string( buffer[:n] ), "@" )

    var resp string     // declare response variable
    switch msg[0] {     // switch/case over the header
        case "la":      // /la/, list answer
            if len( msg ) > 2 {
                resp = "lr@"

            } else {
                // split recv'd message over newline to create quote list
                quotes := strings.Split( msg[1], "\n" )
                for p, v := range quotes {      // iterate over quotes
                    // print num of quote, followed by that quote
                    fmt.Println( p + 1, "\t", v )
                }
                resp = "te@success"
            }
        default:        // something undefined entirely
            // TODO write out to file
            //fmt.Println( "Message header not recognized" )
            //fmt.Println( msg )
            conn.Close()    // TODO how important is this line
            return      // dead return, essentially lets program quit
    }

    // send formulated response back to server
    _, err = conn.Write( []byte( resp ) )
    check( err )    // make sure write worked

    handle( conn )  // call handle again
}
*/


// quick check for errors
func check( err error ) {
    if err != nil {     // err should be nil, !nil is bad
        // print error and ditch
        fmt.Println( "fatal::", err )
        os.Exit( 1 )
    }
}


// level 4 handle
func handle( conn net.Conn ) {
    // create byte array of pre-ordained size
    buf := make( []byte, tsize )

    // read from connection into byte array
    _, err := conn.Read( buf )
    check( err )    // routine check for errors

    msg := string( buf )    // convert to string

/*
    // iterate over message string
    for p, c := range msg {
        // cast char ( rune ) to string and see if it's the end char
        if string( c ) == "|" {
            // if end char, print everything BUT char and exit
            fmt.Println( msg[:p] )
            return
        }
    }
*/

    // iterate over message string
    for _, c := range msg {
        // cast char ( rune ) to string and see if it's the end char
        if string( c ) == "X" {
            // if end char, print everything BUT char and exit
            fmt.Println( msg )
            return
        }
    }

    fmt.Println( msg )      // if char is not found, print entire
    handle( conn )          // message and recursively call level4
}



// globals
var tsize int
var dest string     // provide a new destination ip and port


func main() {

    // define flags
    flag.IntVar( &tsize, "s", 10, "tx size" )
    flag.StringVar( &dest, "ip", ":1300", "ip:port of server" )
    flag.Parse()    // process flags

    service := dest     // declare server ip:port

    // create address object
    addr, err := net.ResolveTCPAddr( "tcp", service )
    check( err )    // check for errors

    // make connection with object
    conn, err := net.DialTCP( "tcp", nil, addr )
    check( err )    // check for errors

    //command := "lr@pls"

    // write command to server, beginning transaction
    //_, err = conn.Write( []byte( command ) )
    //check( err )    // make sure write worked

    handle( conn )  // handle response
    // TODO stack trace when handle returns, should probably be
    // updated
    conn.Close()    // close connection when finished

    os.Exit( 0 )    // exit

}

