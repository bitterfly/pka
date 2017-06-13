package dfa

import (
	"fmt"
	"sort"
)

//states are consecutive numbers
//start state is always 1
type DFA struct {
	maxState    int
	finalStates map[int]struct{}
	delta       DeltaTransitions
}

func NewDFA(maxState int, _finalStates []int, delta map[Transition]int) *DFA {
	finalStates := make(map[int]struct{})
	for _, state := range _finalStates {
		finalStates[state] = struct{}{}
	}

	return &DFA{
		maxState:    maxState,
		finalStates: finalStates,
		delta:       *NewDeltaTransitions(delta),
	}
}

func EmptyAutomaton() *DFA {
	return &DFA{
		maxState:    1,
		finalStates: make(map[int]struct{}),
		delta:       *NewDeltaTransitions(make(map[Transition]int)),
	}
}

func BuildDFAFromDict(dict []string) (*DFA, int) {
	checked := &EquivalenceTree{}
	dfa := EmptyAutomaton()

	for _, word := range dict {
		remaining, lastState := dfa.delta.commonPrefix(word)

		if dfa.delta.hasChildren(lastState) {
			dfa.reduce(lastState, checked)
		}

		if remaining == "" {
			dfa.makeFinal(lastState)
		} else {
			dfa.AddWord(lastState, remaining)
		}

	}
	dfa.reduce(1, checked)
	return dfa, CountStuff(checked)
}

func GetTimesReduce() int {
	return TimesReduce
}

var TimesReduce int = 0

func (d *DFA) reduce(state int, checked *EquivalenceTree) {

	TimesReduce += 1
	children := d.delta.stateToTransitions[state]
	child := children[len(children)-1]
	if d.delta.hasChildren(child.state) {
		d.reduce(child.state, checked)
	}

	childEquivalenceClass := *NewEquivalenceClass(d.isFinal(child.state), d.delta.getChildren(child.state))
	childEquivalenceNode := *NewEquivalenceNode(child.state, childEquivalenceClass)

	// if child.state == 27 {
	// 	fmt.Printf("Equivalence class 27: %v \n", childEquivalenceClass)
	// 	equivalenceClass15 := *NewEquivalenceClass(d.isFinal(15), d.delta.getChildren(15))
	// 	fmt.Printf("Equivalence class 15: %v \n", equivalenceClass15)
	// }

	checked_state, ok := checked.Find(childEquivalenceNode)
	if checked_state == child.state {
		return
	}

	if ok {
		// if child.state == 168484 {
		// 	fmt.Printf("We found equivalent to 27: %d\n", checked_state)
		// 	fmt.Printf("We are removing: %d, %c, %d\n and %d wont be final\n", state, child.letter, child.state, child.state)
		// 	fmt.Printf("We are addint transition: %d, %c, %d\n\n", state, child.letter, checked_state)

		// 	last := childEquivalenceClass.children[len(childEquivalenceClass.children)-1]

		// 	fmt.Printf("We want to remove: %d, %c, %d\n", child.state, last.letter, last.state)
		// 	if len(d.delta.stateToTransitions[child.state]) > 1 {
		// 		panic("We have more than one child ")
		// 	}
		// }
		d.delta.removeTransition(state, child.letter, child.state)
		d.removeState(child.state)
		d.delta.removeTransitionsFor(child.state)

		d.delta.addTransition(state, child.letter, checked_state)
	} else {
		Insert(&checked, childEquivalenceNode)
	}
}

func (d *DFA) AddWord(state int, word string) {
	d.addNewStates(len(word))
	d.finalStates[d.maxState] = struct{}{}
	d.delta.addWord(state, d.maxState-len(word)+1, word)
}

func (d *DFA) isFinal(state int) bool {
	_, ok := d.finalStates[state]
	return ok
}

func (d *DFA) checkEquivalentStates(first int, second int) bool {
	return (d.isFinal(first) == d.isFinal(second)) &&
		(d.delta.compareOutgoing(first, second))
}

func (d *DFA) addNewStates(number int) {
	d.maxState += number
}

func (d *DFA) makeFinal(state int) {
	d.finalStates[state] = struct{}{}
}

func (d *DFA) removeState(state int) {
	if d.isFinal(state) {
		delete(d.finalStates, state)
	}
}

//===========================Human Friendly======================================
func (d *DFA) CountStates() (int, int) {
	states := make(map[int]struct{})

	for tr, i := range d.delta.transitionToState {
		states[tr.state] = struct{}{}
		states[i] = struct{}{}
	}

	states_2 := make(map[int]struct{})
	for i, tr_2 := range d.delta.stateToTransitions {
		states_2[i] = struct{}{}
		for _, t := range tr_2 {
			states_2[t.state] = struct{}{}
		}
	}

	return len(states), len(states_2)
}

func (d *DFA) sortedFinalStates() []int {
	var states []int
	for k, _ := range d.finalStates {
		states = append(states, k)
	}
	sort.Ints(states)
	return states
}

func (d *DFA) Print() {
	fmt.Printf("====DFA====\n")
	fmt.Printf("Max: %d, Final: %v\n", d.maxState, d.sortedFinalStates())
	d.PrintFunction()
	fmt.Printf("\n====AFD====\n")
}

func (d *DFA) PrintFunction() {
	fmt.Printf("(p, a) -> q\n\n")
	for transition, goalState := range d.delta.transitionToState {
		fmt.Printf("(%d, %c) -> %d)\n", transition.state, transition.letter, goalState)
	}
	fmt.Printf("\np -> (a, q)\n\n")
	for initialState, children := range d.delta.stateToTransitions {
		fmt.Printf("%d -> %v \n", initialState, children)
	}
}

func (d *DFA) Traverse(word string) {
	ok, state := d.delta.traverse(word)
	if !ok {
		fmt.Printf("Not in the automation - %s\n", word)
	} else {
		fmt.Printf("%s leads to %d\n", word, state)
	}
}

func (d *DFA) FindCommonPrefix(word string) {
	remaining, state := d.delta.commonPrefix(word)
	fmt.Printf("Word: %s\nRemaining: %s, last_state: %d\n\n", word, remaining, state)
}

func (d *DFA) CheckLanguage(dict []string) bool {
	for _, word := range dict {
		ok, state := d.delta.traverse(word)
		if !ok {
			fmt.Printf("No transition: %s\n", word)
			return false
		}
		if !d.isFinal(state) {
			fmt.Printf("First failing word: %s\n", word)
			return false
		}
	}
	return true
}

func (d *DFA) Check() {
	fmt.Printf("Equiv 1:1 %v\n", d.checkEquivalentStates(1, 1))
	fmt.Printf("Equiv 1:2 %v\n", d.checkEquivalentStates(1, 2))
}

func (d *DFA) GetMaxState() int {
	return d.maxState
}

func (d *DFA) CheckMinimal() bool {
	for s1, tr1 := range d.delta.stateToTransitions {
		for s2, tr2 := range d.delta.stateToTransitions {
			if s1 != s2 {
				if (len(tr1) == len(tr2)) && (d.isFinal(s1) == d.isFinal(s2)) && compareTransitionSlices(tr1, tr2) == 0 {
					fmt.Printf("Equal states: %d, %d\n", s1, s2)
					fmt.Printf("%d: %v\n", s1, d.delta.stateToTransitions[s1])
					fmt.Printf("%d: %v\n", s2, d.delta.stateToTransitions[s2])
					return false
				}
			}
		}
	}

	return true
}
