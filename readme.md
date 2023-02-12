 <div align="center">
    <img src="./yylogo.png">
</div>

YeetYoink (YY for short) is a dynamically typed programming language that combines both functional and imperative programming paradigms. 


# Features

- **Keywords start with 'y'**: with all keywords starting with a big, bold 'y', you'll never forget what language you're coding in.
- **YOLO mode**: when you're ready to let loose and have some fun, switch on YOLO mode and ignore those pesky compiler errors. YY will try its best to execute the code, but results may vary. Play stupid games, win stupid prizes.
- **Everything is an expression:** even statements like control flow constructs (yif, yall, etc.) and yassignments yeet a value. This allows for a more concise and expressive code - peak efficiency unlocked.
- **First-class functions:** functions are first-class citizens, meaning they can be passed around like hot potatoes, yeeted as values, and stored in data structures like arrays and hashmaps.
- **Closures:** functions capture variables from their surrounding scope, making them more powerful and flexible than your above average yak.
- **REPL**: YY comes with a REPL (Read-Eval-Print Loop), allowing you to try out code snippets, experiment, and test your functions on the fly.


# Quick tour

For more details, check out [examples](examples) dir.

```c
a := 5 // variables are declared and assigned using walrus operator ':='
a = 10 // variable assignment, variables must be declared before use

// supported types: integer, string, bool, null

yinteger := 5
yarn     := "how long is a piece of string?"
yup      := true
ylem     := null

// print with 'yelp()' (or 'yell()' if urgent)
yelp("Hello, world!")       // prints "Hello, world!"
yell("Hello, cruel world!") // prints "HELLO, CRUEL WORLD!"

// CONTROL FLOW

// yif requires brackets, but doesn't require parentheses
yif 2*2 > 1 {
   "all good" 
} yels {
    "math.exe stopped working"
}

// there are 2 looping constructs in YY: yall and yet

// yall (Y'all) yeeterates over a collection (array, string or range)
// variable 'yt' (short for yeeterator) is created automatically in the loop's scope
yall 0..3 {
    yelp(yt) // prints '0', '1', '2', '3'
}

yarray := [1, 2, 3]
sum := 0
yall yarray {
    sum = sum + yt
}
yelp(sum) // prints 6

// 2nd type of loop: 'yet' as in 'are we there yet?'
// similar to a 'while' loop in other languages
i := 0
yet i < 5 {
    i = i + 1
}

// ARRAYS

my_yarray := [1, 2, 3, 4]

// yarray can hold values of different types
my_yarray2 := [1, true, "hello", yarn, yup]

// yarrays can be concatenated
mega_yarray := my_yarray + my_yarray2

// HASHMAPS

yak := { 
    "name":    "Jon the Yak", 
    "age":     2, 
    "alive":   true,
    "colours": ["brown", "white"], 
}
yelp(friend["name"], "is", friend["age"], "years old.") // prints "Jon the Yak is 2 years old."

// a key to a hashmap can be a string, an integer, or a bool
foe := { 
    "name": "Giant Chicken", 
    true:   "bool as a key!?",
    42:     "an integer works too",
}

// FUNCTIONS

// anonymous functions (lambdas) are declared using '\'
// (if you squint hard enough, '\' looks kinda like lambda 'λ')

max := \x y {
    yif x > y {
        yeet x // yeets the value and returns from the function
    } 
    yeet y
}

// 'yeet' keyword can be omitted, last statement in a block is yeeted implicitly
max2 := \x y {
    yif x > y {
        x 
    } yels {
        y
    }
}

// recursion

factorial := \n { 
    yif n == 0 { 1 } yels { n * factorial(n-1) }
}

// higher-order functions

add_three      := \x { x + 3 }
call_two_times := \x fn { fn(fn(x)) }
yelp(call_two_times(5)) // prints 11

// closures

new_adder := \x { 
    \n { x + n } 
}
add_two := new_adder(2)
yelp(add_two(5)) // prints 7

// YOLO MODE

// mixing types would normally cause an error:
c := 8 * "string" // runtime error: 'type mismatch: STRING * INTEGER'
// But in Yolo Mode, anything goes:

yolo {
    yelp("2" * 5)       // prints 10
    yelp("troll" * 3)   // prints trolltrolltroll
    yelp([1, 2, 3] * 3) // prints [3, 6, 9]
    yelp("oh boy" / 0)  // see for yourself...
}

```

# Usage

To run a script:

```
$ yy filename
```

To start a REPL session (--debug flag is optional):

```
$ yy [--debug]
```