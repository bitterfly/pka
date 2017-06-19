package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/bitterfly/pka/regex"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func readWord(fileName string) chan string {
	dict := make(chan string, 1000)
	go func() {
		defer close(dict)

		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			dict <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()
	return dict
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if len(os.Args) != 2 {
		fmt.Printf("usage: pka filename\n")
		return
	}

	//===================================

	epsilon := regex.EmptyExpressionNDFA(3)
	epsilon.Print()

	// letter := regex.LetterExpressionNDFA(5, 'a')
	// letter.Print()

	// union := regex.UnionExpressionsNDFA(2, epsilon, letter)

	// epsilon2 := regex.EmptyExpressionNDFA(8)
	// doubleUnion := regex.UnionExpressionsNDFA(1, epsilon2, union)

	// concatenation := regex.ConcatenateExpressionsNDFA(doubleUnion, letter)
	// concatenation.Print()
	// concatenation.Dot("a.dot")

	//=====================
	// dict := readWord(os.Args[1])

	// start_time := time.Now()

	// test := dfa.BuildDFAFromDict(dict)
	// elapsed := time.Since(start_time)
	// //test.Print()

	// dict = readWord(os.Args[1])
	// fmt.Printf("Correct language: %v\nTime: %s\n", test.CheckLanguage(dict), elapsed)
	// //fmt.Printf("Is minimal? %v\n", (i == eq_c))
	// fmt.Printf("Number of states: %d\n", test.GetNumStates())
	// fmt.Printf("Number of eq classes: %d\n", test.GetNumEqClasses())

	// test.DotGraph("a.dot")
	// fmt.Printf("Check real minimality: %v\n", test.CheckMinimal())
}
