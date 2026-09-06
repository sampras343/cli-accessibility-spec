package probe

import "regexp"

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
var colorFGPattern = regexp.MustCompile(`\x1b\[(3[0-7]|9[0-7]|38;5;\d+|38;2;\d+;\d+;\d+)m`)
var colorBGPattern = regexp.MustCompile(`\x1b\[(4[0-7]|10[0-7]|48;5;\d+|48;2;\d+;\d+;\d+)m`)
var eightBitPattern = regexp.MustCompile(`\x1b\[(38|48);5;(\d+)m`)
var twentyFourBitPattern = regexp.MustCompile(`\x1b\[(38|48);2;\d+;\d+;\d+m`)
var fourBitFGPattern = regexp.MustCompile(`\x1b\[(3[0-7]|9[0-7])m`)

type ColorClassification struct {
	FourBit       int
	EightBit      int
	TwentyFourBit int
}

func HasANSI(data []byte) bool {
	return ansiPattern.Match(data)
}

func StripANSI(data []byte) []byte {
	return ansiPattern.ReplaceAll(data, nil)
}

func HasColorCodes(data []byte) bool {
	return colorFGPattern.Match(data) || colorBGPattern.Match(data)
}

func ClassifyColors(data []byte) ColorClassification {
	var c ColorClassification
	c.TwentyFourBit = len(twentyFourBitPattern.FindAll(data, -1))
	c.EightBit = len(eightBitPattern.FindAll(data, -1))
	c.FourBit = len(fourBitFGPattern.FindAll(data, -1))
	return c
}
