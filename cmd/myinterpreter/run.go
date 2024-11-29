package main

func Run(lexer *Lexer) {
	// Parse the tokens to build the AST
	isRunning = true
	parser := NewParser(lexer.tokens)
	parser.parse()
	// Check for syntax errors
	if tokenError == SYNTAX_ERROR {
		return
	}
	// Evaluate the AST and print the result
	env := newEnvironment()
	runStmts(parser.stmts, env)
}

func copyMap(oldMap map[string]Token) map[string]Token {
	newMap := make(map[string]Token)
	for key, value := range oldMap {
		newMap[key] = value
	}
	return newMap
}

func runStmts(stmts []Stmt, env *Environment) {
	for _, stmt := range stmts {
		switch stmt.(type) {
		case PrintStmt:
			printResult(evaluateStmt(stmt, env))
		case ExprStmt, VarStmt:
			evaluateStmt(stmt, env)
		case BlockStmt:
			newEnv := newEnvironment()
			newEnv.values = copyMap(env.values)
			bs := stmt.(BlockStmt)
			runStmts(bs.stmts, newEnv)
		}
	}
}
