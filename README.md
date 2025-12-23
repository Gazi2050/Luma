# Luma – A Minimal Programming Language

**Luma** is a lightweight, interpreted programming language built in **Go**, designed for a clean and simple developer experience while supporting core language features.

---

## Features

- **Variables**: `let` and `const` (with constant protection).
- **Types**: `number`, `string`, `boolean`, `array`, `object`.
- **Built-in functions**: `log()`, `len()`.
- **Advanced Logic**: Full support for nested functions, lexical scoping, and recursion.
- **Control Flow**: `if/else` and `loop (init; condition; post)` syntax.
- **Improved Logging**: Pretty-printed output for arrays and objects.

---

## Installation

Clone the repository and build Luma:

```bash
git clone <your-repo-url>
cd Luma
go build -o luma main.go
```

This produces the `luma` executable.

---

## Usage

### 1. REPL (Interactive Mode)

Start Luma REPL:

```bash
./luma
```

### 2. Running `.lu` Scripts

Run the provided example script:

```bash
./luma run index.lu
```

### Sample `index.lu`

```lu
// 1. Variables & Types
let name = "Luma";
const version = 1;

// 2. Objects & Arrays
let user = { name: name, age: version };
let nums = [10, 20, 30];

log("Hello " + user.name);
log(nums[0]);
log(len(nums));

// 3. Math & Comparisons
let result = (10 + 2) * 3;
log(result > 30); // true

// 4. Functions
fn greet(n) {
    if (len(n) > 0) {
        return "Hi, " + n;
    }
    return "Hi, guest";
}
log(greet(user.name));

// 5. Loops
loop (let i = 0; i < 3; i = i + 1) {
    log(i);
}
```

**Expected Output:**

```text
Hello Luma
10
3
true
Hi, Luma
0
1
2
```

---

### Syntax Highlights

#### Variables

```lu
let age = 25;
const title = "Luma"; // Cannot be reassigned
```

#### Arrays & Objects

```lu
let list = [1, "two", true];
let obj = { key: "value", nested: { a: 1 } };
log(obj.nested.a);
```

#### Functions

```lu
fn add(a, b) {
    return a + b;
}
```

#### Loops

```lu
loop (let i = 0; i < 5; i = i + 1) {
    log(i);
}
```

---

### Notes

- Semicolons `;` are optional after blocks and at the end of many statements.
- The interpreter uses a Pratt Parser for reliable expression evaluation.
- All code is executed via a tree-walking evaluator with scoped environments.
- hello
