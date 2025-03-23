package headers

var allowedSpecialChar = []rune{
	'!',
	'#',
	'$',
	'%',
	'&',
	'\'',
	'*',
	'+',
	'-',
	'.',
	'^',
	'_',
	'`',
	'|',
	'~',
}

func isLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isSpecial(r rune) bool {
	for _, a := range allowedSpecialChar {
		if a == r {
			return true
		}
	}
	return false
}

func checkKey(str string) bool {
	for _, s := range str {
		if !isLetter(s) && !isDigit(s) && !isSpecial(s) {
			return false
		}
	}
	return true
}
