# How Luma Works: Architecture & Internal Logic

This document provides a comprehensive breakdown of the engine behind the Luma programming language. Built entirely in Go, Luma follows a modular architecture where each component has a distinct responsibility in the lifecycle of a script—from raw text to a running program.



## 📂 Project Structure & Package Breakdown

Luma's source code is organized into modular Go packages, each responsible for a specific stage of the interpreter lifecycle.

### 🏠 Root Directory
**File:** [main.go](file:///home/gazi/Projects/my-projects/Luma/main.go)

The starting point of the application. It handles CLI arguments and initializes the execution environment.
- **CLI Handling**: Decides between starting the interactive REPL or running a `.lu` file.
- **Orchestration**: Connects the Lexer, Parser, and Evaluator to process code from start to finish.

### 📦 `token/`
**File:** [token/token.go](file:///home/gazi/Projects/my-projects/Luma/token/token.go)

Defines the fundamental "words" of the Luma language.
- **TokenType**: A custom string type used to categorize symbols (e.g., `LET`, `FN`, `PLUS`).
- **Token Struct**: Stores the type and the actual text (literal) from the source code.
- **Lookup Table**: Maps keywords like `if` and `loop` to their respective token types.

### 📦 `lexer/`
**File:** [lexer/lexer.go](file:///home/gazi/Projects/my-projects/Luma/lexer/lexer.go)

The Lexer (Scanner) converts the raw source string into a stream of tokens.
- **Character Reading**: Iterates through the input character by character using a pointer system.
- **Pattern Matching**: Recognizes identifiers, numbers, and strings while ignoring whitespace and comments.
- **Peek Logic**: Essential for multi-character operators like `==` or `<=`.

### 📦 `ast/`
**File:** [ast/ast.go](file:///home/gazi/Projects/my-projects/Luma/ast/ast.go)

Contains the definition of the Abstract Syntax Tree (AST), the internal logical representation of the code.
- **Node Interfaces**: Every construct must implement the `Node` interface.
- **Hierarchy**: Distinguishes between **Statements** (actions like variable declaration) and **Expressions** (computations that result in values).
- **String Printing**: Every node can reconstruct its source-like representation for debugging.

### 📦 `parser/`
**File:** [parser/parser.go](file:///home/gazi/Projects/my-projects/Luma/parser/parser.go)

Transforms the flat token stream into a nested AST using a **Pratt Parser**.
- **Precedence Management**: Uses a weight-based system to handle the "order of operations" (e.g., multiplication before addition).
- **Recursive Parsing**: Decodes complex structures like nested function calls and object literals.
- **Syntax Validation**: Ensures the code follows Luma's grammatical rules and reports errors.

### 📦 `evaluator/`
**Files:** [evaluator.go](file:///home/gazi/Projects/my-projects/Luma/evaluator/evaluator.go) | [environment.go](file:///home/gazi/Projects/my-projects/Luma/evaluator/environment.go)

The "Engine" that carries out the logic of the program.
- **`evaluator.go`**: Traverses the AST and executes Go logic for each node. It handles arithmetic, function calls, and data structure manipulation.
- **`environment.go`**: Manages the life of variables. It supports **Lexical Scoping** (nested scopes) and strictly enforces `const` reassignment rules.
- **Pretty Printing**: Includes logic to format Luma objects and arrays into human-readable strings for the `log()` function.
