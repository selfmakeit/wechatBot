package utils

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"strconv"

	"github.com/printfcoder/goutils/mathutils"
)

const (
	// SPACE is a String for a space character.
	SPACE = " "

	// EMPTY is the empty String ""
	EMPTY = ""

	// IndexNotFound IndexNotFound
	IndexNotFound = -1
)

// region empty checks

// IsEmpty checks if the address of point cs is empty
func IsEmpty(cs string) bool {
	return len(cs) == 0
}
func TrimAllSpace(src string) (dist string) {
	src = strings.Trim(src, "\u3000")
	if len(src) == 0 {
		return
	}

	r, distR := []rune(src), []rune{}
	for i := 0; i < len(r); i++ {
		/*
			10 \n
			32 空格
			8197 中文空格
		*/

		if r[i] == 8197 || r[i] == 32 {
			continue
		}

		distR = append(distR, r[i])
	}
	dist = string(distR)
	return
}
//去除所有中文空格，包括文本中间的
func TrimAllChineseSpace(src string) (dist string) {
	src = strings.Trim(src, "\u3000")
	if len(src) == 0 {
		return
	}

	r, distR := []rune(src), []rune{}
	for i := 0; i < len(r); i++ {
		/*
			10 \n
			32 空格
			8197 中文空格
		*/

		if r[i] == 8197{
			continue
		}

		distR = append(distR, r[i])
	}
	dist = string(distR)
	return
}

// 去除所有回车和空格，包括中英文
func TrimAllSpaceAndEnter(src string) (dist string) {
	src = strings.Trim(src, "\u3000")
	if len(src) == 0 {
		return
	}

	r, distR := []rune(src), []rune{}
	for i := 0; i < len(r); i++ {
		/*
			10 \n
			32 空格
			8197 中文空格
		*/
		if r[i] == 8197 || r[i] == 32 || r[i] == 10 || r[i] == 9 || r[i] == 9 {
			continue
		}

		distR = append(distR, r[i])
	}
	dist = string(distR)
	return
}
//去除文字开头的问号
func TrimQusetionAtBegin(src string) (dist string) {
	src = strings.Trim(src, "\u3000")
	if len(src) == 0 {
		return 
	}
	r, _ := []rune(src), []rune{}
	if r[0] == 65311 || r[0] == 63 {
		dist,_ = SubStringBetween(src, 1, len(src))
	}else{
		dist = src
	}
	return
}

// IsNotEmpty checks if the address of point cs is empty
func IsNotEmpty(cs string) bool {
	return len(cs) > 0
}

// IsAnyEmpty checks if any of the css are empty or nil point
func IsAnyEmpty(css ...string) bool {

	if len(css) == 0 {
		return false
	}

	for _, str := range css {
		if IsEmpty(str) {
			return true
		}
	}

	return false
}

// IsNoneEmpty checks if none of the css are empty or nil point
func IsNoneEmpty(css ...string) bool {
	return !IsAnyEmpty(css...)
}

// IsAllEmpty checks if all of the css are empty or nil point.
func IsAllEmpty(css ...string) bool {

	if len(css) == 0 {
		return true
	}

	for _, str := range css {
		if IsNotEmpty(str) {
			return false
		}
	}

	return true
}

// IsBlank checks if a cs is empty, or nil point or whitespace only.
func IsBlank(cs string) bool {

	if strLen := len(cs); strLen == 0 {
		return true
	} else {
		return len(strings.TrimSpace(cs)) != strLen
	}

}

// IsNotBlank checks if a cs is not empty, not nill and not whitespace only.
func IsNotBlank(cs string) bool {
	return !IsBlank(cs)
}

// IsAnyBlank checks if any of the css are empty or nill or whitespace only
func IsAnyBlank(css ...string) bool {
	if len(css) == 0 {
		return false
	}

	for _, str := range css {
		if IsBlank(str) {
			return true
		}
	}

	return false
}

// StringsHasOneEmpty checks the input strings has someone empty.
func StringsHasOneEmpty(in ...string) bool {
	for _, str := range in {
		if len(str) == 0 {
			return true
		}
	}
	return false
}

//  endregion

// region Truncate

// Truncate truncates a String
func Truncate(str string, maxWidth int) string {
	ret, _ := TruncateFromWithMaxWith(str, 0, maxWidth)
	return ret
}

// TruncateFromWithMaxWith truncates a String
func TruncateFromWithMaxWith(str string, offset, maxWidth int) (ret string, err error) {
	if offset < 0 {
		return "", fmt.Errorf("offset cannot be negative")
	}

	if maxWidth < 0 {
		return "", fmt.Errorf("maxWidth cannot be negative")
	}

	if IsEmpty(str) {
		return str, nil
	}

	l := RuneLen(str)

	if offset > l {
		return EMPTY, nil
	}

	if l > maxWidth {

		var ix int
		if offset+maxWidth > l {
			ix = l
		} else {
			ix = offset + maxWidth
		}

		return SubStringBetween(str, offset, ix)
	}

	return SubString(str, offset)
}

// SubString returns a new string that is a substring of this string
// beginIndex means the string returned is the substring begins it and extends to the end of input“
func SubString(str string, beginIndex int) (ret string, err error) {

	if beginIndex < 0 {
		return "", fmt.Errorf("beginIndex cannot be negative")
	}

	subLen := RuneLen(str) - beginIndex

	if subLen <= 0 {
		return "", fmt.Errorf("beginIndex out of bound")
	}

	str2 := []rune(str)

	return string(str2[beginIndex:]), nil

}

// SubStringBetween returns a new string that is a substring of this string
func SubStringBetween(str string, beginIndex, endIndex int) (ret string, err error) {

	if beginIndex < 0 {
		return "", fmt.Errorf("beginIndex cannot be negative")
	}

	l := RuneLen(str)

	if endIndex > l {
		return "", fmt.Errorf("endIndex out of bound")
	}

	subLen := endIndex - beginIndex
	if subLen < 0 {
		return "", fmt.Errorf("endIndex must be bigger than beginIndex")
	}

	if beginIndex == 0 && endIndex == l {
		return str, nil
	} else {
		str2 := []rune(str)
		return string(str2[beginIndex:endIndex]), nil
	}
}

// SubstringBefore gets the substring before the first occurrence of a separator.
// stringutils.SubstringBefore("abc", "a")   = ""
// stringutils.SubstringBefore("abcba", "b") = "a"
// stringutils.SubstringBefore("abc", "c")   = "ab"
// stringutils.SubstringBefore("abc", "d")   = "abc"
// stringutils.SubstringBefore("abc", "")    = ""
func SubstringBefore(str, separator string) string {

	if IsEmpty(str) {
		return str
	}

	if IsEmpty(separator) {
		return EMPTY
	}

	pos := IndexOf(str, separator)

	if pos == IndexNotFound {
		return str
	}

	ret, _ := SubStringBetween(str, 0, pos)

	return ret
}

// SubstringAfter gets the substring after the first occurrence of a separator.
// The separator is not returned.
// stringutils.SubstringAfter("abc", "a")   = "bc"
// stringutils.SubstringAfter("abcba", "b") = "cba"
// stringutils.SubstringAfter("abc", "c")   = ""
// stringutils.SubstringAfter("abc", "d")   = ""
// stringutils.SubstringAfter("abc", "")    = "abc"
func SubstringAfter(str, separator string) string {

	if IsEmpty(str) {
		return str
	}

	spL := RuneLen(separator)
	if spL == 0 {
		return str
	}

	pos := IndexOf(str, separator)
	if pos == IndexNotFound {
		return EMPTY
	}

	ret, _ := SubString(str, pos+spL)

	return ret
}

// endregion

// region

// Strip strips whitespace from the start and end of a String
// stringutils.Strip("")       = ""
// stringutils.Strip("   ")    = ""
// stringutils.Strip("abc")    = "abc"
// stringutils.Strip("  abc")  = "abc"
// stringutils.Strip("abc  ")  = "abc"
// stringutils.Strip(" abc ")  = "abc"
// stringutils.Strip(" ab c ") = "ab c"
func Strip(str string) string {
	return StripWithChar(str, " ")
}

// StripWithChar strips any of a set of characters from the start and end of a String.
// This is similar to trim but allows the characters
// to be stripped to be controlled
// stringutils.StripWithChar("", *)            = ""
// stringutils.StripWithChar("  abcyx", "xyz") = "  abc"
func StripWithChar(str, stripChars string) string {
	if IsEmpty(str) {
		return str
	}

	str = StripStart(str, stripChars)
	return StripEnd(str, stripChars)
}

// StripStart strips any of a set of characters from the start of a String.
// stringutils.StripStart("", *)            = ""
// stringutils.StripStart("abc", "")        = "abc"
// stringutils.StripStart("yxabc  ", "xyz") = "abc  "
func StripStart(str, stripChars string) string {
	str2 := []rune(str)
	l := len(str2)
	if l == 0 {
		return str
	}

	if IsEmpty(stripChars) {
		return str
	}

	start := 0
	var ch string
	for start != l {
		ch, _ = CharAt(str, start)
		if IndexOf(stripChars, ch) == -1 {
			break
		}
		start++
	}

	ret, _ := SubString(str, start)
	return ret

}

// StripEnd strips any of a set of characters from the end of a String
//
// An empty string ("") input returns the empty string
// stringutils.StripEnd("", *)            = ""
// stringutils.StripEnd("abc", "")        = "abc"
// stringutils.StripEnd("  abcyx", "xyz") = "  abc"
// stringutils.StripEnd("120.00", ".0")   = "12"
func StripEnd(str, stripChars string) string {
	str2 := []rune(str)
	end := len(str2)
	if end == 0 {
		return str
	}

	if IsEmpty(stripChars) {
		return str
	}

	var ch string
	for end != 0 {
		ch, _ = CharAt(str, end-1)
		if IndexOf(stripChars, ch) == -1 {
			break
		}
		end--
	}

	ret, _ := SubStringBetween(str, 0, end)
	return ret
}

// endregion

// region Equals

// EqualsIgnoreCase compares two CharSequences, returning true if they represent
// equal sequences of characters, ignoring case.
func EqualsIgnoreCase(str1, str2 string) bool {
	if str1 == str2 {
		return true
	} else if len(str1) != len(str2) {
		return false
	} else {
		return RegionMatches(str1, true, 0, str2, 0, len(str2))
	}
}

// endregion

// region compare

// Compare compares two strings lexicographically.
// The comparison is based on the Unicode value of each character in the strings.
// The result is a negative integer if this str1 lexicographically precedes the argument str2.
// The result is a positive integer if this str1 lexicographically follows the argument str2.
// The result is zero if the strings are equal;
func Compare(str1, str2 string) int {
	l1, l2 := len(str1), len(str2)

	lim := mathutils.Min(l1, l2)

	for k := 0; k < lim; k++ {
		c1, _ := CharAt(str1, k)
		c2, _ := CharAt(str2, k)

		cu1, _ := utf8.DecodeRuneInString(c1)
		cu2, _ := utf8.DecodeRuneInString(c2)

		if cu1 != cu2 {
			return int(cu1 - cu2)
		}
	}

	return l1 - l2
}

// CompareIgnoreCase compares two Strings lexicographically, ignoring case differences.
// returning:
//
//	= 0, if str1 is equal to str2 (or both {@code null})
//	< 0, if str1 is less than str2
//	> 0, if str1 is greater than str2
func CompareIgnoreCase(str1, str2 string) int {

	if str1 == str2 {
		return 0
	}

	return Compare(strings.ToUpper(str1), strings.ToUpper(str2))
}

// EqualsAny compares given str1 to a char vararg of searchStrings,
// returning true if the str1 is equal to any of the searchStrings.
func EqualsAny(str1 string, searchStrings ...string) bool {
	if len(searchStrings) > 0 {
		for _, v := range searchStrings {
			if str1 == v {
				return true
			}
		}
	}
	return false
}

// EqualsAnyIgnoreCase compares given str1 to a char vararg of searchStrings
// returning true if the str1 is equal to any of the searchStrings, ignoring case.</p>
func EqualsAnyIgnoreCase(str1 string, searchStrings ...string) bool {
	if len(searchStrings) > 0 {
		for _, v := range searchStrings {
			if strings.ToUpper(str1) == strings.ToUpper(v) {
				return true
			}
		}
	}
	return false
}

// endregion

// region indexof

// IndexOf returns the index within this string of the first occurrence of
// the specified character.
func IndexOf(str, sub string) int {

	str2 := []rune(str)

	sub2 := []rune(sub)
	l2 := len(sub2)

outer:
	for i, s := range str2 {

		for j, su := range sub2 {

			if su == s {

				if j+1 == l2 {
					return i
				}

				continue
			}

			if j+1 == l2 {
				continue outer
			}

			break
		}
	}

	return IndexNotFound
}

// IndexOfFromIndex return the index within this string of the first occurrence of the
// specified substring.
// fromIndex is the index to begin searching from
func IndexOfFromIndex(str, sub string, fromIndex int) int {

	str2 := []rune(str)
	l := len(str2)

	sub2 := []rune(sub)
	l2 := len(sub2)

	return indexOf(str, 0, l, sub, 0, l2, fromIndex)
}

// IndexOfAny searchs a char to find the first index of any char in the given set of searchChars
// stringutils.IndexOfAny("", *)                  = -1
// stringutils.IndexOfAny(*, [])                  = -1
// stringutils.IndexOfAny("zzabyycdxx",['z','a']) = 0
// stringutils.IndexOfAny("zzabyycdxx",['b','y']) = 3
// stringutils.IndexOfAny("aba", ['z'])           = -1
func IndexOfAny(cs string, searchChars ...string) int {

	csLen, searchLen := len(cs), len(searchChars)

	for i := 0; i < csLen; i++ {
		ch, _ := CharAt(cs, i)

		for j := 0; j < searchLen; j++ {
			if searchChars[j] == ch {
				return i
			}
		}
	}

	return IndexNotFound
}

// source       the characters being searched.
// sourceOffset offset of the source string.
// sourceCount  count of the source string.
// target       the characters being searched for.
// targetOffset offset of the target string.
// targetCount  count of the target string.
// fromIndex    the index to begin searching from.
func indexOf(source string, sourceOffset, sourceCount int, target string, targetOffset, targetCount, fromIndex int) int {

	if fromIndex >= sourceCount {
		if targetCount == 0 {
			return sourceCount
		}
		return -1
	}

	if fromIndex < 0 {
		fromIndex = 0
	}

	if targetCount == 0 {
		return fromIndex
	}

	first, _ := CharAt(target, targetOffset)
	max := sourceOffset + (sourceCount - targetCount)

	for i := sourceOffset + fromIndex; i <= max; i++ {

		sI, _ := CharAt(source, i)
		/* Look for first character. */

		for i <= max && sI != first {
			i++
			sI, _ = CharAt(source, i)
		}

		/* Found first character, now look at the rest of v2 */
		if i <= max {
			j := i + 1
			end := j + targetCount - 1
			sJ, _ := CharAt(source, j)

			k := targetOffset + 1
			tk, _ := CharAt(target, k)
			for j < end && sJ == tk {
				j++
				k++
				sJ, _ = CharAt(source, j)
				tk, _ = CharAt(target, k)
			}

			if j == end {
				return i - sourceOffset
			}
		}
	}
	return -1
}

func lastIndexOf(source string, sourceOffset, sourceCount int, target string, targetOffset, targetCount, fromIndex int) int {

	rightIndex := sourceCount - targetCount

	if fromIndex < 0 {
		return -1
	}

	if fromIndex > rightIndex {
		fromIndex = rightIndex
	}

	if targetCount == 0 {
		return fromIndex
	}

	strLastIndex := targetOffset + targetCount - 1
	strLastChar, _ := CharAt(target, strLastIndex)
	min := sourceOffset + targetCount - 1
	i := min + fromIndex

startSearchForLastChar:
	{

		for {
			sI, _ := CharAt(source, i)
			for i >= min && sI != strLastChar {
				i--
				sI, _ = CharAt(source, i)
			}

			if i < min {
				return -1
			}

			j := i - 1
			start := j - (targetCount - 1)
			k := strLastIndex - 1

			for j > start {

				sJ, _ := CharAt(source, j)
				j--
				tK, _ := CharAt(target, k)
				k--

				if sJ != tK {
					i--
					goto startSearchForLastChar
				}

			}

			return start - sourceOffset + 1
		}
	}

}

// OrdinalIndexOf finds the n-th index within searchStr
// The code starts looking for a match at the start of the target,
// incrementing the starting index by one after each successful match
// (unless searchStr is an empty string in which case the position
// is never incremented and '0' is returned immediately).
// This means that matches may overlap.
// stringutils.OrdinalIndexOf("ababab","aba", 1)   = 0
// stringutils.OrdinalIndexOf("ababab","aba", 2)   = 2
// stringutils.OrdinalIndexOf("ababab","aba", 3)   = -1
// stringutils.OrdinalIndexOf("abababab", "abab", 1) = 0
// stringutils.OrdinalIndexOf("abababab", "abab", 2) = 2
// stringutils.OrdinalIndexOf("abababab", "abab", 3) = 4
// stringutils.OrdinalIndexOf("abababab", "abab", 4) = -1
func OrdinalIndexOf(str, searchStr string, ordinal int, lastIndex bool) int {

	if ordinal <= 0 {
		return IndexNotFound
	}

	l1 := len(str)
	l2 := len(searchStr)

	if l2 == 0 {

		if lastIndex {
			return l1
		}

		return 0
	}

	found := 0
	index := 0
	if lastIndex {
		index = l1
	} else {
		index = IndexNotFound
	}

	for found < ordinal {

		if lastIndex {
			index = LastIndexOf(str, searchStr, index-1)
		} else {
			index = IndexOfFromIndex(str, searchStr, index+1)
		}

		if index < 0 {
			return index
		}

		found++
	}

	return index
}

// LastIndexOf returns the index within this string of the last occurrence of the
// specified substring, searching backward starting at the specified index.
// important: The search starts at the startPos and works backwards; matches starting after the start position are ignored.
func LastIndexOf(cs, searchChar string, startPos int) int {
	str2 := []rune(cs)
	l1 := len(str2)

	sub2 := []rune(searchChar)
	l2 := len(sub2)

	return lastIndexOf(cs, 0, l1, searchChar, 0, l2, startPos)
}

// IndexOfAnyBut searches a CharSequence to find the first index of any character not in the given set of characters.
// stringutils.IndexOfAnyBut("", *)                         = -1
// stringutils.IndexOfAnyBut(*, [])                         = -1
// stringutils.IndexOfAnyBut("zzabyycdxx", ['z', 'a'] )     = 3
// stringutils.IndexOfAnyBut("aba",  ['z'] )                = 0
// stringutils.IndexOfAnyBut("aba",  ['a', 'b'])            = -1
func IndexOfAnyBut(cs string, searcChars ...string) int {

	csL, searchL := len(cs), len(searcChars)

outer:

	for i := 0; i < csL; i++ {
		ch, _ := CharAt(cs, i)
		for j := 0; j < searchL; j++ {
			if searcChars[j] == ch {
				continue outer
			}
		}
		return i

	}
	return IndexNotFound
}

// LastIndexOfAny find the latest index of any of a set of potential substrings.
// !! something is different that searching a array containing "" do NOT return the length of the searched array, but the Apache StringUtils will return the length
// stringutils.LastIndexOfAny(*, [])                     = -1
// stringutils.LastIndexOfAny("zzabyycdxx", ["ab","cd"]) = 6
// stringutils.LastIndexOfAny("zzabyycdxx", ["cd","ab"]) = 6
// stringutils.LastIndexOfAny("zzabyycdxx", ["mn","op"]) = -1
// stringutils.LastIndexOfAny("zzabyycdxx", ["mn","op"]) = -1
// stringutils.LastIndexOfAny("zzabyycdxx", ["mn",""])   = 10
func LastIndexOfAny(cs string, searcChars ...string) int {

	csL := len(cs)
	if csL == 0 || len(searcChars) == 0 {
		return IndexNotFound
	}

	ret := IndexNotFound
	tmp := 0

	for _, search := range searcChars {

		if len(search) == 0 {
			continue
		}

		tmp = LastIndexOf(cs, search, csL)
		if tmp > ret {
			ret = tmp
		}
	}

	return ret
}

// endregion

// region Contains

// Contains checks if cs contains a search character sub
// stringutils.Contains("", "")      = true
// stringutils.Contains("abc", "")   = true
// stringutils.Contains("abc", "a")  = true
// stringutils.Contains("abc", "z")  = false
func Contains(cs, sub string) bool {
	if IsEmpty(cs) && IsEmpty(sub) {
		return true
	}

	return IndexOfFromIndex(cs, sub, 0) >= 0
}

// ContainsIgnoreCase checks if str contains a searchStr irrespective of case
// stringutils.ContainsIgnoreCase("", "") = true
// stringutils.ContainsIgnoreCase("abc", "") = true
// stringutils.ContainsIgnoreCase("abc", "a") = true
// stringutils.ContainsIgnoreCase("abc", "z") = false
// stringutils.ContainsIgnoreCase("abc", "A") = true
// stringutils.ContainsIgnoreCase("abc", "Z") = false
func ContainsIgnoreCase(str, searchStr string) bool {

	if IsEmpty(str) && IsEmpty(searchStr) {
		return true
	}

	l1 := RuneLen(searchStr)
	max := RuneLen(str) - l1

	for i := 0; i <= max; i++ {
		if RegionMatches(str, true, i, searchStr, 0, l1) {
			return true
		}
	}

	return false
}

// ContainsWhitespace checks whether the given CharSequence contains any whitespace characters.
func ContainsWhitespace(str string) bool {

	l := len(str)
	if l == 0 {
		return false
	}

	for i := 0; i < l; i++ {
		c, _ := CharAt(str, i)
		if IsWhitespace(c) {
			return true
		}
	}

	return false
}

// ContainsAny checks if the CharSequence contains any character in the given set of characters.
// stringutils.ContainsAny("", *)                  = false
// stringutils.ContainsAny(*, [])                  = false
// stringutils.ContainsAny("zzabyycdxx",['z','a']) = true
// stringutils.ContainsAny("zzabyycdxx",['b','y']) = true
// stringutils.ContainsAny("zzabyycdxx",['z','y']) = true
// stringutils.ContainsAny("aba", ['z'])           = false
func ContainsAny(cs string, searchChars ...string) bool {

	csL, searchL := len(cs), len(searchChars)

	for i := 0; i < csL; i++ {
		ch, _ := CharAt(cs, i)
		for j := 0; j < searchL; j++ {
			if searchChars[j] == ch {
				return true
			}
		}
	}

	return false
}

// ContainsOnly checks if the CharSequence contains only certain characters.
// stringutils.ContainsOnly("", *)         = true
// stringutils.ContainsOnly("ab", "")      = false
// stringutils.ContainsOnly("abab", "a", "b", "c") = true
// stringutils.ContainsOnly("ab1", "a", "b", "c")  = false
// stringutils.ContainsOnly("abz", "a", "b", "c")  = false
func ContainsOnly(cs string, valid ...string) bool {

	if len(cs) == 0 {
		return true
	}
	if len(valid) == 0 {
		return false
	}
	return IndexOfAnyBut(cs, valid...) == IndexNotFound
}

// ContainsNone checks that the CharSequence does not contain certain characters.
// stringutils.ContainsNone("", *)         = true
// stringutils.ContainsNone("ab", "")      = true
// stringutils.ContainsNone("abab", "x", "y", "z") = true
// stringutils.ContainsNone("ab1", "x", "y", "z")  = true
// stringutils.ContainsNone("abz", "x", "y", "z")  = false
func ContainsNone(cs string, searchChars ...string) bool {

	csL, searchL := len(cs), len(searchChars)

	for i := 0; i < csL; i++ {
		ch, _ := CharAt(cs, i)

		for j := 0; j < searchL; j++ {
			if searchChars[j] == ch {
				return false
			}
		}
	}

	return true
}

// endregion

// region Left/Right/Mid

// Left gets the leftmost l characters of a String.
// stringutils.Left("abc", 0)   = ""
// stringutils.Left("abc", 2)   = "ab"
// stringutils.Left("abc", 4)   = "abc"
func Left(str string, l int) string {

	sl := RuneLen(str)

	if sl == 0 || l < 0 {
		return EMPTY
	}

	if sl <= l {
		return str
	}

	ret, _ := SubStringBetween(str, 0, l)

	return ret
}

// Right gets the rightmost r(length) characters of a String.
// stringutils.Right("abc", 0)   = ""
// stringutils.Right("abc", 2)   = "bc"
// stringutils.Right("abc", 4)   = "abc"
func Right(str string, r int) string {
	sl := RuneLen(str)

	if sl == 0 || r < 0 {
		return EMPTY
	}

	if sl <= r {
		return str
	}

	ret, _ := SubString(str, sl-r)

	return ret
}

// Mid gets length l characters from the middle of a String.
// stringutils.Mid("abc", 0, 2)   = "ab"
// stringutils.Mid("abc", 0, 4)   = "abc"
// stringutils.Mid("abc", 2, 4)   = "c"
// stringutils.Mid("abc", 4, 2)   = ""
// stringutils.Mid("abc", -2, 2)  = "ab"
func Mid(str string, pos, l int) string {

	sl := RuneLen(str)

	if l < 0 || pos > sl {
		return EMPTY
	}

	if pos < 0 {
		pos = 0
	}

	if sl <= pos+l {
		ret, _ := SubString(str, pos)
		return ret
	}

	ret, _ := SubStringBetween(str, pos, pos+l)

	return ret
}

// endregion

// region Numeric

// IsNumeric Checks if the str string contains only Unicode digits.
// A decimal point is not a Unicode digit and returns false.</p>
// StringUtils.IsNumeric("")     = false
// StringUtils.IsNumeric("  ")   = false
// StringUtils.IsNumeric("123")  = true
// StringUtils.IsNumeric("12 3") = false
// StringUtils.IsNumeric("ab2c") = false
// StringUtils.IsNumeric("12-3") = false
// StringUtils.IsNumeric("12.3") = false
// StringUtils.IsNumeric("-123") = false
// StringUtils.IsNumeric("+123") = false
func IsNumeric(str string) bool {
	if IsEmpty(str) {
		return false
	}
	sz := RuneLen(str)
	for i := 0; i < sz; i++ {
		c, _ := CharAt(str, i)
		if !unicode.IsDigit([]rune(c)[0]) {
			return false
		}
	}

	return true
}

// ToInt converts str to int
func ToInt32(str string) (ret int32, err error) {
	ret2, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return
	}
	return int32(ret2), nil
}

// ToInt64 converts str to int64
func ToInt64(str string) (ret int64, err error) {
	ret, err = strconv.ParseInt(str, 10, 64)
	return
}

// FromInt64 converts int64 to str
func FromInt64(in int64) (ret string) {
	return strconv.FormatInt(in, 16)
}

// endregion

// region base64

// ToBase64 encodes the str in to new base64 string.
func ToBase64(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

// FromBase64 decodes the str in encoded to original string.
func FromBase64(in string) (str string, err error) {

	decoded, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}

	return string(decoded), nil
}

// endregion base64

/* region trim*/

// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
func TrimSpace(in string) string {
	return strings.TrimSpace(in)
}

/* endregion */

// IDArrayToSQLInString combiles id nums to a sql "in" clause string
func IDArrayToSQLInString(in []int) string {
	if len(in) > 0 {
		ret := strings.Join(IntArrayToStringArray(in), "','")
		return "'" + ret + "'"
	}
	return ""
}

// IntArrayToStringArray converts int an array to new string array
func IntArrayToStringArray(in []int) []string {
	ret := make([]string, 0, len(in))
	for _, v := range in {
		ret = append(ret, strconv.Itoa(v))
	}
	return ret
}

// IsWhitespace returns if the str is a whitespace of latin-1
//
//	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP) are all true
func IsWhitespace(str string) bool {

	l := len(str)
	if l != 1 {
		return false
	}

	return unicode.IsSpace([]rune(str)[0])
}

// RegionMatches tests if two string regions are equal.
// cs the char to be processed
// ignoreCase whether or not to be case insensitive
// thisStart the index to start on the param cs
// substring the String/subString to be looked for
// start the index to start on the subString
// length character length of the region
// whether the region matched
func RegionMatches(cs string, ignoreCase bool, thisStart int,
	substring string, start int, length int) bool {

	index1 := thisStart
	index2 := start

	srcLen := RuneLen(cs) - thisStart
	otherLen := RuneLen(substring) - start

	ta := ToCharArray(cs)
	pa := ToCharArray(substring)

	// Check for invalid parameters
	if thisStart < 0 || start < 0 || length < 0 {
		return false
	}

	// Check that the regions are long enough
	if srcLen < length || otherLen < length {
		return false
	}

	for ; length > 0; length-- {

		c1 := ta[index1]
		index1++
		c2 := pa[index2]
		index2++

		if c1 == c2 {
			continue
		}
		if ignoreCase {

			u1 := strings.ToUpper(c1)
			u2 := strings.ToUpper(c2)
			if u1 == u2 {
				continue
			}

			if strings.ToLower(u1) == strings.ToLower(u2) {
				continue
			}
		}
		return false
	}

	return true
}

// CharAt returns the char value at the specified index.
func CharAt(str string, index int) (ret string, err error) {

	if index < 0 || index >= len(str) {
		return "", fmt.Errorf("%d index out of bound %d", index, len(str))
	}

	str2 := []rune(str)

	return string(str2[index : index+1]), nil
}

// ToCharArray returns a char array contains all string chars;
func ToCharArray(str string) []string {
	l := len(str)
	ret := make([]string, 0, l)
	for _, r := range str {
		c := string(r)
		ret = append(ret, c)
	}

	return ret
}

// RuneLen returns length of str's rune array
func RuneLen(str string) int {
	return len([]rune(str))
}

// StringChinesePhoneNumOrEmail phone all email
// if return 2 means 'in' is a phone num, 1 means it is a email address
func StringChinesePhoneNumOrEmail(in string) int {
	reg, _ := regexp.Compile(`^1[345789]\d{9}$`)
	if reg.MatchString(in) {
		return 2
	}

	reg = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if reg.MatchString(in) {
		return 1
	}

	return 0
}

// ToString changes the string "in" to a string.
func ToString(in interface{}) string {

	if v, ok := in.(string); ok {
		return v
	}

	return fmt.Sprintf("%v", in)
}

// region join

// Join joins the elements of the provided array into a single String containing the provided list of elements
// stringutils.Join([], *)                 = ""
// stringutils.Join(nil, *)             = ""
// stringutils.Join([1, 2, 3], ';')  = "1;2;3"
// stringutils.Join([1, 2, 3], "") = "123"
func Join(in []string, separator string) string {
	if in == nil {
		return ""
	}

	return JoinBetween(in, separator, 0, len(in))
}

// JoinBetween joins the elements of the provided array into a single String containing the provided list of elements.
// stringutils.JoinBetween([], *)              = ""
// stringutils.JoinBetween(nil, *)             = ""
// stringutils.JoinBetween([1, 2, 3], ';')  = "1;2;3"
// stringutils.JoinBetween([1, 2, 3], "") = "123"
func JoinBetween(in []string, separator string, startIndex, endIndex int) string {
	if in == nil {
		return ""
	}

	noOfItems := endIndex - startIndex

	if noOfItems <= 0 {
		return EMPTY
	}

	var builder strings.Builder

	for i := startIndex; i < endIndex; i++ {
		if i > startIndex {
			builder.WriteString(separator)
		}
		builder.WriteString(in[i])
	}
	return builder.String()
}

// endregion

// region startWith
// StartsWith checks if a CharSequence starts with a specified prefix.
// stringutils.StartsWith("", "")      = true
// stringutilu.StartsWith("", "abc")     = false
// stringutilss.StartsWith("abcdef", "")  = false
// stringutils.StartsWith("abcdef", "abc") = true
// stringutils.StartsWith("ABCDEF", "abc") = false
func StartsWith(str, prefix string) bool {
	return StartsWithIgnoreCase(str, prefix, false)
}

func StartsWithIgnoreCase(str, prefix string, ignoreCase bool) bool {
	if str == "" || prefix == "" {
		return prefix == ""
	}
	if len(prefix) > len(str) {
		return false
	}
	return RegionMatches(str, ignoreCase, 0, prefix, 0, RuneLen(prefix))
}

// endregion

// region split

// Split slices s into all substrings separated by sep and returns a slice of
// the substrings between those separators.
func Split(str, sep string) []string {
	return strings.Split(str, sep)
}

// endregion

// region Delete

// DeleteWhitespace deletes all whitespaces from a String
// stringutils.DeleteWhitespace("")           = ""
// stringutils.DeleteWhitespace("abc")        = "abc"
// stringutils.DeleteWhitespace("   ab  c  ") = "abc"
func DeleteWhitespace(str string) string {
	if str == "" {
		return ""
	}

	sz := RuneLen(str)
	var builder strings.Builder

	count := 0
	for i := 0; i < sz; i++ {
		charT, _ := CharAt(str, i)
		if !IsWhitespace(charT) {
			builder.WriteString(charT)
		}
	}
	if count == sz {
		return str
	}
	return builder.String()
}

// RemoveStart removes a substring only if it is at the beginning of a source string,
// otherwise returns the source string.
// case when len(remove) is greater than len(str), then str will be returned
// stringutils.RemoveStart("", *)        = ""
// stringutils.RemoveStart("www.domain.com", "www.")   = "domain.com"
// stringutils.RemoveStart("domain.com", "www.")       = "domain.com"
// stringutils.RemoveStart("www.domain.com", "domain") = "www.domain.com"
// stringutils.RemoveStart("abc", "")    = "abc"// s
func RemoveStart(str, remove string) string {

	if str == "" || remove == "" || len(remove) > len(str) {
		return str
	}

	if StartsWith(str, remove) {
		ret, _ := SubString(str, RuneLen(remove))
		return ret
	}

	return str
}

// RemoveStartIgnoreCase cases insensitive removal of a substring if it is at the beginning of a source string,
// otherwise returns the source string
// case when len(remove) is greater than len(str), then str will be returned
//
// stringUtils.RemoveStartIgnoreCase("", *)        = ""
// stringUtils.RemoveStartIgnoreCase("www.domain.com", "www.")   = "domain.com"
// stringUtils.RemoveStartIgnoreCase("www.domain.com", "WWW.")   = "domain.com"
// stringUtils.RemoveStartIgnoreCase("domain.com", "www.")       = "domain.com"
// stringUtils.RemoveStartIgnoreCase("www.domain.com", "domain") = "www.domain.com"
// stringUtils.RemoveStartIgnoreCase("abc", "")    = "abc"// s
func RemoveStartIgnoreCase(str, remove string) string {
	if str == "" || remove == "" || len(remove) > len(str) {
		return str
	}

	if StartsWithIgnoreCase(str, remove, true) {
		ret, _ := SubString(str, RuneLen(remove))
		return ret
	}

	return str
}

// RemoveEnd removes a substring only if it is at the end of a source string,
// otherwise returns the source string

// case when len(remove) is greater than len(str), then str will be returned

// stringutils.RemoveEnd("", *)        = ""
// stringutils.RemoveEnd("www.domain.com", ".com.")  = "www.domain.com"
// stringutils.RemoveEnd("www.domain.com", ".com")   = "www.domain"
// stringutils.RemoveEnd("www.domain.com", "domain") = "www.domain.com"
// stringutils.RemoveEnd("abc", "")    = "abc"
func RemoveEnd(str, remove string) string {

	if str == "" || remove == "" || len(remove) > len(str) {
		return str
	}

	if EndsWith(str, remove) {
		ret, _ := SubStringBetween(str, 0, len(str)-len(remove))
		return ret
	}

	return str
}

// RemoveEndIgnoreCase cases insensitive removal of a substring if it is at the end of a source string,
// otherwise returns the source string.</p>

// case when len(remove) is greater than len(str), then str will be returned

// stringutils.RemoveEndIgnoreCase("", *)        = ""
// stringutils.RemoveEndIgnoreCase("www.domain.com", ".com.")  = "www.domain.com"
// stringutils.RemoveEndIgnoreCase("www.domain.com", ".com")   = "www.domain"
// stringutils.RemoveEndIgnoreCase("www.domain.com", "domain") = "www.domain.com"
// stringutils.RemoveEndIgnoreCase("abc", "")    = "abc"
// stringutils.RemoveEndIgnoreCase("www.domain.com", ".COM") = "www.domain")
// stringutils.RemoveEndIgnoreCase("www.domain.COM", ".com") = "www.domain")
func RemoveEndIgnoreCase(str, remove string) string {

	if str == "" || remove == "" || len(remove) > len(str) {
		return str
	}

	if EndsWithIgnoreCase(str, remove) {
		ret, _ := SubStringBetween(str, 0, len(str)-len(remove))
		return ret
	}

	return str
}

// endregion

// region endWith

// EndsWith checks if a CharSequence ends with a specified suffix.</p>
// The comparison is case sensitive.

// stringutils.EndsWith("", "def")     = false
// stringutils.EndsWith("abcdef", "def") = true
// stringutils.EndsWith("ABCDEF", "def") = false
// stringutils.EndsWith("ABCDEF", "cde") = false
// stringutils.EndsWith("ABCDEF", "")    = true
func EndsWith(str, suffix string) bool {
	return endsWith(str, suffix, false)
}

// EndsWithIgnoreCase case insensitive checks if a CharSequence ends with a specified suffix.
//
// The comparison is case insensitive.
// stringutils.EndsWithIgnoreCase("", "def")     = false
// stringutils.EndsWithIgnoreCase("abcdef", "")  = true
// stringutils.EndsWithIgnoreCase("abcdef", "def") = true
// stringutils.EndsWithIgnoreCase("ABCDEF", "def") = true
// stringutils.EndsWithIgnoreCase("ABCDEF", "cde") = false
func EndsWithIgnoreCase(str, suffix string) bool {
	return endsWith(str, suffix, true)
}

// EndsWithAny checks if a CharSequence ends with any of the provided case-sensitive suffixes
// stringutils.EndsWithAny("")      = false
// stringutils.EndsWithAny("abcxyz", "") = true
// stringutils.EndsWithAny("abcxyz", "xyz") = true
// stringutils.EndsWithAny("abcxyz", "xyz", "abc") = true
// stringutils.EndsWithAny("abcXYZ", "def", "XYZ") = true
// stringutils.EndsWithAny("abcXYZ", "def", "xyz") = false
func EndsWithAny(sequence string, searchStrings ...string) bool {
	if sequence == "" || len(searchStrings) == 0 {
		return false
	}

	for _, v := range searchStrings {
		if endsWith(sequence, v, false) {
			return true
		}
	}
	return false
}

// EndsWithAny checks if a CharSequence ends with any of the provided case-insensitive suffixes
// stringutils.EndsWithAny("")      = false
// stringutils.EndsWithAny("abcxyz", "") = true
// stringutils.EndsWithAny("abcxyz", "xyz") = true
// stringutils.EndsWithAny("abcxyz", "xyz", "abc") = true
// stringutils.EndsWithAny("abcXYZ", "def", "XYZ") = true
// stringutils.EndsWithAny("abcXYZ", "def", "xyz") = true
func EndsWithAnyIgnoreCase(sequence string, searchStrings ...string) bool {
	if sequence == "" || len(searchStrings) == 0 {
		return false
	}

	for _, v := range searchStrings {
		if endsWith(sequence, v, true) {
			return true
		}
	}
	return false
}

// endsWith checks if a string ends with a specified suffix (optionally case insensitive).
func endsWith(str, suffix string, ignoreCase bool) bool {
	if str == "" || suffix == "" {
		return suffix == ""
	}

	sufLen := RuneLen(suffix)
	strLen := RuneLen(str)

	if sufLen > strLen {
		return false
	}

	strOffset := strLen - sufLen
	return RegionMatches(str, ignoreCase, strOffset, suffix, 0, sufLen)
}

// endregion

// region append

// appendIfMissing appends the suffix to the end of the string if the string does not
// already end with the suffix which in suffixes
// A new String if suffix was appended, the same string otherwise.
//
// str The string.
// suffix The suffix to append to the end of the string.
// ignoreCase Indicates whether the compare should ignore case.
// suffixes Additional suffixes that are valid terminators (optional).
func appendIfMissing(str, suffix string, ignoreCase bool, suffixes ...string) string {
	if len(suffix) == 0 || endsWith(str, suffix, ignoreCase) {
		return str
	}

	if len(suffixes) > 0 {
		for _, s := range suffixes {
			if endsWith(str, s, ignoreCase) {
				return str
			}
		}
	}

	return str + suffix
}

// Appends the suffix to the end of the string if the string does not
// already end with any of the suffixes.
// stringutils.AppendIfMissing("abc", "") = "abc"
// stringutils.AppendIfMissing("", "xyz") = "xyz"
// stringutils.AppendIfMissing("abc", "xyz") = "abcxyz"
// stringutils.AppendIfMissing("abcxyz", "xyz") = "abcxyz"
// stringutils.AppendIfMissing("abcXYZ", "xyz") = "abcXYZxyz"
// stringutils.AppendIfMissing("abc", "xyz", "") = "abc"
// stringutils.AppendIfMissing("abc", "xyz", "mno") = "abcxyz"
// stringutils.AppendIfMissing("abcxyz", "xyz", "mno") = "abcxyz"
// stringutils.AppendIfMissing("abcmno", "xyz", "mno") = "abcmno"
// stringutils.AppendIfMissing("abcXYZ", "xyz", "mno") = "abcXYZxyz"
// stringutils.AppendIfMissing("abcMNO", "xyz", "mno") = "abcMNOxyz"
func AppendIfMissing(str, suffix string, suffixes ...string) string {
	return appendIfMissing(str, suffix, false, suffixes...)
}

// Appends the suffix to the end of the string if the string does not
// already end, case insensitive, with any of the suffixes.
//
// stringutils.AppendIfMissingIgnoreCase("", "xyz") = "xyz"
// stringutils.AppendIfMissingIgnoreCase("abc", "xyz") = "abcxyz"
// stringutils.AppendIfMissingIgnoreCase("abcxyz", "xyz") = "abcxyz"
// stringutils.AppendIfMissingIgnoreCase("abcXYZ", "xyz") = "abcXYZ"
// stringutils.AppendIfMissingIgnoreCase("abc", "xyz", "") = "abc"
// stringutils.AppendIfMissingIgnoreCase("abc", "xyz", "mno") = "abcxyz"
// stringutils.AppendIfMissingIgnoreCase("abcxyz", "xyz", "mno") = "abcxyz"
// stringutils.AppendIfMissingIgnoreCase("abcmno", "xyz", "mno") = "abcmno"
// stringutils.AppendIfMissingIgnoreCase("abcXYZ", "xyz", "mno") = "abcXYZ"
// stringutils.AppendIfMissingIgnoreCase("abcMNO", "xyz", "mno") = "abcMNO"
func AppendIfMissingIgnoreCase(str, suffix string, suffixes ...string) string {
	return appendIfMissing(str, suffix, true, suffixes...)
}

// endregion

// region prepend

// prependIfMissing prepends the prefix to the start of the string if the string does not
// already start with any of the prefixes.
//
// return A new String if prefix was prepended, the same string otherwise.
//
// str The string.
// prefix The prefix to prepend to the start of the string.
// ignoreCase Indicates whether the compare should ignore case.
// prefixes Additional prefixes that are valid (optional).
func prependIfMissing(str, prefix string, ignoreCase bool, prefixes ...string) string {
	if prefix == "" || StartsWithIgnoreCase(str, prefix, ignoreCase) {
		return str
	}

	if len(prefixes) > 0 {

		for _, s := range prefixes {
			if StartsWithIgnoreCase(str, s, ignoreCase) {
				return str
			}
		}

	}
	return prefix + str
}

// prependIfMissing Prepends the prefix to the start of the string if the string does not
// already start with any of the prefixes.
//
// A new String if prefix was prepended, the same string otherwise.

// stringutils.PrependIfMissing("", "xyz") = "xyz"
// stringutils.PrependIfMissing("abc", "xyz") = "xyzabc"
// stringutils.PrependIfMissing("xyzabc", "xyz") = "xyzabc"
// stringutils.PrependIfMissing("XYZabc", "xyz") = "xyzXYZabc"
// stringutils.PrependIfMissing("abc", "xyz", "") = "abc"
// stringutils.PrependIfMissing("abc", "xyz", "mno") = "xyzabc"
// stringutils.PrependIfMissing("xyzabc", "xyz", "mno") = "xyzabc"
// stringutils.PrependIfMissing("mnoabc", "xyz", "mno") = "mnoabc"
// stringutils.PrependIfMissing("XYZabc", "xyz", "mno") = "xyzXYZabc"
// stringutils.PrependIfMissing("MNOabc", "xyz", "mno") = "xyzMNOabc"

// str The string.
// prefix The prefix to prepend to the start of the string.
// prefixes Additional prefixes that are valid.
func PrependIfMissing(str, prefix string, prefixes ...string) string {
	return prependIfMissing(str, prefix, false, prefixes...)
}

/**
 * PrependIfMissingIgnoreCase prepends the prefix to the start of the string if the string does not
 * already start, case insensitive, with any of the prefixes.
 *
 * stringutils.PrependIfMissingIgnoreCase("", "xyz") = "xyz"
 * stringutils.PrependIfMissingIgnoreCase("abc", "xyz") = "xyzabc"
 * stringutils.PrependIfMissingIgnoreCase("xyzabc", "xyz") = "xyzabc"
 * stringutils.PrependIfMissingIgnoreCase("XYZabc", "xyz") = "XYZabc"
 * stringutils.PrependIfMissingIgnoreCase("abc", "xyz", "") = "abc"
 * stringutils.PrependIfMissingIgnoreCase("abc", "xyz", "mno") = "xyzabc"
 * stringutils.PrependIfMissingIgnoreCase("xyzabc", "xyz", "mno") = "xyzabc"
 * stringutils.PrependIfMissingIgnoreCase("mnoabc", "xyz", "mno") = "mnoabc"
 * stringutils.PrependIfMissingIgnoreCase("XYZabc", "xyz", "mno") = "XYZabc"
 * stringutils.PrependIfMissingIgnoreCase("MNOabc", "xyz", "mno") = "MNOabc"
 * str The string.
 * prefix The prefix to prepend to the start of the string.
 * prefixes Additional prefixes that are valid (optional).
 *
 * A new String if prefix was prepended, the same string otherwise.
 */
func PrependIfMissingIgnoreCase(str, prefix string, prefixes ...string) string {
	return prependIfMissing(str, prefix, true, prefixes...)
}

// endregion

// Wrap wraps a String with another String.
// stringutils.Wrap("", *)           = ""
// stringutils.Wrap("ab", "x")       = "xabx"
// stringutils.Wrap("ab", "\"")      = "\"ab\""
// stringutils.Wrap("\"ab\"", "\"")  = "\"\"ab\"\""
// stringutils.Wrap("ab", "'")       = "'ab'"
// stringutils.Wrap("'abcd'", "'")   = "”abcd”"
// stringutils.Wrap("\"abcd\"", "'") = "'\"abcd\"'"
// stringutils.Wrap("'abcd'", "\"")  = "\"'abcd'\""
func Wrap(str, wrapWith string) string {

	if str == "" || wrapWith == "" {
		return str
	}

	return wrapWith + str + wrapWith
}

// WrapIfMissing wraps a string with a string if that string is missing from the start or end of the given string.
// stringutils.WrapIfMissing("", *)           = ""
// stringutils.WrapIfMissing("ab", "x")       = "xabx"
// stringutils.WrapIfMissing("ab", "\"")      = "\"ab\""
// stringutils.WrapIfMissing("\"ab\"", "\"")  = "\"ab\""
// stringutils.WrapIfMissing("ab", "'")       = "'ab'"
// stringutils.WrapIfMissing("'abcd'", "'")   = "'abcd'"
// stringutils.WrapIfMissing("\"abcd\"", "'") = "'\"abcd\"'"
// stringutils.WrapIfMissing("'abcd'", "\"")  = "\"'abcd'\""
// stringutils.WrapIfMissing("/", "/")  = "/"
// stringutils.UrapIfMissing("a/b/c", "/")  = "/a/b/c/"
// stringutils.WrapIfMissing("/a/b/c", "/")  = "/a/b/c/"
// stringutils.WrapIfMissing("a/b/c/", "/")  = "/a/b/c/"
func WrapIfMissing(str, wrapWith string) string {
	if str == "" || wrapWith == "" {
		return str
	}

	var builder strings.Builder
	if !StartsWith(str, wrapWith) {
		builder.WriteString(wrapWith)
	}
	builder.WriteString(str)
	if !endsWith(str, wrapWith, false) {
		builder.WriteString(wrapWith)
	}
	return builder.String()
}

// Unwrap unwraps a given string from anther string.

// stringutils.Unwrap("\'abc\'", "\'")    = "abc"
// stringutils.Unwrap("\"abc\"", "\"")    = "abc"
// stringutils.Unwrap("AABabcBAA", "AA")  = "BabcB"
// stringutils.Unwrap("A", "#")           = "A"
// stringutils.Unwrap("#A", "#")          = "#A"
// stringutils.Unwrap("A#", "#")          = "A#"
func Unwrap(str, wrapToken string) string {
	if str == "" || wrapToken == "" {
		return str
	}

	if StartsWith(str, wrapToken) && endsWith(str, wrapToken, false) {
		startIndex := IndexOf(str, wrapToken)
		endIndex := LastIndexOf(str, wrapToken, RuneLen(str))
		wrapLength := RuneLen(wrapToken)
		if startIndex != -1 && endIndex != -1 {
			ret, _ := SubStringBetween(str, startIndex+wrapLength, endIndex)
			return ret
		}
	}

	return str
}

// region wrap

// endregion

// region rand

const (
	STR_RAND_KIND_NUM   = 0 // 纯数字
	STR_RAND_KIND_LOWER = 1 // 小写字母
	STR_RAND_KIND_UPPER = 2 // 大写字母
	STR_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// Rand returns a string that the length is size
func Rand(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}

	return string(result)
}

// endregion
