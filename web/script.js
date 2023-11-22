function run() {
    if (!source.value) return

    output.innerText = ''

    const result = interpret(source.value)
    if (result?.error) {
        if (output.innerText) output.innerText += '\n'
        output.innerText += result.error
    }
}

function captureLog(msg) {
    const li = document.createElement('li')
    li.innerText = msg
    output.appendChild(li)
}

window.console.log = captureLog

// for better UX, we listen for ctrl+enter on the main window rather than textbox only
window.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && e.ctrlKey) run()
})

source.addEventListener('keydown', (e) => {
    if (e.key === 'Tab') {
        document.execCommand('insertText', false, '    ') // tab is 4 spaces
        e.preventDefault() // prevent tabbing out from textarea
    }
})

function buildSampleSelect() {
    const sampleSelect = document.querySelector('#sample-select')

    for (const sample in samples) {
        const option = document.createElement('option')
        option.value = sample
        option.textContent = sample
        sampleSelect.appendChild(option)
    }

    sampleSelect.addEventListener('change', (e) => setSample(e.target.value))
}

function setSample(sampleName) {
    source.value = samples[sampleName]
    source.scrollTo(0, 0)
}

const samples = {
    'hello world': `// You can edit this code, or select a sample from the dropdown on the right.
//
// To run the code press Ctrl+Enter or click 'Run' button.

name := "Yennefer"
yap("Yo, {name}!")`,

    fizzbuzz: `// Ah, FizzBuzz, the timeless test that weeds out the 10x engineers from the wannabes in programming
// interviews. But fear not, YY is here to help you slay this beast. And rather than printing
// the mundane FizzBuzz, we'll print out the magnificent YeetYoink instead. 

// This example illustartes the use of implicitly defined 'yt' variable.
// 'yt' is short for 'yeeterator' and it's created autmatically inside a yall loop.

yall 1..100 {
    result := ""

    yif yt % 3 == 0 { result += "Yeet" }  // append "Yeet" if the current number is divisible by 3
    yif yt % 5 == 0 { result += "Yoink" } // append "Yoink" if the current number is divisible by 5

    yif result {    // check if the result isn't empty 
        yap(result) // if so, print the result (Yeet, Yoink, or YeetYoink)
    } yels {
        yap(yt)     // print the current number instead
    }
}
`,

    fibonacci: `// Implementation of Fibbonacci numbers using two ways: recursion and closure. Just like choosing
// between pizza and spaghetti, there is no right or wrong way to do it, both are equally satisfying.
// And while these methods may not be the fastest, they add some spicy flavor to this demo.

// Recursion
fib := \\n {
    yif n < 2 { 
        n                   // if n is less than 2, return n (base case)
    } yels { 
        fib(n-1) + fib(n-2) // otherwise, calculate Fibonacci recursively
    }
}

yap("seventh Fibonacci number:", fib(7))


// Closure
fib_gen := \\{
    a := 0
    b := 1
    // return a closure that will return next Fibonnacci number
    \\{
        tmp := a
        a = b
        b += tmp
        a
    }
}

f := fib_gen() // create a closure instance 'f' for generating Fibonacci numbers
yap("consecutive Fibonacci numbers:", f(), f(), f(), f(), f())
`,

    'map et al': `// Map, filter, and reduce are The Three Musketeers of functional programming, banding together
// to process and transform collections with finesse and style.

// Map transforms all elements and returns a shiny new list.
map := \\arr, fn {
    acc := []
    yall arr {
        acc << fn(yt)
    }
}

// Reduce violently smashes a list into a single value.
reduce := \\arr, initial, fn {
    result := initial
    yall arr {
        result = fn(result, yt)
    }
}

// Filter picks the juiciest elements, leaving the rest to wither away into obscurity.
filter := \\arr, fn {
    acc := []
    yall arr {
        yif fn(yt) {
            acc << yt
        }
    }
    acc
}

arr := [1, 2, 3, 4]
yap("original array:", arr)

tripled := map(arr, \\x { x * 3 })
yap("tripled:", tripled)

summed := reduce(tripled, 0, \\x, y { x + y })
yap("sum of elements:", summed)

avg  := summed / len(arr)
smol := \\x { x < avg }
yap("smaller than average:", filter(tripled, smol))
`,

    yolo: `// Yolo Mode allows you to do things that would make other self-respecting languages blush.
// (Not JavaScript though. JavaScript wouldn't bat an eye.)
//
// Types can be mismatched, strings can be negated, variables don't have to be declared before use.
// But be warned, the return value is anyone's guess. What about the Principle of Least Surprise you ask?
// Exactly, what about it?
//
// Go ahead and experiment, but remember Uncle Ben's words of wisdom: "Play stupid games, win stupid prizes".

yolo {
    // you can multiply a string by an integer
    yap("'tree' * 18 =", "tree" * 18)
    yap("'troll' * 3 =", "troll" * 3)
    yap("'2' * 5 =", "2" * 5)

    // you can multiply an array by an integer
    yap("[1, 2, 3] * 3 = ", [1, 2, 3] * 3)

    // you can negate a string
    yap("-'i am a string' =", -"i am a string")

    // you can do useful stuff too, like bake an argument into a function
    // (check out 'bake' example for more details)
    greet     := \\name { yap("Hello, {name}!") }
    greet_yan := greet + "Yan"
    greet_yan() // look ma, no args!

    // but even in yolo mode, division by zero doesn't end well (what did you expect?)
    yap("division by zero:", "weee" / 0)
}`,

    bake: `// Brace yourselves, we're about to go into YOLO mode! We'll be adding numbers, arrays, and hashmaps
// to a function like a mad scientist adding ingredients to a cauldron. This magically bakes the
// arguments into the function, turning it into a deliciously self-contained recipe for success.
//
// This is a powerful technique that can make your code more concise and easier to read, especially
// when you have functions with many arguments that are frequently used with certain fixed values.
//
// Some fancy folks call this 'partial function application' or 'currying', we'll just call it baking.

// Exhibit A
greet := \\name, message {
    "Hello {name}! {message}"
}

greet_alice := yolo { greet + "Alice" }
greet_bob   := yolo { greet + "Bob" }

yap(greet_bob("How are you doing?"))
yap(greet_alice("Nice to meet you!"))

// To specify which arguments you want to bake in, add a hashmap.
rude_greet := yolo { greet + %{ "message": "I don't like your face." } }
yap(rude_greet("Bob"))

// Exhibit B
converter := \\symbol, factor, offset, input {
    "{(offset + input) * factor} {symbol}"
}

// To bake multiple arguments, add an array.
miles_to_km          := yolo { converter + ["km", 1.60936, 0] }
pounds_to_kg         := yolo { converter + ["kg", 0.45460, 0] }
farenheit_to_celsius := yolo { converter + ["C", 0.5556, -32] }

yap(miles_to_km(15))
yap(pounds_to_kg(5.5))
yap(farenheit_to_celsius(97))
`,

    mandelbrot: `// Looking to create some fractal fun? With YY, you can easily draw your very own Mandelbrot set.
// Perfect for impressing your math-loving friends or showing off to your imaginary ones. Just sit
// back and let YY do the heavy lifting while you bask in the glory of your own infinite intricacies.

// set up the width and height of the output 'image'
width  := 70
height := 24

// set up the complex plane boundaries for the fractal
real_min := -2.0
real_max := 0.5
imag_min := -1.1
imag_max := 1.1

// a palette of characters to use for the fractal, with increasing 'intensity'
palette := "..--~~:;+=!*#%@"

// the maximum number of iterations to perform, which determines the 'intensity' of each pixel
max_iter := len(palette) - 1

// loop through each pixel in the output image
yall py: height {
    yall px: width {
        // calculate the corresponding complex number for the current pixel
        real := (float(px) / width)  * (real_max - real_min) + real_min
        imag := (float(py) / height) * (imag_max - imag_min) + imag_min

        // set up initial values for the complex number sequence
        x := y := 0.0

        // loop through the complex number sequence until the sequence diverges
        // or the maximum number of iterations is reached
        i := 0
        yoyo x*x + y*y < 4.0 && i < max_iter {
            tmp := x*x - y*y + real
            y   = 2*x*y + imag
            x   = tmp
            i += 1
        }

        // output the corresponding color for the current pixel
        yelp(palette[i])
    }

    // output a newline to move to the next row in the output image
    yap()
}`,

    brainfuck: `// An interpreter for the Brainfuck programming language, written in the YY programming language.
// An interpreter within an interpreter. Interpreter Inception. Interpreception. Interception?
// 
// *BWOOOONNNNGNGGGG* <- Inception's horn sound effect
// -_-                <- DiCaprio's face

// this is an actual "Hello World!" program in Brainfuck
code := "
++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>
---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.
"

// our very own Brainfuck Virtual Machine
ip  := 0  // instruction pointer
dp  := 0  // data pointer
mem := [] // memory

// initialise memory
yall 0..100 { mem << 0 }

// loop through the code and execute each instruction
yoyo ip < len(code) {
    ins := code[ip]
    yif ins == "+" {        // increment the value in memory at the data pointer
        mem[dp] += 1
    } yels yif ins == "-" { // decrement the value in memory at the data pointer
        mem[dp] -= 1
    } yels yif ins == ">" { // move the data pointer to the right
        dp += 1
    } yels yif ins == "<" { // move the data pointer to the left
        dp -= 1
    } yels yif ins == "." { // print the ASCII character for the value in memory at the data pointer
        yelp(chr(mem[dp]))
    } yels yif ins == "[" { // if the memory at the data pointer is 0, jump to the corresponding "]" char
        yif mem[dp] == 0 {
            depth := 1
            yoyo depth != 0 {
                ip += 1
                yif code[ip] == "[" {
                    depth += 1
                } yels yif code[ip] == "]" {
                    depth -= 1
                }
            }
        }
    } yels yif ins == "]" { // if the memory at the data pointer isn't 0, jump back to the corresponding "[" char
        yif mem[dp] != 0 {
            depth := 1
            yoyo depth != 0 {
                ip -= 1
                yif code[ip] == "[" {
                    depth -= 1
                } yels yif code[ip] == "]" {
                    depth += 1
                }
            }
        }
    }

    // move to the next instruction
    ip += 1
}
`,

    maze: `// Have you ever got lost in a supermarket as a child? Perfect!
// We'll recreate that traumatic event by building a maze solver in YY.

maze := [
    "@S@@@@@@@@@@@@@@@@@@@@@@@@@@@@@",
    "@     @   @ @         @       @",
    "@@@@@ @@@ @ @ @@@@@ @@@@@ @@@ @",
    "@   @ @   @ @   @ @     @   @ @",
    "@ @ @ @ @ @ @ @@@ @@@@@@@@@@@ @",
    "@ @     @ @     @ @           @",
    "@@@@@@@@@ @@@ @@@ @@@ @@@ @@@@@",
    "@       @       @       @   @ @",
    "@@@ @ @@@ @@@@@ @@@ @@@@@ @ @ @",
    "@   @     @          @    @   @",
    "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@E@",
]

// locate the starting position by searching for the 'S' character
find_start := \\maze {
    yall row: len(maze)-1 {
        yall col: len(maze[row])-1 {
            yif maze[row][col] == "S" {
                yeet [row, col]
            }
        }
    }

    // 'yikes' terminates the program
    yikes("invalid maze: no starting position found")
}

solve := \\maze {
    start := find_start(maze)
    queue := [start]

    // YY doesn't support sets, so we'll use a hashmap instead
    seen := %{ start: true }

    // keep track of the path to reconstruct our way through the maze
    path := %{ start: null }

    // run until the queue is empty or we found a way out
    yoyo queue {
        // since we're using depth-first seach, we'll get the next position by taking
        // the last element from the queue (we're using queue as a stack)
        cur := yoink(queue)

        // we could change this algorithm to breadth-first search by taking the first element like so
        // cur := yoink(queue, 0)

       // check if we have reached the end
        yif maze[cur[0]][cur[1]] == "E" {
            // backtrack to find and mark the path
            yoyo cur != start {
                maze[cur[0]][cur[1]] = "."
                cur = path[cur]
            }

            maze[cur[0]][cur[1]] = "."

            // exit early, we're done here
            yeet true
        }

        // get neighbours of the current position
        neighbours := []
        yif cur[0] > 0 {
            neighbours << [cur[0]-1, cur[1]]
        }
        yif cur[0] < len(maze)-1 {
            neighbours << [cur[0]+1, cur[1]]
        }
        yif cur[1] > 0 {
            neighbours << [cur[0], cur[1]-1]
        }
        yif cur[1] < len(maze[0])-1 {
            neighbours << [cur[0], cur[1]+1]
        }

        // add unseen neighbours to the queue
        yall neighbours {
            yif !seen[yt] && maze[yt[0]][yt[1]] != "@" {
                seen[yt] = true
                path[yt] = cur
                queue << yt
            }
        }
    }
}

yif solve(maze) {
    // print out the maze with our path
    yall row: maze {
        yall col: row {
            yelp(col)
        }
        yap()
    }
} yels {
    yap("there's no way out :(")
}`,

    regex: `//  “Some people, when confronted with a problem, think: 'I know, I'll use regular expressions'.
//   Now they have two problems.”
//                           -- Jamie Zawinski
//
// But we're not satisfied with just two problems - we like to live dangerously. How about using
// regular expressions in a language that doesn't have them, so we'll have to implement them ourselves
// from scratch? Contgratulations, now we have 3 problems.
//
// We'll yoink Rob Pike's beautifully simple regex matcher from 'The Practice of Programming' (1998).
// It supports 4 special characters: '*', '^', '$' and '.', which account for 95% of real use.
// More details: https://www.cs.princeton.edu/courses/archive/spr09/cos333/beautiful.html

// search for regex anywhere in text
match := \\regex, text {
    yif regex && regex[0] == "^" {
        yeet match_here(regex[1..-1], text)
    }
    yoyo {
        yif match_here(regex, text) {
            yeet true
        }
        yif !text {
            yeet false
        }
        text = text[1..-1]
    }
}

// search for regex at beginning of text
match_here := \\regex, text {
    yif !regex {
        yeet true
    }
    yif regex == "$" {
        yeet text == ""
    }
    yif len(regex) > 1 && regex[1] == "*" {
        yeet match_star(regex[0], regex[2..-1], text)
    }
    yif text && (regex[0] == "." || regex[0] == text[0]) {
        yeet match_here(regex[1..-1], text[1..-1])
    }
    yeet false
}

// search for c*regex at beginning of text
match_star := \\c, regex, text {
    yoyo {
        yif match_here(regex, text) {
            yeet true
        }
        yif !text || (text[0] != c && c != ".") {
            yeet false
        }
        text = text[1..-1]
    }
}

regexes := [ "cat", "^cat", "cat$", "c.*t", "^c.*t$" ]
words   := [ "cat", "cult", "concat", "category", "concatenation" ]

yap("all words:", words)

yall re: regexes {
    result := []
    yall words {
        yif match(re, yt) { result << yt }
    }
    yap("words matching /{re}/: {result}")
}
`,

    sort: `// Sorting algorithms of different sorts.

// Bubble sort: silly name, sillier algorithm (O(n^2)). As useful as 'g' in 'lasagna'. But it is a staple.
bubble_sort := \\arr {
    yall i: len(arr)-2..0 {         // iterate over each element from the second-to-last to the first element
        yall j: 0..i {              // iterate over each element from the first to the i-th element
            yif arr[j] > arr[j+1] { // if the current element is greater than the next element, swap them
                swap(arr, j, j+1)
            }
        }
    }

    arr // return the sorted array
}

// Insertion sort: again, not the fastest (O(n^2)), but it's a patient sorter and gets the job done.
insertion_sort := \\arr {
    i := 1                                // start from the second element
    yoyo i < len(arr) {                   // iterate until the end of the array
        j := i                            // set j as the current element's index
        yoyo j > 0 && arr[j-1] > arr[j] { // keep swapping the element with the previous one
            swap(arr, j, j-1)             // until it's in the correct order

            j -= 1
        }
        i += 1 // move to the next element
    }

    arr // return the sorted array
}

// Quick sort: unlike previous sorts, is quick and nimble like a young yak yodelling in a yurt (O(n log n)).
quick_sort := \\arr, lo, hi {
    yif lo < hi && lo >= 0 { // ensure indices are in correct order
        pivot := arr[hi]     // choose the last element as the pivot
        i := lo - 1          // temporary pivot index

        yall j: lo .. hi-1 {
            yif arr[j] <= pivot {
                i += 1          // move the temporary pivot index forward
                swap(arr, i, j) // swap the current element with the element at the temporary pivot index
            }
        }

        // move the pivot element to the correct pivot position (between the smaller and larger elements)
        i += 1
        swap(arr, i, hi)

        // sort the two partitions
        quick_sort(arr, lo, i - 1) // left side of pivot
        quick_sort(arr, i + 1, hi) // right side of pivot
    }

    arr // return the sorted array
}

// Last, but not least, merge sort.
merge_sort := \\arr {
    // base case: if the array has 0 or 1 element, it is already sorted
    yif len(arr) <= 1 {
        yeet arr
    }

    mid   := len(arr) / 2                   // find the middle index of the array
    left  := merge_sort(arr[0..mid])        // recursively sort the left half of the array
    right := merge_sort(arr[mid..len(arr)]) // recursively sort the right half of the array

    result := [] // initialize an empty list to store the merged result
    i := j := 0  // initialize indices for left and right subarrays

    // merge the sorted left and right subarrays into the result list
    yoyo i < len(left) && j < len(right) {
        yif left[i] < right[j] {
            result << left[i] // append the element from the left subarray to the result
            i += 1            // move to the next element in the left subarray
        } yels {
            result << right[j] // append the element from the right subarray to the result
            j += 1             // move to the next element in the right subarray
        }
    }

    // add any remaining elements from left or right subarrays
    yoyo i < len(left) {
        result << left[i]
        i += 1
    }
    yoyo j < len(right) {
        result << right[j]
        j += 1
    }

    result // return the merged and sorted result
}

// And lastly, let's setup two helper functions, these will come in handy.
// As you probably noticed, functions in YY can be called before their declaration, this isn't C.

// swap function swaps the positions of two elements in an array 
swap := \\arr, i, j {
    tmp    := arr[i]
    arr[i] = arr[j]
    arr[j] = tmp
}

// copy function copies an entire array, using slicing operator
copy := \\arr { arr[0..-1] }

nums := [3, 6, 9, 1, 5, 4, 2, 0, 8, 7]
yap("Original nums:", nums)

// since bubble, insertion, and quick sort are an in-place algorithms,
// we'll pass a copy of the nums array to preserve the original
yap("Bubble sorted:", bubble_sort(copy(nums)))
yap("Quick sorted: ", quick_sort(copy(nums), 0, len(nums)-1))
yap("Insert sorted:", insertion_sort(copy(nums)))

// merge sort doesn't modify the array passed as an arg so we don't have to copy 'nums'
yap("Merge sorted: ", merge_sort(nums))

yap("original array is still a beautiful mess:", nums)
`,

    random: `// A password generator so good, it'll make even the most nefarious hackers throw in the towel.

alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
digits   := "0123456789"
special  := "!?^&*~@#$%"
charset  := alphabet + digits + special

length   := 12
password := ""

// Builtin yahtzee, the magical randomiser, is responsible for generating all things unpredictable.

yall length {
    password += yahtzee(charset)
}

yap("your first secret password:", password)

// Like a genie, yahtzee accepts integers, ranges, arrays, and even strings as offerings to its
// unpredictable power. Just for shit and giggles, we can rewrite the generator to use charset's
// length (integer) as input to yahtzee:

password = ""

yall length {
    idx      := yahtzee(len(charset)-1)
    password += charset[idx]
}

yap("your other secret password:", password)
`,
}

buildSampleSelect()
setSample('hello world')
