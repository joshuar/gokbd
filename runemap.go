// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package gokbd

// CharVariants represents the upper and lower case variants for a charactor
type CharVariants struct {
	lc rune
	uc rune
	// cc rune
}

var runeMap = map[int]CharVariants{
	2:  {lc: '1', uc: '!'},
	3:  {lc: '2', uc: '@'},
	4:  {lc: '3', uc: '#'},
	5:  {lc: '4', uc: '$'},
	6:  {lc: '5', uc: '%'},
	7:  {lc: '6', uc: '^'},
	8:  {lc: '7', uc: '&'},
	9:  {lc: '8', uc: '*'},
	10: {lc: '9', uc: '('},
	11: {lc: '0', uc: ')'},
	12: {lc: '-', uc: '_'},
	13: {lc: '=', uc: '+'},
	14: {lc: '\b', uc: '\b'},
	15: {lc: '\t', uc: '\t'},
	16: {lc: 'q', uc: 'Q'},
	17: {lc: 'w', uc: 'W'},
	18: {lc: 'e', uc: 'E'},
	19: {lc: 'r', uc: 'R'},
	20: {lc: 't', uc: 'T'},
	21: {lc: 'y', uc: 'Y'},
	22: {lc: 'u', uc: 'U'},
	23: {lc: 'i', uc: 'I'},
	24: {lc: 'o', uc: 'O'},
	25: {lc: 'p', uc: 'P'},
	26: {lc: '[', uc: '{'},
	27: {lc: ']', uc: '}'},
	28: {lc: '\n', uc: '\n'},
	30: {lc: 'a', uc: 'A'},
	31: {lc: 's', uc: 'S'},
	32: {lc: 'd', uc: 'D'},
	33: {lc: 'f', uc: 'F'},
	34: {lc: 'g', uc: 'G'},
	35: {lc: 'h', uc: 'H'},
	36: {lc: 'j', uc: 'J'},
	37: {lc: 'k', uc: 'K'},
	38: {lc: 'l', uc: 'L'},
	39: {lc: ';', uc: ':'},
	40: {lc: '\'', uc: '"'},
	41: {lc: '`', uc: '~'},
	43: {lc: '\\', uc: '|'},
	44: {lc: 'z', uc: 'Z'},
	45: {lc: 'x', uc: 'X'},
	46: {lc: 'c', uc: 'C'},
	47: {lc: 'v', uc: 'V'},
	48: {lc: 'b', uc: 'B'},
	49: {lc: 'n', uc: 'N'},
	50: {lc: 'm', uc: 'M'},
	51: {lc: ',', uc: '<'},
	52: {lc: '.', uc: '>'},
	53: {lc: '/', uc: '?'},
	57: {lc: ' ', uc: ' '},
}

// CodeAndCase returns the keycode and whether the key was
// an upper or lowercase rune for the typed key
func CodeAndCase(r rune) (int, bool) {
	for k, v := range runeMap {
		switch {
		case r == v.lc:
			return k, false
		case r == v.uc:
			return k, true
		}
	}
	return 0, false
}
