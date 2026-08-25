package weathertool

import "bytes"

func cleanString(b []byte) []byte {
	/*
		After adding this -> b = bytes.ReplaceAll(b, []byte("\n"), []byte(""))
		The output was still outputting escape sequence chars. Claude suggested
		adding the subsequent two lines, works!
	*/
	b = bytes.ReplaceAll(b, []byte("\n"), []byte(""))  // strip real newline bytes (pretty-printed XML)
	b = bytes.ReplaceAll(b, []byte(`\n`), []byte(" ")) // replace literal "\n" text markers with a space
	b = bytes.Join(bytes.Fields(b), []byte(" "))       // collapse any resulting runs of whitespace

	return b
}
