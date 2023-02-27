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

function setSample(sampleName) {
    const { content } = samples[sampleName] || samples['hello']
    source.value = content;
}

function setupSampleSelect() {
    const sampleSelect = document.querySelector('#sample-select')

    for (const sample in samples) {
        const option = document.createElement('option')
        option.value = sample
        option.textContent = sample
        sampleSelect.appendChild(option)
    }

    sampleSelect.addEventListener('change', e => setSample(e.target.value))
}

const samples = {
    'hello': {
        content:
            `// You can edit this code, or select a sample from the dropdown on the right.
// 
// To run the code click 'Run' button.

name := "Yennefer"
yap("Hello, {name}!")` },

    'fibonacci': {
        content:
            `// recursion
fib1 := \\n {
  yif n < 2 { n } yels { fib1(n-1) + fib1(n-2) }
}

yap("seventh Fibonacci number:", fib1(7))

// closure
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
` },

    'map et al': {
        content:
            `// Map, filter, and reduce are like the three musketeers of functional programming, 
// banding together to process and transform collections with finesse and style.

// Map transforms all elements and returns a shiny new list.
map := \\arr fn {
    acc := []
    yall arr {
        acc = push(acc, fn(yt))
    }
}

arr    := [1, 2, 3, 4]
double := \\x { x * 2 }
yap("doubled:", map(arr, double))

// Reduce takes a list and violently smashes it into a single value.
reduce := \\arr initial f {
    result := initial
    yall arr {
        result = f(result, yt)
    }
}

sum := \\arr {
    reduce(arr, 0, \\initial el { initial + el })
}

yap("sum:", sum([1, 2, 3, 4]))

// Filter picks the juiciest elements, leaving all the rest to wither away into obscurity.
filter := \\arr fn {
    acc := []
    yall arr {
        yif fn(yt) {
            acc = push(acc, yt)
        }
    }
    acc
}

arr  := [1, 2, 3, 4, 5, 6, 7]
avg  := sum(arr) / len(arr)
smol := \\x { x < avg }
yap("smaller than average:", filter(arr, smol)) 
`
    }
}

setupSampleSelect()
setSample('hello')
