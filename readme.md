 <div align="center">
    <img src="./yylogo.png">
</div>

YeetYoink (YY for short) is a dynamically typed programming language that combines both functional and imperative programming paradigms. 


# Key features

- **Keywords start with 'y'**. With all keywords starting with a big, bold 'y', you'll never forget what language you're coding in.
- **YOLO mode**. When you're ready to let loose and have some fun, switch on YOLO mode and ignore those pesky compiler errors. YY will try its best to execute the code, but results may vary. Play stupid games, win stupid prizes.
- **Everything is an expression.** Even statements like control flow constructs (yif, yall, etc.) and yassignments yeet a value. This allows for a more concise and expressive code - peak efficiency unlocked.
- **First-class functions.** Functions are first-class citizens, meaning they can be passed around like hot potatoes, yeeted as values, and stored in data structures like arrays and hashmaps.
- **Closures.** Functions capture variables from their surrounding scope, making them more powerful and flexible than your above average yak.
- **REPL**. YY comes with a REPL (Read-Eval-Print Loop), allowing you to try out code snippets, experiment, and test your functions on the fly.


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
yif 2 * 2 > 1 {
   "all good" 
} yels yif 8 + 8 < 4 {
   "still all good" 
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
sum    := 0
yall yarray {
    sum += yt
}

yelp(sum) // prints 6

// 2nd type of loop: 'yet' as in 'are we there yet?'
// similar to a 'while' loop in other languages
i := 0
yet i < 5 {
    i += 1
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
    42:        "yes, definitely",
    "colours": ["brown", "white"], 
}

yelp(friend["name"], "is", friend["age"], "years old.") // prints "Jon the Yak is 2 years old."
yelp(friend[42]) // prints "yes, definitely"

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

// assignment to a variable that hasn't been declared would normally error:
new_var = 5 // runtime error: 'identifier not found: new_var'

// but in Yolo Mode, anything goes:

yolo {
    yelp("tree" * 18)   // prints forest
    yelp("2" * 5)       // prints 10
    yelp("troll" * 3)   // prints trolltrolltroll
    yelp([1, 2, 3] * 3) // prints [3, 6, 9]

    new_var = 5 // this works even though new_var hasn't been declared first

    yelp("oh boy" / 0)  // see for yourself...
}

```

# Usage

Build with 

```
go build
```

Run a script

```
$ ./yy [--debug] filename
```

Or start a REPL session 

```
$ ./yy [--debug]
```

Note: `--debug` flag is optional

# Other, less exciting, features

- **Two data structures.** YY supports arrays and hashmaps, providing more data structure options than Lua.
- **Very basic data types**. YY supports the basic data types of yinteger, string, bool and null. And yes, null isn't technically a data type.
- **Optional semicolons.** YY has taken the modern approach of making semicolons optional, allowing for a cleaner codebase (semicolons are so 1970s anyway). 
- **Garbage collected**. YY uses automated memory management, meaning you don't have to worry about freeing up memory that is no longer being used. It's like having a personal janitor for your code!
- **Not Object-Oriented**. You don't have to wrap your head around inheritance hierarchy if there's no inheritance hierarchy.
- **Built-in functions.** YY includes a number of built-in functions for common tasks.
- **Pass by value**. Variables are passed by value, meaning that when a variable is passed as an argument to a function, a copy of its value is used in the function. No more worrying about your function making a mess of your precious data.


# FAQ

### Q: Why is the project called YeetYoink?

Yeet and Yoink symbolise two complementary, yet opposing forces that exist in the universe. Everything has both Yeet and Yoink aspects and they are interconnected and interdependent. Together, Yeet and Yoink form a whole, and the balance between the two is necessary for harmony and balance in the universe.

### Q: What's the deal with the fishies in the logo?

A: They are two koi fish and they are the mascot for the project. The dark one is called Yeet and is passive, negative, and feminine. The bright one is called Yoink and is active, positive, and masculine.

### Q: Can this language be used for serious business?

Of course, why do you ask.

### Q: Has anyone actually asked these questions?

A: No, I made this section up.
