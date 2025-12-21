# Luma – A Minimal Programming Language

**Luma** is a lightweight, interpreted programming language built in **Go**, designed to learn **data structures, algorithms, and core language concepts** through hands-on implementation.

---

## Features (v1)

- Variables: `let` and `const`
- Types: `number`, `string`, `boolean`, `array`, `object`
- Built-in functions: `log()`, `len()`
- Math operations: `+`, `-`, `*`, `/`
- Boolean comparisons: `==`, `!=`, `<`, `>`, `<=`, `>=`
- Loops using `loop` keyword
- User-defined functions (`fn`) with parameters and return
- Run `.lu` scripts
- Interactive REPL

---

## Installation

Clone the repository and build Luma:

```bash
git clone <your-repo-url>
cd luma
go build -o luma luma.go
```

This produces the `luma` executable.

---

## Usage

### 1. REPL (Interactive Mode)

Start Luma REPL:

```bash
./luma
```

Example:

```text
Luma REPL v1 - Type your code
>> let x = 10;
>> let name = "Luma";
>> log(x);
10
>> log(name);
Luma
>> let isReady = true;
>> log(isReady);
true
```

### 2. Running `.lu` Scripts

Create a Luma script file: `index.lu`

```lu
let x = 10;
let y = 20;
let arr = [x, y, 30];
let user = { name: "Luma", age: 1 };

fn add(a, b) {
    return a + b;
}

loop (let i = 0; i < len(arr); i = i + 1) {
    log(arr[i]);
}

log(add(x, y));
log(user.name);
```

Run the script using:

```bash
./luma run index.lu
```

Expected output:

```
10
20
30
30
Luma
```

### 3. Supported Syntax

#### Variables

```lu
let age = 25;       // number
const title = "Luma"; // string
let ready = true;   // boolean
```

#### Arrays

```lu
let nums = [1, 2, 3];
log(nums[0]); // 1
```

#### Objects

```lu
let user = { name: "Luma", age: 1 };
log(user.name); // Luma
```

#### Loops

```lu
loop (let i = 0; i < 5; i = i + 1) {
    log(i);
}
```

#### Functions

```lu
fn add(a, b) {
    return a + b;
}
log(add(2, 3)); // 5
```

#### Math & Boolean Operations

```lu
let result = 10 + 2 * 3;
log(result); // 16
let isBig = result > 10;
log(isBig); // true
```

---

### Notes

- Scripts must use semicolons `;` to terminate statements.
- Currently supports core features; more advanced features may be added in future versions.
