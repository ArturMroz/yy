function run() {
    if (!source.value) return;

    output.innerText = ''

    const result = interpret(source.value)
    if (result?.error) {
        output.innerText = result.error
    }
}

function captureLog(msg) {
    const li = document.createElement("li")
    li.innerText = msg
    output.appendChild(li)
}

window.console.log = captureLog

source.addEventListener('keydown', (e) => {
    if (e.keyCode === 9) { // tab
        document.execCommand('insertText', false, '    ') // tab is 4 spaces
        e.preventDefault(); // prevent tabbing out from textarea
    } else if (e.keyCode === 13 && e.ctrlKey) { // ctrl + enter
        run()
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
}

const samples = {
    'hello':
        `// You can edit this code, or select a sample from the dropdown on the right.
//
// To run the code press Ctrl+Enter or click 'Run' button.

name := "Yennefer"
yap("Hello, {name}!")`,

    'fibonacci':
        `// Implementation of Fibbonacci nubmers using two ways: recursion and closure.
// And yes, these aren't the most efficent ways to calculate the sequence but they make the demo more interesting.

// Recursion
fib1 := \\n {
    yif n < 2 { n } yels { fib1(n-1) + fib1(n-2) }
}

yap("seventh Fibonacci number:", fib1(7))

// Closure
fib2 := \\{
    a := 0
    b := 1
    \\{
        temp := a
        a = b
        b += temp
        a
    }
}

f := fib2()
yap("consecutive Fibonacci numbers:", f(), f(), f(), f(), f())
` ,

    'map et al':
        `// Map, filter, and reduce are The Three Musketeers of functional programming,
// banding together to process and transform collections with finesse and style.

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

    // you can do useful stuff too, like baking an argument into a function
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

    'sort':
        `// Sorting algorithms of different sorts.

// To start we capture this unexpecting array from the wild, and make it our unwilling test subject.
nums := [6, 3, 9, 1, 0, 7, 2, 5, 8, 4]

// Bubble sort: silly name, sillier algorithm (O(n^2)). It is a staple though.
bubble_sort := \\arr {
    yall i: len(arr) - 1..0 {
        yall j: 0..i {
            yif arr[j] > arr[j+1] {
                arr = swap(arr, j, j+1)
            }
        }
    }
    arr
}

yap("Bubble sorted nums:", bubble_sort(nums))
yap("Original array is still there, untouched:", nums)

// Quick sort, unlike bubble sort, is quick and nimble like a young yak frolicking in a field (O(n log n)).
// It uses divide-and-conquer strategy to slice and dice an array into smaller, more manageable pieces.
qsort := \\arr {
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

    qsort(left) + middle + qsort(right)
}

yap("Quick sorted nums:", qsort(nums))
`
}

buildSampleSelect()
setSample('hello')
