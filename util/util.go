package util

import (
  "os"
  "strings"
)

var ENCODING_UNENCODED_TOKENS = []string{"%", ":", "[", "]", ","}
var ENCODING_ENCODED_TOKENS = []string{"%25", "%3A", "%5B", "%5D", "%2C"}
var DECODING_UNENCODED_TOKENS = []string{":", "[", "]", ",", "%"}
var DECODING_ENCODED_TOKENS = []string{"%3A", "%5B", "%5D", "%2C", "%25"}

// fail if an error is provided and print out the message
func CheckForError(err error, message string) {
  if err != nil {
      println(message + ": ", err.Error())
      os.Exit(1)
  }
}

// simple http-ish encoding to handle special characters
func Encode(value string) (string) {
  return replace(ENCODING_UNENCODED_TOKENS, ENCODING_ENCODED_TOKENS, value)
}

// simple http-ish decoding to handle special characters
func Decode(value string) (string) {
  return replace(DECODING_ENCODED_TOKENS, DECODING_UNENCODED_TOKENS, value)
}

// replace the from tokens to the to tokens (both arrays must be the same length)
func replace(fromTokens []string, toTokens []string, value string) (string) {
  for i:=0; i<len(fromTokens); i++ {
      value = strings.Replace(value, fromTokens[i], toTokens[i], -1)
  }
  return value;
}