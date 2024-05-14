// =============================================================================
// Auth: alex
// File: client.go
// Revn: 10-05-2023  2.0
// Func: ask server for a message of the day quote
//
// TODO: document
//       expand on /ba/ failures
//       add more flags?
//       log errors in logfile
// =============================================================================
// CHANGE LOG
// -----------------------------------------------------------------------------
// 05-11-2022: init
// 05-17-2022: changed an entry in the response table
// 08-16-2023: copied from speak.go
//             gutted, rewrote
// 09-12-2023: added size conversation
//             added -d, -l, and -dest flags
//             commented
// 09-21-2023: began complete-ass overhaul of handle() system
//             wrote draft of draw()
// 09-26-2023: made draw return size and reference globals
//             imported check() from soary
// 09-29-2023: commented handle()
//*10-02-2023: commented main(), and friends
//*10-05-2023: /remove/ works
//             removed -d flag
//             commented
//
// =============================================================================

package main

import ( 
    "io"        // io.EOF
    "flag"      // BoolVar, StringVar, Parse
    "fmt"       // Println
    "math/rand" // NewSource, New, Intn, *rand.Rand
    "net"       // ResolveTCPAddr, DialTCP, conn.Write,Read
    "os"        // Exit
    "strconv"   // Atoi, Itoa
    "strings"   // Split
    "time"      // Now, UnixNano
)


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


// get user input
func in( prompt string ) string {
    fmt.Print( prompt )     // print user prompt
    var input string        // declare var to hold user input
    fmt.Scanln( &input )    // get user input
    // try to convert to an int, discard results
    _, err := strconv.Atoi( input )
    // if there was an error, user didn't input an int
    if err != nil {
        // warn user against this
        fmt.Println( "Input must be a number" )
        input = in( prompt )    // try to get user input again
        // this, like handle() calling handle(), seems like a super
        // bad idea, stack-wise...
    }
    // once user has properly inputted a number, return string of num
    return input
}


// seed and get a new rng
func seed() *rand.Rand {
    // use current time to seed random number
    s := rand.NewSource( time.Now().UnixNano() )
    r := rand.New( s )  // create new random number from seed
    return r            // return random number
}


// get a random number
func draw() string {
    rn := seed()            // get a random seed
    r := rn.Intn( qsize )   // get a random number using seed
    // convert random number to ascii, because converting
    // ints to bytes ( for networking ) and back is obnoxious
    rsz := strconv.Itoa( r )
    return rsz      // return rand size
}


// get server response and print
func handle( conn net.Conn ) {

    var buffer [512]byte    // create buffer to hold response
    // read n bytes from server into buffer ( byte slice )
    n, err := conn.Read( buffer[:] )
    check( err )    // check for errors

    // split recv'd message over @, separate header from message
    msg := strings.Split( string( buffer[:n] ), "@" )

    var resp string     // declare response variable
    switch msg[0] {     // switch/case over the header
        // if message is...
        case "sa":      // /sa/, size answer
            // convert back half to int, catch error
            size, nerr := strconv.Atoi( msg[1] )
            if nerr == nil {    // if there was no error
                qsize = size    // set size global to recv'd size
            } else {
                // if int conversion failed, just quit
                // TODO make more robust
                os.Exit( 1 )
            }
            req := draw()           // draw a random number
            resp = "qr@" + req      // create quote query with rng
        case "ba":      // /ba/, bad answer
            // TODO create server-side, and handle client-side
            fmt.Println( "error in " + msg[1] + ": " + msg[2] )
            resp = "te@error"   // set transaction end message
        case "qa":      // /qa/, quote answer
            fmt.Println( msg[1] )   // whole-ass just print quote
            resp = "te@success"     // set transaction end message
        case "la":      // /la/, list answer
            // split recv'd message over newline to create quote list
            quotes := strings.Split( msg[1], "\n" )
            for p, v := range quotes {      // iterate over quotes
                // print num of quote, followed by that quote
                fmt.Println( p + 1, "\t", v )
            }
            resp = "te@success"     // set transaction end message
        case "aa":      // /aa/, add answer
            // if server says add succeeded
            if msg[1] == "success" {
                fmt.Println( "Add successful" )     // print success
                resp = "te@success" // set transaction end message
            // if server says add failed
            } else if msg[1] == "error" {
                fmt.Println( "Add failed" )     // print fail
                resp = "te@error"   // set transaction end message
            // unknown status message
            } else {
                // print unknown status
                fmt.Println( "Unknown status code: " + msg[1] )
                // TODO should this be unknown instead of error?
                resp = "te@error"   // set transaction end message
            }
        case "r1":      // /r1/, remove ( first transaction )
        // the idea is that if the user doesn't specify which quote to
        // remove, list all quotes and let the user choose the quote
        // to remove. There was probably a way to handle this a little
        // more simply, but this, I feel, has some nuance to it that I
        // actually like
            // redo /la/, list all quotes
            // split recv'd message over newline to create quote list
            quotes := strings.Split( msg[1], "\n" )
            for p, v := range quotes {      // iterate over quotes
                // print num of quote, followed by that quote
                fmt.Println( p + 1, "\t", v )
            }
            // ask user which quote to remove
            resp = "r2@" + in( "\n-> " )
        case "r2":      // /r2/, remove ( second transaction )
            // if server says remove succeeded
            if msg[1] == "success" {
                fmt.Println( "Remove successful\n" )    // success
                // list all quotes to validate remove succeeded
                // TODO should /aa/ also make /lr/ before /te/?
                resp = "lr@dumvar"  // set transaction end message
            // if server says remove failed
            } else if msg[1] == "error" {
                fmt.Println( "Remove failed" )      // print failure
                resp = "te@error"   // set transaction end message
            // if unknown status message
            } else {
                // print unknown status
                fmt.Println( "Unknown error code: " + msg[1] )
                // TODO should this be unknown instead of error?
                resp = "te@error"   // set transaction end message
            }
        case "br":      // /br/, bad request
        // client has asked for something nonsensical
            // TODO write out to file
            fmt.Println( "Client sent a bad request" )
            resp = "te@error"       // set transaction end message
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


// globals
var qsize int       // keep track of num of quotes
var conn net.Conn   // keep track of connection


// global flag variables
var add bool        // add new quotes to quote file ( on server )
var remove bool     // remove quotes from quote file ( on server )
var list bool       // request entire quote file from server
var dest string     // provide a new destination ip and port


func main() {

    // define flags
    flag.BoolVar( &add, "a", false, "add quote to list" )
    flag.BoolVar( &remove, "r", false, "remove quote from list" )
    flag.BoolVar( &list, "l", false, "print quote list" )
    flag.StringVar( &dest, "ip", ":1300", "ip:port of server" )
    flag.Parse()    // process flags

    service := dest     // declare server ip:port

    // create address object
    addr, err := net.ResolveTCPAddr( "tcp", service )
    check( err )    // check for errors

    // make connection with object
    conn, err = net.DialTCP( "tcp", nil, addr )
    check( err )    // check for errors

    // block for all flags
    var command string
    switch {        // empty switch statement is basically an if
        case list:  // if list flag is active
            command = "lr@dumvar"   // list request
            // TODO come up with a useful argument, instead of dumvar
        case add:   // if add flag is active
            // if arguments were specified
            if len( flag.Args() ) != 0 {
                command = "ar@"     // preface command w/ add request
                // need to concat flag.args
                for _, v := range flag.Args() {
                    // add the string and a space to command
                    command = command + v + " "
                }
                // remove the trailing space
                command = command[:len( command ) - 1]
            } else {    // if no arguments were specified
                // print that you need args and bail
                fmt.Println( "/a/ flag must be used with an arugment" )
                os.Exit( 1 )
            }
        case remove:    // if remove flag is active
            // if arguments were specified
            if len( flag.Args() ) > 0 {
                // argument is supposed to the number of the quote to
                // be removed, get first arg
                // cast to int just to make sure it's a valid number
                // assuming the user knows what they're doing
                _, rerr := strconv.Atoi( flag.Args()[0] )
                // if there was an error, just pretend there's no args
                if rerr != nil {
                    command = "r1@dumvar"   // remove request 1
                    // TODO come up w/ useful argument, replace dumvar
                } else {
                    // if user knows the number, take care of it
                    // arg[0] is a string, so concat is cool
                    // jump right into remove request 2
                    command = "r2@" + flag.Args()[0]
                }
            } else {    // if no argument was specified for remove
                // jump right into remove request 1
                command = "r1@dumvar"
            }
        default:    // if no flag was specified
            // begin transaction for quote request
            command = "sr@please"   // size request
    }

    // write command to server, beginning transaction
    _, err = conn.Write( []byte( command ) )
    check( err )    // make sure write worked

    handle( conn )  // handle response
    // TODO stack trace when handle returns, should probably be
    // updated
    conn.Close()    // close connection when finished

    os.Exit( 0 )    // exit

}

