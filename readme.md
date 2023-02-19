 <div align="center">
    <img src="./yylogo.png">
</div>

YeetYoink (YY for short) is a dynamically typed programming language that combines functional and imperative programming paradigms.


# Key features

- **Keywords start with 'y'**. With all keywords starting with a big, bold 'y', you'll never forget what language you're coding in.
- **YOLO mode**. When you're ready to let loose and have some fun, switch on YOLO mode and ignore those pesky compiler errors. YY will try its best to execute the code, but results may vary. Play stupid games, win stupid prizes.
- **Everything is an expression.** Even statements like control flow constructs (yif, yall, etc.) and yassignments yeet a value. This allows for a more concise and expressive code - peak efficiency unlocked.
- **First-class functions.** Functions are first-class citizens, meaning they can be passed around like hot potatoes, yeeted as values, and stored in data structures like arrays and hashmaps.
- **Closures.** Functions capture variables from their surrounding scope, making them more powerful and flexible than your above average yak.
- **REPL**. YY comes with a REPL (Read-Eval-Print Loop), allowing you to try out code snippets, experiment, and test your functions on the fly.


# Quick tour

For more details, check out [examples](examples) directory.

## Hello world

```c
// print with 'yelp()' (or 'yell()' if urgent)
name := "Yennefer"
yelp("Hello, {name}!") // prints "Hello, Yennefer!"
yell("Hello, {name}!") // prints "HELLO, YENNEFER!"
```

## Variables

```c
a := 5 // variables are declared and assigned using walrus operator ':='
a = 10 // variable assignment, variables must be declared before use

// supported types: integer, string, bool, null
my_yinteger := 5
my_yarn     := "how long is a piece of string?"
my_yup      := true
my_void     := null
```

## Control flow

```c
yif 2 * 2 > 10 {
   "that's untrue"
} yels yif 8 + 8 < 4 {
   "yup, that's our stop"
} yels {
    "math.exe stopped working"
}
```

## Loops

```c
// yall (y'all) yeeterates over a collection
// variable 'yt' (short for yeeterator) is created automatically in the loop's scope
yall 0..3 {
    yelp(yt) // prints '0', '1', '2', '3'
}

yarray := [1, 2, 3]
sum    := 0
yall elt: yarray { // optionally, you can name the yeeterator
    sum += elt
}

yelp(sum) // prints 6
```


```c
// 2nd type of loop: 'yet' as in 'are we there yet?'
// similar to a 'while' loop in other languages
i := 0
yet i < 5 {
    i += 1
}
```

## Data structures

```c
my_yarray := [1, 2, 3, 4]

// yarray can hold values of different types
my_yarray2 := [1, true, "hello"]

// yarrays can be concatenated
mega_yarray := my_yarray + my_yarray2
```

```c
my_hashmap := {
    "name":  "Yakub the Yak",
    "age":   2,
    "alive": true,
    42:      "yes, definitely",
}

yelp(my_hashmap["name"], "is", my_hashmap["age"], "years old.") // prints "Yakub the Yak is 2 years old."
```

## Functions

```c
// anonymous functions (lambdas) are declared using '\'
// (if you squint hard enough, '\' looks kinda like a lambda 'λ')
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
```

```c
// recursion
factorial := \n {
    yif n == 0 {
        1
    } yels {
        n * factorial(n-1)
    }
}

// higher-order functions
add_three  := \x { x + 3 }
call_twice := \x fn { fn(fn(x)) }
yelp(call_twice(5, add_three)) // prints 11

// closures
new_adder := \x {
    \n { x + n }
}
add_two := new_adder(2)
yelp(add_two(5)) // prints 7
```

## Yolo Mode

```c
// mixing types would normally cause an error:
c := 8 * "string" // runtime error: 'type mismatch: STRING * INTEGER'

// assignment to a variable that hasn't been declared would normally error:
new_var = 5 // runtime error: 'identifier not found: new_var'

// but in Yolo Mode, anything goes:

yolo {
    yelp("tree" * 18)   // prints 'forest'
    yelp("2" * 5)       // prints 10
    yelp("troll" * 3)   // prints 'trolltrolltroll'
    yelp([1, 2, 3] * 3) // prints [3, 6, 9]

    // this works even though new_var hasn't been declared first
    new_var = 5

    // but even in yolo mode, division by zero doesn't end well
    yelp("weee" / 0) // prints "Stare at the abyss long enough, and it starts to stare back at you."
}
```


# Usage

Build with

```
$ go build
```

Run a YY script

```
$ ./yy [--debug] filename
```

Or start a REPL session

```
$ ./yy [--debug]
```

Note: `--debug` flag is optional


# More features

- **Two data structures.** YY supports arrays and hashmaps, providing more data structure options than Lua.
- **Very basic data types**. YY supports the basic data types of yinteger, string, bool and null. And yes, null isn't technically a data type.
- **Optional semicolons.** YY has taken the modern approach of making semicolons optional, allowing for a cleaner codebase (since semicolons are so 1970s).
- **Garbage collected**. YY uses automated memory management, meaning you don't have to worry about freeing up memory that is no longer being used. It's like having a personal janitor for your code!
- **Not Object-Oriented**. You don't have to wrap your head around inheritance hierarchy if there's no inheritance hierarchy.
- **No exception handling**. No more wrangling with complex error handling mechanisms. In YY, you can throw an exception, but there is no mechanism for catching it (we're not half-assing it like Go, with its weird panic-recover mechanism).
- **Built-in functions.** YY includes a number of built-in functions for common tasks.


# FAQ

## Why is the project called YeetYoink?

Yeet and Yoink symbolise two complementary, yet opposing forces that exist in the universe. Everything has both Yeet and Yoink aspects and they are interconnected and interdependent. Together, Yeet and Yoink form a whole, and the balance between the two is necessary for harmony and balance in the universe.

## What's the deal with the fishies in the logo?

They are two koi fish and they are the mascot for the project. The dark one is called Yeet and is passive, negative, and feminine. The bright one is called Yoink and is active, positive, and masculine.

## Can this language be used for serious business?

Of course, why do you ask.

## Has anyone actually asked these questions?

No, I made this section up.
