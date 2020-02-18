// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2010 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
//
// This is a modification of the camel case generator by google to translate
// ID and other words into all caps as it's meant to be in golang.

// Package generator provides a CamelCase generator that allows you to
// uppercase certain words if they are found in the identifier
package generator

import (
	"fmt"
	"regexp"
)

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// CamelCase returns a camel case generator that checks for
// any words that we might want to make uppercase
func CamelCase(uppercased []string) (func(string) string, error) {
	rgx := make([]*regexp.Regexp, 0, len(uppercased))
	for _, word := range uppercased {
		exp, err := regexp.Compile(fmt.Sprintf("^%v($|_)", word))
		if err != nil {
			return nil, fmt.Errorf("error converting to the regexp %v: %w", fmt.Sprintf("^%v($|_)", word), err)
		}
		rgx = append(rgx, exp)
	}

	return func(s string) string {
		if s == "" {
			return ""
		}
		t := make([]byte, 0, 32)
		i := 0
		if s[0] == '_' {
			// Need a capital letter; drop the '_'.
			t = append(t, 'X')
			i++
		}
		// Invariant: if the next letter is lower case, it must be converted
		// to upper case.
		// That is, we process a word at a time, where words are marked by _ or
		// upper case letter. Digits are treated as words.
		for ; i < len(s); i++ {
			c := s[i]
			if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
				continue // Skip the underscore in s.
			}
			if isASCIIDigit(c) {
				t = append(t, c)
				continue
			}
			// Perform Regex matching
			var ident []byte
			i, ident = uppercaseExpr(rgx, i, s[i:])
			if ident != nil {
				t = append(t, ident...)
				continue
			}
			// Assume we have a letter now - if not, it's a bogus identifier.
			// The next word is a sequence of characters that must start upper case.
			if isASCIILower(c) {
				c ^= ' ' // Make it a capital letter.
			}
			t = append(t, c) // Guaranteed not lower case.
			// Accept lower case sequence that follows.
			for i+1 < len(s) && isASCIILower(s[i+1]) {
				i++
				t = append(t, s[i])
			}
		}
		return string(t)
	}, nil
}

func uppercaseExpr(rgx []*regexp.Regexp, i int, s string) (int, []byte) {
	for _, exp := range rgx {
		newIndex, identifier := uppercaseWord(exp, i, s)
		if identifier != nil {
			return newIndex, identifier
		}
	}
	return i, nil
}

func uppercaseWord(exp *regexp.Regexp, index int, s string) (int, []byte) {
	identifier := make([]byte, 0, 32)
	indexes := exp.FindStringIndex(s)
	if indexes == nil {
		return index, nil
	}
	for i, c := range s[indexes[0]:indexes[1]] {
		if c == '_' {
			return index + i - 1, identifier
		}
		c ^= ' ' // Make it a capital leter
		identifier = append(identifier, byte(c))
	}
	return indexes[1], identifier
}
