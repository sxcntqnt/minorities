// find strings in a dir. can be a list of strings

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// inDir is dir to start. must be full path
var inDir = "lsingine/"

var dirsToSkip = []string{
	".git"}

// list of string to search
// each item is intepreted as regex
var findList = []string{
	"haskell",
	"lisp",
	"perl",
	"python",
	"Dagoretti",
}

// fnameRegex. only these are searched
const fnameRegex = `\.json$`

// number of chars (actually bytes) to show before the found string
const charsBefore = 100
const charsAfter = 100

const fileSep = "ff━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"

const occurSep = "oo────────────────────────────────────────────────────────────\n"

const occurBracketL = '〖'
const occurBracketR = '〗'

const posBracketL = '❪'
const posBracketR = '❫'

const fileBracketL = '〘'
const fileBracketR = '〙'

var bigRegexStr = strings.Join(findList, "|")

// var rgx = regexp.MustCompile(regexp.QuoteMeta(bigRegexStr))
var rgx = regexp.MustCompile(bigRegexStr)

func printSliceStr(sliceX []string) error {
	for k, v := range sliceX {
		fmt.Printf("%v %v\n", k, v)
	}
	return nil
}

// scriptPath returns the current running script path
// version 2018-10-07
func scriptPath() string {
	name, errPath := os.Executable()
	if errPath != nil {
		panic(errPath)
	}
	return name
}

func max(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}
func min(x int, y int) int {
	if x > y {
		return y
	} else {
		return x
	}
}

// stringMatchAny return true if x equals any of y
func stringMatchAny(x string, y []string) bool {
	for _, v := range y {
		if x == v {
			return true
		}
	}
	return false
}

func doFile(path string) error {
	var textBytes, er = ioutil.ReadFile(path)
	if er != nil {
		panic(er)
	}

	var indexes = rgx.FindAllIndex(textBytes, -1)

	var bytesLength = len(textBytes)
	if len(indexes) != 0 {
		for _, k := range indexes {
			var foundStart = k[0]
			var foundEnd = k[1]
			var showStart = max(foundStart-charsBefore, 0)
			for !utf8.RuneStart(textBytes[showStart]) {
				showStart = max(showStart-1, 0)
			}

			var showEnd = min(foundEnd+charsAfter, bytesLength-1)
			for !utf8.RuneStart(textBytes[showEnd]) {
				showEnd = min(showEnd-1, bytesLength)
			}

			// 			fmt.Printf("%s〖%s〗%s\n", textBytes[showStart:foundStart],
			fmt.Printf("%c%d%c %s%c%s%c%s\n",
				posBracketL,
				utf8.RuneCount(textBytes[0:foundStart+1]),
				posBracketR,
				textBytes[showStart:foundStart],
				occurBracketL,
				textBytes[foundStart:foundEnd],
				occurBracketR,
				textBytes[foundEnd:showEnd])
			fmt.Println(occurSep)
		}
		fmt.Printf("%v %c%v%c\n", len(indexes), fileBracketL, path, fileBracketR)
		fmt.Println(fileSep)
	}
	return nil
}

func main() {

	fmt.Println("-*- coding: utf-8; mode: xah-find-output -*-")
	fmt.Printf("%v\n", time.Now())
	fmt.Printf("Script: %v\n", filepath.Base(scriptPath()))
	fmt.Printf("in dir: %v\n", inDir)
	fmt.Printf("file regex filter: %v\n", fnameRegex)
	fmt.Printf("findList:\n")
	printSliceStr(findList)
	fmt.Println()
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")

	var pWalker = func(pathX string, infoX os.FileInfo, errX error) error {
		// first thing to do, check error. and decide what to do about it
		if errX != nil {
			fmt.Printf("error 「%v」 at a path 「%q」\n", errX, pathX)
			return errX
		}
		if infoX.IsDir() {
			if stringMatchAny(filepath.Base(pathX), dirsToSkip) {
				return filepath.SkipDir
			}
		} else {
			var x, err = regexp.MatchString(fnameRegex, filepath.Base(pathX))
			if err != nil {
				panic("stupid MatchString error 59767")
			}
			if x {
				doFile(pathX)
			}
		}
		return nil
	}
	err := filepath.Walk(inDir, pWalker)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", inDir, err)
	}
	fmt.Println("\nDone.")
}
