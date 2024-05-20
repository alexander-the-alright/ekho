// =============================================================================
// Auth: Alex Celani
// File: test.go
// Revn: 05-19-2024  0.4
// Func: explore bounding on piecewise file transmission ( PWFT )
//
// TODO: implement
// =============================================================================
// CHANGE LOG
// -----------------------------------------------------------------------------
// 04-19-2024: bones
//             level 0 PWFT
//             level 1 PWFT
//             level 2 PWFT
// 05-13-2024: comment level 1 and level 2 PWFT
//             level 3 PWFT
// 05-18-2024: began level 4 changes
// 05-19-2024: got level 4 working
//             commented
//             wrote level 5
//             commented
//             wrote level 6!
//             commented
//
// =============================================================================

package main

import ( 
    "flag"
    "fmt"
    "math/rand"
    "net"       // ResovleTCPAddr, ListenTCP, listener.Accept
                // conn.Read,Write
    "os"        // Exit
    "strconv"
)


/*
func makedata() {
    for c := 0; c < dsize; c++ {
        datum[c] = rand.Intn( dsize )
    }
}
*/

// special make function for strings
func makedatatx() {
    // classic iteration of length dsize, toss in same-sized array
    for c := 0; c < dsize; c++ {
        // toss in same size array
        // rand int picks letter, 0x61 -> lowercase ascii, cast byte
        datumtx[c] = byte( rand.Intn( 26 ) + 0x61 )
    }
    // cast to string and print ( as slice )
    fmt.Println( string( datumtx[:] ) )
}


// handle errors catastrophically
func check( err error ) {
    if err != nil {    // error is not nil on error
        // print error
        fmt.Println( "Fatal: ", err.Error() )
        os.Exit( 1 )    // bail
    }
}


// Levels of PWFT
// 0. print everything indescriminately
// 1. printing piecewise
// 2. printing piecewise with header
// 3. printing piecewise with variable size header
// 4. file tx piecewise
// 5. file tx piecewise with header
// 6. file tx piecewise with variable header


// level 0 PWFT
func level0( data []int) {
    fmt.Println( data )    // level 0 sux
}


// level 1 PWFT
func level1() {
    tsize := 5               // print 5 chars at a time
    txn := dsize / tsize    // amount of prints needed
    for c := 0; c < txn; c++ {      // iterate that amount of times
        // print ( using wrapper ) n-1 data points of size tsize
        level0( datum[c*tsize:(c+1)*tsize] )
    }
    level0( datum[txn*tsize:] )     // print nth data point
}


// level 2 PWFT
func level2() {
    var bsize int = 10              // total transmission size
    var header string = "Lr@x@"     // static header
    // calculate usable header size
    var tsize int = bsize - len( header )
    txn := dsize / tsize            // amount of prints needed
    for c := 0; c < txn; c++ {      // iterate that amount of times
        fmt.Print( header )         // print static header
        // exact same print statement from level 1
        level0( datum[c*tsize:(c+1)*tsize] )
    }
    fmt.Print( header )             // print nth static header
    level0( datum[txn*tsize:] )     // print nth data point again
}


// level 3 PWFT
func level3() {
    var header string       // declare header
    var bsize int = 10      // total transmission size
    var hhead, htail string = "Lr@", "@"    // static pieces of header
    txn := 1    // init transmission number
    pos := 0    // init position tracker

    for {   // easier to do an infinite loop and break later
        // combine both static parts of the header with the tx number
        header = hhead + strconv.Itoa( txn ) + htail
        hlen := len( header )   // keep track of total header length
        fmt.Print( header )     // print header

        // level0 print ( bsize - hlen ) units of data, starting at
        // the noted position
        level0( datum[ pos:(pos + bsize - hlen) ] )
        // update position as last position printed
        pos = pos + bsize - hlen
        txn++   // update transmission number

        // check to see if *updated* position goes out of bounds
        if pos + bsize - hlen > dsize {
            break   // break if'n
        }
    }
    // print nth header ( X denotes final transmission )
    fmt.Print( hhead + "X" + htail )
    level0( datum[pos:] )   // print from final position to end
}


// level 4 PWFT
func level4( conn net.Conn ) {
    txn := dsize / tsize    // amount of transmissions

    // iterate for number of transmissions
    for c := 0; c < txn; c++ {
        // write that amount of times, very specific slice
        _, err := conn.Write( datumtx[c*tsize:(c+1)*tsize] )
        check( err )    // routine error check
    }

    // write final of last chunk of data
    _, err := conn.Write( datumtx[txn * tsize:] )
    check( err )    // routine error check
    _, err = conn.Write( []byte( "|" ) )    // write end char
    check( err )    // routine error check
}


// level 5 PWFT
func level5( conn net.Conn ) {
    header := "la@x@"               // static header
    psize := tsize - len( header )  // calculate payload size

    txn := dsize / psize    // amount of transmissions

    // iterate for number of transmissions
    for c := 0; c < txn; c++ {
        // concat header with payload, easier as string, cast to byte
        packet := header + string( datumtx[c*psize:(c+1)*psize] )
        // write that amount of times, very specific slice
        _, err := conn.Write( []byte( packet ) )
        check( err )    // routine error check
    }

    // concat last payload
    packet := header + string( datumtx[txn*psize:] )
    // write final of last chunk of data
    _, err := conn.Write( []byte( packet ) )
    check( err )    // routine error check
    _, err = conn.Write( []byte( "|" ) )    // write end char
    check( err )    // routine error check
}


// level 6 PWFT
func level6( conn net.Conn ) {
    var header string       // declare header
    var hhead, htail string = "la@", "@"    // static pieces of header
    txn := 1    // init transmission number
    pos := 0    // init position tracker

    for {   // easier to do an infinite loop and break later
        // combine both static parts of the header with the tx number
        header = hhead + strconv.Itoa( txn ) + htail
        hlen := len( header )   // keep track of total header length

        // craft packet
        packet := header + string( datumtx[pos:( pos + tsize - hlen )] )

        // cast packet to byte array and send
        _, err := conn.Write( []byte( packet ) )
        check( err )    // routine error check

        // update position as last position printed
        pos = pos + tsize - hlen
        txn++   // update transmission number

        // check to see if *updated* position goes out of bounds
        if pos + tsize - hlen > dsize {
            break   // break if'n
        }
    }
    // print nth header ( X denotes final transmission )
    header = hhead + "X" + htail
    packet := header + string( datumtx[pos:] )
    // cast data to string, concat, cast to byte array, send
    _, err := conn.Write( []byte( packet ) )
    check( err )    // routine error check

}


// size variables
var dsize int = 104
// data global
var datum = make( []int, dsize )
var datumtx = make( []byte, dsize )
var tsize int


func main() {

    flag.IntVar( &tsize, "s", 10, "tx size" )
    flag.Parse()

    makedatatx()  // fill global with data

    service := ":1300"  // create service on ip and port

    // resolve ip address and port
    addr, err := net.ResolveTCPAddr( "tcp", service )
    check( err )    // make sure ip resolves

    // create listener object from ip:port
    listener, err := net.ListenTCP( "tcp", addr )
    
    // print global data to show that PWFT works as intended
    //level0( datum )

    //level3()    // do the thing

    for {
        // wait, create connection when found
        conn, err := listener.Accept()
        check( err )    // make sure connection works
    
        go level6( conn )  // handle connection
    }
}

