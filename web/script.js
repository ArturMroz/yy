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

    sampleSelect.addEventListener('change', e => {
        const sampleName = e.target.value
        setSample(sampleName)
    })
}

const samples = {
    'hello': {
        content:
            `// you can edit this code, or select a sample from the dropdown on the right
// 
// to run the code click 'Run' button

name := "Yan"\nyap("Hello, {name}!")` },

    'fibonacci': {
        content:
            `// recursion
fib1 := \\n {
  yif n < 2 { n } yels { fib1(n-1) + fib1(n-2) }
};

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
}

setupSampleSelect()
setSample('hello')
