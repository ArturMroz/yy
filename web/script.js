function run() {
    if (!source.value) return;

    output.innerText = ''

    const result = interpret(source.value)
    if (result?.error) {
        if (output.innerText) output.innerText += '\n'
        output.innerText += result.error
    }
}

function captureLog(msg) {
    const li = document.createElement("li")
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
        e.preventDefault(); // prevent tabbing out from textarea
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
    'hello world':
        `// You can edit this code, or select a sample from the dropdown on the right.
//
// To run the code press Ctrl+Enter or click 'Run' button.

name := "Yennefer"
yap("Yo, {name}!")`,

    "fizzbuzz":
        `// Ah, FizzBuzz, the timeless test that weeds out the 10x engineers from the wannabes in programming
// interviews. But fear not, YY is here to help you slay this beast. And rather than printing
// the mundane FizzBuzz, we'll print out the magnificent YeetYoink instead, for it truly captures
// the essence of YY.

yall 1..100 {
    yif yt % 15 == 0 {
        yap("YeetYoink")
    } yels yif yt % 3 == 0 {
        yap("Yeet")
    } yels yif yt % 5 == 0 {
        yap("Yoink")
    } yels {
        yap(yt)
    }
}
`,

    'fibonacci':
        `// Implementation of Fibbonacci numbers using two ways: recursion and closure. Just like choosing
// between pizza and tacos, there is no right or wrong way to do it, both are equally satisfying.
// And while these methods may not be the fastest, they add some spicy flavor to this demo.

// Recursion
fib := \\n {
    yif n < 2 { n } yels { fib(n-1) + fib(n-2) }
}

yap("seventh Fibonacci number:", fib(7))

// Closure
fib_gen := \\{
    a := 0
    b := 1
    \\{
        temp := a
        a = b
        b += temp
        a
    }
}

f := fib_gen()
yap("consecutive Fibonacci numbers:", f(), f(), f(), f(), f())
` ,

    'yolo':
        `// Yolo Mode allows you to do things that would make other self-respecting languages blush.
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
    // (check out 'baking' example for more details)
    greet     := \\name { yap("Hello, {name}!") }
    greet_yan := greet + "Yan"
    greet_yan() // look ma, no args!

    // but even in yolo mode, division by zero doesn't end well (what did you expect?)
    yap("division by zero:", "weee" / 0)
}`,

    'bake':
        `// Brace yourselves, we're about to go into YOLO mode! We'll be adding numbers, arrays, and hashmaps
// to a function like a mad scientist adding ingredients to a cauldron. This magically bakes the
// arguments into the function, turning it into a deliciously self-contained recipe for success.
// Some fancy folks call this 'partial function application' or 'currying', we'll just call it baking.

// Exhibit A
greet := \\name, message {
    "Hello {name}! {message}"
}

greet_alice := yolo { greet + "Alice" }
greet_bob   := yolo { greet + "Bob" }

yap(greet_bob("How are you doing?"))
yap(greet_alice("Nice to meet you!"))

// To specify which arguments you want to bake in, add the function to a hashmap.
rude_greet := yolo { greet + %{ "message": "I don't like your face." } }
yap(rude_greet("Bob"))

// Exhibit B
converter := \\symbol, factor, offset, input {
    result := (offset + input) * factor
    "{result} {symbol}"
}

// To bake multiple arguments, add an array to a function.
miles_to_km          := yolo { converter + ["km", 1.60936, 0] }
pounds_to_kg         := yolo { converter + ["kg", 0.45460, 0] }
farenheit_to_celsius := yolo { converter + ["C", 0.5556, -32] }

yap(miles_to_km(15))
yap(pounds_to_kg(5.5))
yap(farenheit_to_celsius(97))
`,

    'map et al':
        `// Map, filter, and reduce are The Three Musketeers of functional programming, banding together
// to process and transform collections with finesse and style.

// Map transforms all elements and returns a shiny new list.
map := \\arr fn {
    acc := []
    yall arr {
        acc = push(acc, fn(yt))
    }
}

// Reduce violently smashes a list into a single value.
reduce := \\arr initial fn {
    result := initial
    yall arr {
        result = fn(result, yt)
    }
}

sum := \\arr {
    add := \\a b { a + b }
    reduce(arr, 0, add)
}

// Filter picks the juiciest elements, leaving the rest to wither away into obscurity.
filter := \\arr fn {
    acc := []
    yall arr {
        yif fn(yt) {
            acc = push(acc, yt)
        }
    }
    acc
}

arr := [1, 2, 3, 4]
yap("original array:", arr)

tripled := map(arr, \\x { x * 3 })
yap("tripled:", tripled)

summed := sum(tripled)
yap("sum of elements:", summed)

avg  := summed / len(arr)
smol := \\x { x < avg }
yap("smaller than average:", filter(tripled, smol))
`,

    'mandelbrot':
        `// Looking to create some fractal fun? With YY, you can easily draw your very own Mandelbrot set.
// Perfect for impressing your math-loving friends or showing off to your imaginary ones. Just sit
// back and let YY do the heavy lifting while you bask in the glory of your own infinite intricacies.

width  := 90
height := 24

real_min := -2.0
real_max := 0.5
imag_min := -1.1
imag_max := 1.1

palette  := "...--~~+:;=*!#$%W@"
max_iter := len(palette) - 1

yall py: 0..height {
    yall px: 0..width {
        real := (float(px) / width)  * (real_max - real_min) + real_min
        imag := (float(py) / height) * (imag_max - imag_min) + imag_min

        x := y := 0.0

        i := 0
        yet x*x + y*y < 4.0 && i < max_iter {
            tmp := x*x - y*y + real
            y   = 2*x*y + imag
            x   = tmp
            i += 1
        }

        yelp(palette[i])
    }

    yap()
}
`,

    'brainfuck':
        `// An interpreter for the Brainfuck programming language, written in the YY programming language.
// An interpreter within an interpreter. Interpreter Inception. Interpreception. Interception?
// *BWOOOONNNNGNGGGG* <- Inception's horn sound effect
// -_-                <- DiCaprio's face

// This is an actual "Hello World!" program in Brainfuck
hello_world := "
++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>
---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.
"

// Our very own Brainfuck VM
ip  := 0  // instruction pointer
dp  := 0  // data pointer
mem := [] // memory

// initialise memory
yall 0..100 { mem = push(mem, 0) }

code := hello_world

yet ip < len(code) {
    ins := code[ip]
    yif ins == "+" {
        mem[dp] += 1
    } yels yif ins == "-" {
        mem[dp] -= 1
    } yels yif ins == ">" {
        dp += 1
    } yels yif ins == "<" {
        dp -= 1
    } yels yif ins == "." {
        yelp(chr(mem[dp]))
    } yels yif ins == "[" {
        yif mem[dp] == 0 {
            depth := 1
            yet depth != 0 {
                ip += 1
                yif code[ip] == "[" {
                    depth += 1
                } yels yif code[ip] == "]" {
                    depth -= 1
                }
            }
        }
    } yels yif ins == "]" {
        yif mem[dp] != 0 {
            depth := 1
            yet depth != 0 {
                ip -= 1
                yif code[ip] == "[" {
                    depth -= 1
                } yels yif code[ip] == "]" {
                    depth += 1
                }
            }
        }
    }

    ip += 1
}
`,

    'sort':
        `// Sorting algorithms of different sorts.

// Bubble sort: silly name, sillier algorithm (O(n^2)). As useful as 'g' in 'lasagna'. But it is a staple.
bubble_sort := \\arr {
    yall i: len(arr)-2..0 {
        yall j: 0..i {
            yif arr[j] > arr[j+1] {
                tmp      := arr[j]
                arr[j]   = arr[j+1]
                arr[j+1] = tmp
            }
        }
    }

    arr
}

// Quick sort, unlike bubble sort, is quick and nimble like a young yak yodelling in a yurt (O(n log n)).
quick_sort := \\arr {
    yif len(arr) < 2 {
        yeet arr
    }

    pivot  := arr[len(arr) / 2]
    left   := []
    right  := []
    middle := []

    yall arr {
        yif yt < pivot {
            left = push(left, yt)
        } yels yif yt > pivot {
            right = push(right, yt)
        } yels {
            middle = push(middle, yt)
        }
    }

    quick_sort(left) + middle + quick_sort(right)
}

nums := [3, 6, 9, 1, 5, 4, 2, 0, 8, 7]

yap("Bubble sorted:", bubble_sort(nums))
yap("Quick sorted: ", quick_sort(nums))
yap("Btw, original array is still there, untouched:", nums)
`,

    "random":
        `// A password generator so good, it'll make even the most nefarious hackers throw in the towel.

alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
digits   := "0123456789"
special  := "!?^&*~@#$%"
charset  := alphabet + digits + special

length   := 12
password := ""

// Builtin yahtzee, the magical randomiser, is responsible for generating all things unpredictable.

yall 0..length {
    password += yahtzee(charset)
}

yap("your first secret password:", password)

// Like a genie, yahtzee accepts integers, ranges, arrays, and even strings as offerings to its
// unpredictable power. Just for shit and giggles, we can rewrite the generator to use charset's
// length (integer) as input to yahtzee:

password = ""

yall 0..length {
    idx := yahtzee(len(charset)-1)
    password += charset[idx]
}

yap("your other secret password:", password)
`,



}

buildSampleSelect()
setSample('hello world')
