function run() {
    if (!source.value) return;

    output.innerText = ''

    const result = interpret(source.value)
    if (result?.error) {
        output.innerText += "\n" + result.error
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
// interviews. But fear not, for YY is here to help you slay this beast. And rather than printing
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

    'sort':
        `// Sorting algorithms of different sorts.

// Bubble sort: silly name, sillier algorithm (O(n^2)). As useful as 'g' in 'lasagna'. But it is a staple.
bubble_sort := \\arr {
    yall i: len(arr)-2..0 {
        yall j: 0..i {
            yif arr[j] > arr[j+1] {
                arr = swap(arr, j, j+1)
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

    'map et al':
        `// Map, filter, and reduce are The Three Musketeers of functional programming, banding together 
// to process and transform collections with finesse and style.

arr := [1, 2, 3, 4, 5]
yap("original array:", arr)

// Map transforms all elements and returns a shiny new list.
map := \\arr fn {
    acc := []
    yall arr {
        acc = push(acc, fn(yt))
    }
}

double := \\x { x * 2 }
yap("doubled:", map(arr, double))

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

yap("sum of elements:", sum(arr))

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

avg  := sum(arr) / len(arr)
smol := \\x { x < avg }
yap("smaller than average:", filter(arr, smol))
`,

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
    greet     := \\name { yap("Hello, {name}!") }
    greet_yan := greet + "Yan"
    greet_yan() // look ma, no args!

    // you can specify which argument you want to bake in by adding a function to a hashmap
    add   := \\a b { a + b }
    add11 := add + %{ "b": 11 } // baking 'b' into 'add'
    add11(6) // 17

    // but even in yolo mode, division by zero doesn't end well (what did you expect?)
    yap("division by zero:", "weee" / 0)
}`,

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

    "palindrome":
        `// Mirror words, also known as palindromes, are the true testament to the power of written word.
// With every letter and syllable, they silently reflect their greatness back at us.

is_palindrome := \\str {
    n := len(str)
    yall 0..n/2 {
        yif str[yt] != str[n-yt-1] {
            yeet false
        }
    }

    true
}

yall ["racecar", "level", "hello", "world", "1221", "1337"] {
    yif is_palindrome(yt) {
        yap(yt, "is a palindrome")
    } yels {
        yap(yt, "is not a palindrome")
    }
}
`,

    'primes':
        `// If you ever wondered if your favorite number is a prime, wonder no more! YY is here for you to do
// the heavy lifting and separate the primes from the imposter numbers, in a very inefficient manner.

is_prime := \\n {
    yif n < 2  { yeet false }
    yif n == 2 { yeet true }

    yif n % 2 == 0 { yeet false }

    yall 3 .. n/2 + 1 {
        yif n % yt == 0 { yeet false }
    }

    true
}

yall 1..20 {
    yif is_prime(yt) {
        yap("number", yt, "is prime")
    }
}
`,

    'bin conv':
        `// Don't you just hate it when you're sorting a box of uncooked spaghetti by length, and sudddenly a
// stranger comes up to you with a piece of paper covered in ones and zeros? Like, seriously, dude,
// can't you see I'm busy here? But alas, you know you can't resist the temptation of converting
// that binary number to decimal right then and there. Thankfully, with YY, you can swiftly convert
// those ones and zeros without losing your precious pasta-sorting momentum.

bin_to_dec := \\bin {
    dec := 0
    pow := 1

    yall len(bin)-1..0 {
        digit := yif bin[yt] == "0" {
            0 
        } yels yif bin[yt] == "1" {
            1
        } yels {
            yikes("good heavens,", bin, "is not a valid binary number!")
        }

        dec += digit * pow
        pow *= 2
    }

    dec
}


bin := "1000101"
dec := bin_to_dec(bin)
yap(bin, "in binary equals", dec, "in decimal") 
`
}

buildSampleSelect()
setSample('hello world')
