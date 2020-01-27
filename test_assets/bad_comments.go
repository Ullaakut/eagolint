package test_assets

// This file only contains comments that are meant to trigger warnings in
// eagolint.

// This  comment has a multi-space typo.

// This   comment   has    multiple   ones.

// This comment does not have any multi-space issues.

// This comment is missing punctuation

// This comment has a typo after punctuation!s

// This comment ends with parentheses, which are considered
// valid punctuation by this linter as long as comments end on a
// closing bracket/parentheses (like in this case)

// This comment is a multi-line comment with  space typos in
// it. It does not have any  punctuation issues though.

// This comment is a multi-line comment with missing punctuation
// so the first line should not trigger a warning, only the last
// one should

// This is a japanese comment. It should work just like any other language.
// 日本語!

// This is a japanese comment with a double space typo.
// 私   はカモメが好きです！!

// This is an arabic comment. For arabic, punctuation is inverted
// and the check should be done differently. This is not supported
// for now so this comment should trigger an error.
// عربى!
