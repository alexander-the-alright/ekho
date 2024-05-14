// =============================================================================
// Auth: Alex Celani
// File: test.go
// Revn: 05-13-2024  0.2
// Func: explore bounding on piecewise file transmission ( PWFT )
//
// TODO: rewrite makedata() to write to a file
//       write level4()
// =============================================================================
// CHANGE LOG
// -----------------------------------------------------------------------------
// 04-19-2024: bones
//             level 0 PWFT
//             level 1 PWFT
//             level 2 PWFT
// 05-13-2024: comment level 1 and level 2 PWFT
//             level 3 PWFT
//
// =============================================================================

package main

import ( 
    "fmt"
    "math/rand"
    "strconv"
)

func makedata() {
    for c := 0; c < dsize; c++ {
        datum[c] = rand.Intn( dsize )
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



// size variables
var dsize int = 104
// data global
var datum = make( []int, dsize )


func main() {
    makedata()  // fill global with data

    // print global data to show that PWFT works as intended
    level0( datum )

    level3()    // do the thing
}

