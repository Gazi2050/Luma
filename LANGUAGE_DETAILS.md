# How Luma Works: Architecture & Internal Logic

This document provides a comprehensive breakdown of the engine behind the Luma programming language. Built entirely in Go, Luma follows a modular architecture where each component has a distinct responsibility in the lifecycle of a script—from raw text to a running program.

---

## 🏗️ Core Architecture Overview

Luma operates as a **Tree-Walking Interpreter**. The execution flow follows these four stages:

1.  **Lexing**: Raw text is converted into tokens.
2.  **Parsing**: Tokens are organized into an Abstract Syntax Tree (AST).
3.  **Evaluating**: The AST is traversed to perform logic and compute values.
4.  **Environment Management**: Variables and scopes are tracked during evaluation.

---

## 📦 Component Breakdown

### 1. The Entry Point
**File:** [main.go](file:///home/gazi/Projects/my-projects/Luma/main.go)

The entry point orchestrates the entire process. It handles command-line arguments (deciding whether to start the REPL or run a file) and initializes the standard environment.

- **REPL Mode**: Repeatedly takes user input and passes it through the lexer → parser → evaluator pipeline.
- **Run Mode**: Reads an entire file from disk and evaluates it in a single pass.

---

### 2. Lexical Analysis (The Lexer)
**Package:** `lexer` | **File:** [lexer/lexer.go](file:///home/gazi/Projects/my-projects/Luma/lexer/lexer.go)

The Lexer is a state machine that scans source code character by character.

- **Tokens**: Every meaningful unit of code (like `let`, `+`, or `123`) is mapped to a constant defined in the [token/token.go](file:///home/gazi/Projects/my-projects/Luma/token/token.go) file.
- **Scanning Logic**: It handles complex cases like multi-character operators (`==`, `!=`) by peeking at the next character before deciding which token to emit.
- **Keyword Recognition**: It uses a lookup table to distinguish between user-defined variables (`x`) and language keywords (`if`, `loop`).

---

### 3. The Abstract Syntax Tree (AST)
**Package:** `ast` | **File:** [ast/ast.go](file:///home/gazi/Projects/my-projects/Luma/ast/ast.go)

The AST is the internal data structure that represents your code.

- **Nodes**: Every construct in Luma (a number literal, a function call, a mathematical expression) implements the `ast.Node` interface.
- **Statements vs Expressions**: Luma distinguishes between code that produces a value (Expressions) and code that performs an action (Statements). This distinction is critical for the parser and evaluator.

---

### 4. Parsing (The Pratt Parser)
**Package:** `parser` | **File:** [parser/parser.go](file:///home/gazi/Projects/my-projects/Luma/parser/parser.go)

Luma uses a **Pratt Parser** (Top-Down Operator Precedence). Unlike a simple recursive descent parser, a Pratt parser handles operator precedence (e.g., ensuring `*` happens before `+`) with high efficiency and readability.

- **Precedences**: Defined in a map, assigning values to tokens (e.g., `PRODUCT` is higher than `SUM`).
- **Infix & Prefix**: The parser registers functions for "prefix" operators (like `-` or `!`) and "infix" operators (like `+` or `*`).
- **Error Handling**: Collects informative error messages if the syntax is invalid.

---

### 5. Scoping & Environments
**Package:** `evaluator` | **File:** [evaluator/environment.go](file:///home/gazi/Projects/my-projects/Luma/evaluator/environment.go)

Luma features **Lexical Scoping**, managed by the `Env` struct.

- **Nested Scopes**: Every time a function or block is executed, a new "enclosed" environment is created. This environment keeps a pointer to its "outer" parent.
- **Lookup Process**: When you use a variable, Luma checks the local scope first, then moves up the chain of parent environments until it finds the variable or hits the global scope.
- **Constant Protection**: The environment tracks whether a variable was defined with `const`, preventing re-assignment at runtime.

---

### 6. The Evaluator (Tree-Walking)
**Package:** `evaluator` | **File:** [evaluator/evaluator.go](file:///home/gazi/Projects/my-projects/Luma/evaluator/evaluator.go)

The Evaluator is where the code actually "runs". It uses recursion to visit every node in the AST and perform the corresponding Go operations.

- **Go Mapping**: Luma values are mapped directly to Go types (e.g., Luma `number` → Go `int`).
- **First-Class Functions**: Functions are stored as AST nodes within the environment, allowing them to be passed as arguments or returned from other functions.
- **Formatters**: Includes a custom `FormatLumaValue` function to pretty-print objects and arrays in a human-readable way.

---

## 🛠️ Implementation Details

### Why Go?
We chose Go for Luma because of its:
- **Fast Execution**: Compiles to native code, making the recursive evaluator very responsive.
- **Robust Standard Library**: Handling strings, maps, and file I/O is seamless.
- **Static Typing**: Ensures the interpreter's internal logic is type-safe.

### Handling `loop`
The `loop` statement in Luma is uniquely implemented to behave like a standard C `for` loop but within a single keyword. During evaluation, it creates a dedicated scope to handle the initialization variable, ensuring it doesn't "leak" into the outer scope.
