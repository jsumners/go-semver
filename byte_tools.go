package semver

type char = int

const (
	dash char = 0x2d // -
	dot       = 0x2e // .
	plus      = 0x2b // +
	star      = 0x2a // *

	lessThan    = 0x3c // <
	equal       = 0x3d // =
	greaterThan = 0x3e // >

	numeral0 = 0x30
	numeral1 = 0x31
	numeral2 = 0x32
	numeral3 = 0x33
	numeral4 = 0x34
	numeral5 = 0x35
	numeral6 = 0x36
	numeral7 = 0x37
	numeral8 = 0x38
	numeral9 = 0x39

	capitalA = 0x41
	capitalB = 0x42
	capitalC = 0x43
	capitalD = 0x44
	capitalE = 0x45
	capitalF = 0x46
	capitalG = 0x47
	capitalH = 0x48
	capitalI = 0x49
	capitalJ = 0x4a
	capitalK = 0x4b
	capitalL = 0x4c
	capitalM = 0x4d
	capitalN = 0x4e
	capitalO = 0x4f
	capitalP = 0x50
	capitalQ = 0x51
	capitalR = 0x52
	capitalS = 0x53
	capitalT = 0x54
	capitalU = 0x55
	capitalV = 0x56
	capitalW = 0x57
	capitalX = 0x58
	capitalY = 0x59
	capitalZ = 0x5a

	lowerA = 0x61
	lowerB = 0x62
	lowerC = 0x63
	lowerD = 0x64
	lowerE = 0x65
	lowerF = 0x66
	lowerG = 0x67
	lowerH = 0x68
	lowerI = 0x69
	lowerJ = 0x6a
	lowerK = 0x6b
	lowerL = 0x6c
	lowerM = 0x6d
	lowerN = 0x6e
	lowerO = 0x6f
	lowerP = 0x70
	lowerQ = 0x71
	lowerR = 0x72
	lowerS = 0x73
	lowerT = 0x74
	lowerU = 0x75
	lowerV = 0x76
	lowerW = 0x77
	lowerX = 0x78
	lowerY = 0x79
	lowerZ = 0x7a
)

func isAlphaChar(i byte) bool {
	return isLowerChar(i) || isCapitalChar(i)
}

func isIntegerChar(i byte) bool {
	c := char(i)
	return c >= numeral0 && c <= numeral9
}

func isCapitalChar(i byte) bool {
	c := char(i)
	return c >= capitalA && c <= capitalZ
}

func isLowerChar(i byte) bool {
	c := char(i)
	return c >= lowerA && c <= lowerZ
}

func isOperatorChar(i byte) bool {
	c := char(i)
	return c == equal || c == lessThan || c == greaterThan
}

func isXRangeChar(i byte) bool {
	c := char(i)
	return c == lowerX || c == capitalX || c == star
}
