package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Standard ML lexer.
var StandardML = internal.Register(MustNewLexer(
	&Config{
		Name:      "Standard ML",
		Aliases:   []string{"sml"},
		Filenames: []string{"*.sml", "*.sig", "*.fun"},
		MimeTypes: []string{"text/x-standardml", "application/x-standardml"},
	},
	Rules{
		"whitespace": {
			{`\s+`, Text, nil},
			{`\(\*`, CommentMultiline, Push("comment")},
		},
		"delimiters": {
			{`\(|\[|\{`, Punctuation, Push("main")},
			{`\)|\]|\}`, Punctuation, Pop(1)},
			{`\b(let|if|local)\b(?!\')`, KeywordReserved, Push("main", "main")},
			{`\b(struct|sig|while)\b(?!\')`, KeywordReserved, Push("main")},
			{`\b(do|else|end|in|then)\b(?!\')`, KeywordReserved, Pop(1)},
		},
		"core": {
			{`(_|\}|\{|\)|;|,|\[|\(|\]|\.\.\.)`, Punctuation, nil},
			{`#"`, LiteralStringChar, Push("char")},
			{`"`, LiteralStringDouble, Push("string")},
			{`~?0x[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`0wx[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`0w\d+`, LiteralNumberInteger, nil},
			{`~?\d+\.\d+[eE]~?\d+`, LiteralNumberFloat, nil},
			{`~?\d+\.\d+`, LiteralNumberFloat, nil},
			{`~?\d+[eE]~?\d+`, LiteralNumberFloat, nil},
			{`~?\d+`, LiteralNumberInteger, nil},
			{`#\s*[1-9][0-9]*`, NameLabel, nil},
			{`#\s*([a-zA-Z][\w']*)`, NameLabel, nil},
			{"#\\s+([!%&$#+\\-/:<=>?@\\\\~`^|*]+)", NameLabel, nil},
			{`\b(datatype|abstype)\b(?!\')`, KeywordReserved, Push("dname")},
			{`(?=\b(exception)\b(?!\'))`, Text, Push("ename")},
			{`\b(functor|include|open|signature|structure)\b(?!\')`, KeywordReserved, Push("sname")},
			{`\b(type|eqtype)\b(?!\')`, KeywordReserved, Push("tname")},
			{`\'[\w\']*`, NameDecorator, nil},
			{`([a-zA-Z][\w']*)(\.)`, NameNamespace, Push("dotted")},
			{`\b(abstype|and|andalso|as|case|datatype|do|else|end|exception|fn|fun|handle|if|in|infix|infixr|let|local|nonfix|of|op|open|orelse|raise|rec|then|type|val|with|withtype|while|eqtype|functor|include|sharing|sig|signature|struct|structure|where)\b`, KeywordReserved, nil},
			{`([a-zA-Z][\w']*)`, Name, nil},
			{`\b(:|\|,=|=>|->|#|:>)\b`, KeywordReserved, nil},
			{"([!%&$#+\\-/:<=>?@\\\\~`^|*]+)", Name, nil},
		},
		"dotted": {
			{`([a-zA-Z][\w']*)(\.)`, NameNamespace, nil},
			// ignoring reserved words
			{`([a-zA-Z][\w']*)`, Name, Pop(1)},
			// ignoring reserved words
			{"([!%&$#+\\-/:<=>?@\\\\~`^|*]+)", Name, Pop(1)},
			{`\s+`, Error, nil},
			{`\S+`, Error, nil},
		},
		"root": {
			Default(Push("main")),
		},
		"main": {
			Include("whitespace"),
			{`\b(val|and)\b(?!\')`, KeywordReserved, Push("vname")},
			{`\b(fun)\b(?!\')`, KeywordReserved, Push("#pop", "main-fun", "fname")},
			Include("delimiters"),
			Include("core"),
			{`\S+`, Error, nil},
		},
		"main-fun": {
			Include("whitespace"),
			{`\s`, Text, nil},
			{`\(\*`, CommentMultiline, Push("comment")},
			{`\b(fun|and)\b(?!\')`, KeywordReserved, Push("fname")},
			{`\b(val)\b(?!\')`, KeywordReserved, Push("#pop", "main", "vname")},
			{`\|`, Punctuation, Push("fname")},
			{`\b(case|handle)\b(?!\')`, KeywordReserved, Push("#pop", "main")},
			Include("delimiters"),
			Include("core"),
			{`\S+`, Error, nil},
		},
		"char": {
			{`[^"\\]`, LiteralStringChar, nil},
			{`\\[\\"abtnvfr]`, LiteralStringEscape, nil},
			{`\\\^[\x40-\x5e]`, LiteralStringEscape, nil},
			{`\\[0-9]{3}`, LiteralStringEscape, nil},
			{`\\u[0-9a-fA-F]{4}`, LiteralStringEscape, nil},
			{`\\\s+\\`, LiteralStringInterpol, nil},
			{`"`, LiteralStringChar, Pop(1)},
		},
		"string": {
			{`[^"\\]`, LiteralStringDouble, nil},
			{`\\[\\"abtnvfr]`, LiteralStringEscape, nil},
			{`\\\^[\x40-\x5e]`, LiteralStringEscape, nil},
			{`\\[0-9]{3}`, LiteralStringEscape, nil},
			{`\\u[0-9a-fA-F]{4}`, LiteralStringEscape, nil},
			{`\\\s+\\`, LiteralStringInterpol, nil},
			{`"`, LiteralStringDouble, Pop(1)},
		},
		"breakout": {
			{`(?=\b(where|do|handle|if|sig|op|while|case|as|else|signature|andalso|struct|infixr|functor|in|structure|then|local|rec|end|fun|of|orelse|val|include|fn|with|exception|let|and|infix|sharing|datatype|type|abstype|withtype|eqtype|nonfix|raise|open)\b(?!\'))`, Text, Pop(1)},
		},
		"sname": {
			Include("whitespace"),
			Include("breakout"),
			{`([a-zA-Z][\w']*)`, NameNamespace, nil},
			Default(Pop(1)),
		},
		"fname": {
			Include("whitespace"),
			{`\'[\w\']*`, NameDecorator, nil},
			{`\(`, Punctuation, Push("tyvarseq")},
			{`([a-zA-Z][\w']*)`, NameFunction, Pop(1)},
			{"([!%&$#+\\-/:<=>?@\\\\~`^|*]+)", NameFunction, Pop(1)},
			Default(Pop(1)),
		},
		"vname": {
			Include("whitespace"),
			{`\'[\w\']*`, NameDecorator, nil},
			{`\(`, Punctuation, Push("tyvarseq")},
			{"([a-zA-Z][\\w']*)(\\s*)(=(?![!%&$#+\\-/:<=>?@\\\\~`^|*]+))", ByGroups(NameVariable, Text, Punctuation), Pop(1)},
			{"([!%&$#+\\-/:<=>?@\\\\~`^|*]+)(\\s*)(=(?![!%&$#+\\-/:<=>?@\\\\~`^|*]+))", ByGroups(NameVariable, Text, Punctuation), Pop(1)},
			{`([a-zA-Z][\w']*)`, NameVariable, Pop(1)},
			{"([!%&$#+\\-/:<=>?@\\\\~`^|*]+)", NameVariable, Pop(1)},
			Default(Pop(1)),
		},
		"tname": {
			Include("whitespace"),
			Include("breakout"),
			{`\'[\w\']*`, NameDecorator, nil},
			{`\(`, Punctuation, Push("tyvarseq")},
			{"=(?![!%&$#+\\-/:<=>?@\\\\~`^|*]+)", Punctuation, Push("#pop", "typbind")},
			{`([a-zA-Z][\w']*)`, KeywordType, nil},
			{"([!%&$#+\\-/:<=>?@\\\\~`^|*]+)", KeywordType, nil},
			{`\S+`, Error, Pop(1)},
		},
		"typbind": {
			Include("whitespace"),
			{`\b(and)\b(?!\')`, KeywordReserved, Push("#pop", "tname")},
			Include("breakout"),
			Include("core"),
			{`\S+`, Error, Pop(1)},
		},
		"dname": {
			Include("whitespace"),
			Include("breakout"),
			{`\'[\w\']*`, NameDecorator, nil},
			{`\(`, Punctuation, Push("tyvarseq")},
			{`(=)(\s*)(datatype)`, ByGroups(Punctuation, Text, KeywordReserved), Pop(1)},
			{"=(?![!%&$#+\\-/:<=>?@\\\\~`^|*]+)", Punctuation, Push("#pop", "datbind", "datcon")},
			{`([a-zA-Z][\w']*)`, KeywordType, nil},
			{"([!%&$#+\\-/:<=>?@\\\\~`^|*]+)", KeywordType, nil},
			{`\S+`, Error, Pop(1)},
		},
		"datbind": {
			Include("whitespace"),
			{`\b(and)\b(?!\')`, KeywordReserved, Push("#pop", "dname")},
			{`\b(withtype)\b(?!\')`, KeywordReserved, Push("#pop", "tname")},
			{`\b(of)\b(?!\')`, KeywordReserved, nil},
			{`(\|)(\s*)([a-zA-Z][\w']*)`, ByGroups(Punctuation, Text, NameClass), nil},
			{"(\\|)(\\s+)([!%&$#+\\-/:<=>?@\\\\~`^|*]+)", ByGroups(Punctuation, Text, NameClass), nil},
			Include("breakout"),
			Include("core"),
			{`\S+`, Error, nil},
		},
		"ename": {
			Include("whitespace"),
			{`(exception|and)\b(\s+)([a-zA-Z][\w']*)`, ByGroups(KeywordReserved, Text, NameClass), nil},
			{"(exception|and)\\b(\\s*)([!%&$#+\\-/:<=>?@\\\\~`^|*]+)", ByGroups(KeywordReserved, Text, NameClass), nil},
			{`\b(of)\b(?!\')`, KeywordReserved, nil},
			Include("breakout"),
			Include("core"),
			{`\S+`, Error, nil},
		},
		"datcon": {
			Include("whitespace"),
			{`([a-zA-Z][\w']*)`, NameClass, Pop(1)},
			{"([!%&$#+\\-/:<=>?@\\\\~`^|*]+)", NameClass, Pop(1)},
			{`\S+`, Error, Pop(1)},
		},
		"tyvarseq": {
			{`\s`, Text, nil},
			{`\(\*`, CommentMultiline, Push("comment")},
			{`\'[\w\']*`, NameDecorator, nil},
			{`[a-zA-Z][\w']*`, Name, nil},
			{`,`, Punctuation, nil},
			{`\)`, Punctuation, Pop(1)},
			{"[!%&$#+\\-/:<=>?@\\\\~`^|*]+", Name, nil},
		},
		"comment": {
			{`[^(*)]`, CommentMultiline, nil},
			{`\(\*`, CommentMultiline, Push()},
			{`\*\)`, CommentMultiline, Pop(1)},
			{`[(*)]`, CommentMultiline, nil},
		},
	},
))
