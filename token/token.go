package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Literals
	IDENT  = "IDENT"  // variable names, function names
	NUMBER = "NUMBER" // numeric literal (digits)
	STRING = "STRING" // quoted string literal

	// Keywords - Variables
	SET = "SET"
	TO  = "TO"

	// Keywords - Arithmetic expressions
	PLUS    = "PLUS"
	TIMES   = "TIMES"
	DIVIDED = "DIVIDED"
	MINUS   = "MINUS"

	// Keywords - In-place arithmetic
	INCREASE = "INCREASE"
	DECREASE = "DECREASE"
	BY       = "BY"

	// Keywords - Conditionals
	IF        = "IF"
	THEN      = "THEN"
	OTHERWISE = "OTHERWISE"
	EQUALS    = "EQUALS"
	IS        = "IS"
	GREATER   = "GREATER"
	LESS      = "LESS"
	THAN      = "THAN"

	// Keywords - Logical
	AND = "AND"
	OR  = "OR"
	NOT = "NOT"

	// Keywords - Loops
	WHILE = "WHILE"
	DO    = "DO"
	FOR   = "FOR"
	EACH  = "EACH"
	IN    = "IN"

	// Keywords - Blocks
	DONE = "DONE"

	// Keywords - Functions
	RETURN = "RETURN"
	CALL   = "CALL"
	WITH   = "WITH"

	// Keywords - I/O
	SAY = "SAY"
	ASK = "ASK"

	// Keywords - Lists
	A      = "A"
	LIST   = "LIST"
	OF     = "OF"
	LENGTH = "LENGTH"
	APPEND = "APPEND"
	GET    = "GET"
	ITEM   = "ITEM"
	FROM   = "FROM"

	// Keywords - Comparison helpers
	INTO = "INTO"

	// Keywords - HTTP
	FETCH   = "FETCH"
	SEND    = "SEND"
	PUT     = "PUT"
	DELETE  = "DELETE"
	BODY    = "BODY"
	STATUS  = "STATUS"
	HEADER  = "HEADER"
	HEADERS = "HEADERS"

	// Keywords - JSON
	PARSE  = "PARSE"
	JSON   = "JSON"
	FIELD  = "FIELD"
	ENCODE = "ENCODE"
	AS     = "AS"

	// Keywords - Web Server
	SERVE      = "SERVE"
	ON         = "ON"
	WHEN       = "WHEN"
	REQUEST    = "REQUEST"
	AT         = "AT"
	USING      = "USING"
	REPLY      = "REPLY"
	ROUTE      = "ROUTE"
	BACKGROUND = "BACKGROUND"
	STOP       = "STOP"
	SERVER     = "SERVER"
	QUERY      = "QUERY"
	METHOD     = "METHOD"
	PATH       = "PATH"

	// Number words (0-19)
	ZERO      = "ZERO"
	ONE       = "ONE"
	TWO       = "TWO"
	THREE     = "THREE"
	FOUR      = "FOUR"
	FIVE      = "FIVE"
	SIX       = "SIX"
	SEVEN     = "SEVEN"
	EIGHT     = "EIGHT"
	NINE      = "NINE"
	TEN       = "TEN"
	ELEVEN    = "ELEVEN"
	TWELVE    = "TWELVE"
	THIRTEEN  = "THIRTEEN"
	FOURTEEN  = "FOURTEEN"
	FIFTEEN   = "FIFTEEN"
	SIXTEEN   = "SIXTEEN"
	SEVENTEEN = "SEVENTEEN"
	EIGHTEEN  = "EIGHTEEN"
	NINETEEN  = "NINETEEN"

	// Number words (tens)
	TWENTY  = "TWENTY"
	THIRTY  = "THIRTY"
	FORTY   = "FORTY"
	FIFTY   = "FIFTY"
	SIXTY   = "SIXTY"
	SEVENTY = "SEVENTY"
	EIGHTY  = "EIGHTY"
	NINETY  = "NINETY"

	// Number words (large)
	HUNDRED  = "HUNDRED"
	THOUSAND = "THOUSAND"
	MILLION  = "MILLION"
)

var keywords = map[string]TokenType{
	// Core keywords
	"set":       SET,
	"to":        TO,
	"plus":      PLUS,
	"times":     TIMES,
	"divided":   DIVIDED,
	"minus":     MINUS,
	"increase":  INCREASE,
	"decrease":  DECREASE,
	"by":        BY,
	"if":        IF,
	"then":      THEN,
	"otherwise": OTHERWISE,
	"equals":    EQUALS,
	"is":        IS,
	"greater":   GREATER,
	"less":      LESS,
	"than":      THAN,
	"and":       AND,
	"or":        OR,
	"not":       NOT,
	"while":     WHILE,
	"do":        DO,
	"for":       FOR,
	"each":      EACH,
	"in":        IN,
	"done":      DONE,
	"return":    RETURN,
	"call":      CALL,
	"with":      WITH,
	"say":       SAY,
	"ask":       ASK,
	"a":         A,
	"list":      LIST,
	"of":        OF,
	"length":    LENGTH,
	"append":    APPEND,
	"get":       GET,
	"item":      ITEM,
	"from":      FROM,
	"into":      INTO,

	// HTTP keywords
	"fetch":   FETCH,
	"send":    SEND,
	"put":     PUT,
	"delete":  DELETE,
	"body":    BODY,
	"status":  STATUS,
	"header":  HEADER,
	"headers": HEADERS,

	// JSON keywords
	"parse":  PARSE,
	"json":   JSON,
	"field":  FIELD,
	"encode": ENCODE,
	"as":     AS,

	// Web server keywords
	"serve":      SERVE,
	"on":         ON,
	"when":       WHEN,
	"request":    REQUEST,
	"at":         AT,
	"using":      USING,
	"reply":      REPLY,
	"route":      ROUTE,
	"background": BACKGROUND,
	"stop":       STOP,
	"server":     SERVER,
	"query":      QUERY,
	"method":     METHOD,
	"path":       PATH,

	// Number words
	"zero":      ZERO,
	"one":       ONE,
	"two":       TWO,
	"three":     THREE,
	"four":      FOUR,
	"five":      FIVE,
	"six":       SIX,
	"seven":     SEVEN,
	"eight":     EIGHT,
	"nine":      NINE,
	"ten":       TEN,
	"eleven":    ELEVEN,
	"twelve":    TWELVE,
	"thirteen":  THIRTEEN,
	"fourteen":  FOURTEEN,
	"fifteen":   FIFTEEN,
	"sixteen":   SIXTEEN,
	"seventeen": SEVENTEEN,
	"eighteen":  EIGHTEEN,
	"nineteen":  NINETEEN,
	"twenty":    TWENTY,
	"thirty":    THIRTY,
	"forty":     FORTY,
	"fifty":     FIFTY,
	"sixty":     SIXTY,
	"seventy":   SEVENTY,
	"eighty":    EIGHTY,
	"ninety":    NINETY,
	"hundred":   HUNDRED,
	"thousand":  THOUSAND,
	"million":   MILLION,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func IsNumberWord(t TokenType) bool {
	switch t {
	case ZERO, ONE, TWO, THREE, FOUR, FIVE, SIX, SEVEN, EIGHT, NINE,
		TEN, ELEVEN, TWELVE, THIRTEEN, FOURTEEN, FIFTEEN, SIXTEEN,
		SEVENTEEN, EIGHTEEN, NINETEEN, TWENTY, THIRTY, FORTY, FIFTY,
		SIXTY, SEVENTY, EIGHTY, NINETY, HUNDRED, THOUSAND, MILLION:
		return true
	}
	return false
}

func NumberWordValue(t TokenType) int64 {
	switch t {
	case ZERO:
		return 0
	case ONE:
		return 1
	case TWO:
		return 2
	case THREE:
		return 3
	case FOUR:
		return 4
	case FIVE:
		return 5
	case SIX:
		return 6
	case SEVEN:
		return 7
	case EIGHT:
		return 8
	case NINE:
		return 9
	case TEN:
		return 10
	case ELEVEN:
		return 11
	case TWELVE:
		return 12
	case THIRTEEN:
		return 13
	case FOURTEEN:
		return 14
	case FIFTEEN:
		return 15
	case SIXTEEN:
		return 16
	case SEVENTEEN:
		return 17
	case EIGHTEEN:
		return 18
	case NINETEEN:
		return 19
	case TWENTY:
		return 20
	case THIRTY:
		return 30
	case FORTY:
		return 40
	case FIFTY:
		return 50
	case SIXTY:
		return 60
	case SEVENTY:
		return 70
	case EIGHTY:
		return 80
	case NINETY:
		return 90
	case HUNDRED:
		return 100
	case THOUSAND:
		return 1000
	case MILLION:
		return 1000000
	}
	return 0
}

func IsMultiplier(t TokenType) bool {
	return t == HUNDRED || t == THOUSAND || t == MILLION
}

func IsArithmeticOperator(t TokenType) bool {
	return t == PLUS || t == MINUS || t == TIMES || t == DIVIDED
}
