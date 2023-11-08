package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"

	"yy/eval"
	"yy/lexer"
	"yy/object"
	"yy/parser"
	"yy/yikes"
)

const version = "v0.0.1"

var debug = false

func main() {
	switch len(os.Args) {
	case 1:
		repl()

	case 2:
		runFile(os.Args[1])

	default:
		fmt.Println("usage: yy [path_to_script]")
	}
}

func runFile(f string) {
	src, err := os.ReadFile(f)
	if err != nil {
		fmt.Println("error: couldn't read file: " + f)
		os.Exit(1)
	}

	l := lexer.New(string(src))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, errMsg := range p.Errors() {
			fmt.Println(yikes.PrettyError(src, errMsg.Offset, errMsg.Msg))
		}
		os.Exit(1)
	}

	env := object.NewEnvironment()

	result := eval.Eval(program, env)
	if evalError, ok := result.(*object.Error); ok {
		fmt.Println(yikes.PrettyError(src, evalError.Pos, evalError.Msg))
		os.Exit(1)
	}
}

const (
	greet   = "YeetYoink " + version
	prompt  = "yy> "
	padLeft = "    "
)

func repl() {
	in := os.Stdin
	out := os.Stdout
	scanner := bufio.NewScanner(in)

	env := object.NewEnvironment()

	fmt.Println(greet)

	for {
		fmt.Fprint(out, prompt)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if debug {
			io.WriteString(out, padLeft)
			io.WriteString(out, program.String())
			io.WriteString(out, "\n")
		}

		if len(p.Errors()) > 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, msg.Error()+"\n")
			}
			continue
		}

		result := eval.Eval(program, env)
		if result != nil {
			io.WriteString(out, result.String())
			io.WriteString(out, "\n")
		}

		// ever so often, delight the user with a random yak fact
		if rand.Intn(7) == 1 {
			idx := rand.Intn(len(yakFacts))
			msg := fmt.Sprintf("Yak Fact #%d: %s\n", idx+1, yakFacts[idx])
			io.WriteString(out, msg)
		}
	}
}

var yakFacts = [...]string{
	"Yaks are large mammals that are native to the Himalayan region of Central Asia.",
	"Yaks are the lumberjacks of the Himalayas - they have big horns, a hump on their back, and they're not afraid of a little cold weather.",
	"Yak milk is so nutritious that it's the whey protein powder of the Himalayas. No wonder those Sherpas can climb mountains like they're walking on flat ground.",
	"Yak wool is so soft and warm that it's like wearing a hug from a big, fluffy friend who's also really good at surviving in sub-zero temperatures.",
	"Yaks are like the SUVs of the animal kingdom - they can carry a ton of stuff and handle any terrain, all while looking stylish with their shaggy coats.",
	"Yaks are the Ron Swanson of the animal kingdom - tough, hardworking, and they don't take any crap from anyone.",
	"Yak meat is so lean and high in protein that it's the chicken breast of the Himalayas. Just don't tell the yaks that - they're already self-conscious enough about their figure.",
	"Yaks are like the Swiss Army knives of the animal kingdom - they can provide milk, meat, wool, and transportation, all while looking like they're ready for a fashion shoot in GQ.",
	"Yaks are well-adapted to high-altitude environments, with a thick coat of hair that protects them from the cold and harsh weather.",
	"Yaks are used by local communities in the Himalayas for their milk, meat, and wool.",
	"The milk from yaks is rich in protein and fat, and is often used to make butter, cheese, and yogurt.",
	"Yak wool is used to make clothing, blankets, and other textiles, and is prized for its softness and warmth.",
	"Yaks are also used as pack animals, carrying goods and supplies across rugged mountain terrain.",
	"Despite their tough exterior, yaks are also known for their gentle and docile nature, and are often kept as pets or used for ceremonial purposes in some cultures.",
	"Despite their shaggy coats, yaks are surprisingly good swimmers. They're like the Michael Phelps of the Himalayas.",
	"Yaks are great at social distancing. They've been practicing it for centuries, long before it was cool.",
	"If you ever need to cross a rickety old bridge over a roaring river, bring a yak with you. They have the balance and poise of a ballerina on a tightrope.",
	"Yaks are known for their strong digestive systems, capable of breaking down tough vegetation. It's like they have a built-in garbage disposal.",
	"If yaks ever decided to start a boy band, they could call themselves the Yakstreet Boys. (I'm sorry, that one was bad.)",
	"Yaks may look docile and cuddly, but don't be fooled - they're fierce protectors of their herds. They're the bouncers of the Himalayas.",
	"Yaks are pretty low-maintenance animals, but they do have one weakness: they're terrible at texting. Their hooves are just too big for those tiny touchscreens.",
	"Yaks are expert mountaineers, scaling steep slopes with ease. If they were climbers, they'd be sponsored by Red Bull.",
	"Yaks are surprisingly fast runners, able to reach speeds of up to 25 miles per hour. They could probably give Usain Bolt a run for his money.",
}
