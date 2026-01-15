# Luma – A Minimal Programming Language

**Luma** is a minimal interpreted programming language built from scratch in **Go**. It is purely designed to provide a clean and simple syntax without extra complexity.

## 🚀 Installation

Ensure you have [Go](https://go.dev/) (version 1.18+) installed.

### 1. Clone the Repository

```bash
git clone https://github.com/Gazi2050/Luma.git
cd Luma
```

### 2. Global Installation

| Platform       | One-Line Install Command                                                         | 🗑️ Uninstall                       |
| :------------- | :------------------------------------------------------------------------------- | :--------------------------------- |
| **🐧 Linux**   | `go build -o luma main.go && sudo mv luma /usr/local/bin/luma`                   | `sudo rm /usr/local/bin/luma`      |
| **🍎 macOS**   | `go build -o luma main.go && sudo mv luma /usr/local/bin/luma`                   | `sudo rm /usr/local/bin/luma`      |
| **🪟 Windows** | `go build -o luma.exe main.go && move luma.exe C:\Windows\System32\` (CMD Admin) | `del C:\Windows\System32\luma.exe` |

The `luma` command is now ready!

## 🕹️ Usage

Luma provides two primary environments for development.

### 1. Interactive REPL

Execute code line by line for quick testing.

```bash
./luma
```

**Example:**

```text
>> let name = "Luma";
>> log("Hello " + name)
Hello Luma
```

_(Use `Ctrl+C` to exit)_

### 2. Script Execution

Run entire `.lu` files containing your Luma code.

```bash
./luma run index.lu
```

## ✨ Features

### Variables

Luma provides two ways to store data depending on whether you need to change it later.

- **How it works:**
  - `let` declarations create mutable variables. You can reassign them using the `=` operator (e.g., `x = 20;`).
  - `const` declarations are read-only. Once assigned, their value is locked.
- **Usage Tips:**
  - ✅ Use `const` for values that never change (like configuration or math constants).
  - ✅ Use `let` for counters or values that must be updated.
  - ❌ Don't try to reassign a `const`. The interpreter will throw an error.

```lu
let score = 100;    // Can be updated
const PI = 3.14;    // Fixed value
```

### Arrays

Represent ordered lists of values.

- **How it works:** Arrays in Luma can store a mix of any type (numbers, strings, booleans, even other arrays or objects). Accessing elements starts at index `0`.
- **Usage Tips:**
  - ✅ Use `len()` to get the size of an array before looping.
  - ✅ Use arrays for lists of similar items.
  - ❌ Don't access an index larger than the array size; it will return `null`.

```lu
let list = [1, "Luma", true];
log(list[0]); // Output: 1
```

### Objects

Store data in key-value pairs, similar to a dictionary or map.

- **How it works:** You can access property values using the **Dot Notation** (e.g., `user.name`).
- **Usage Tips:**
  - ✅ Use objects to group related data (like a user's profile).
  - ✅ Keys are strings by default.
  - ❌ Don't use reserved keywords as object keys.

```lu
let user = { name: "Gazi", id: 1 };
log(user.name); // Output: Gazi
```

### Condition

Control the flow of your program based on logic.

- **How it works:** The `if` statement evaluates a condition. If it is `true`, the first block runs. If `false`, the optional `else` block runs.
- **Usage Tips:**
  - ✅ Always wrap your condition in parentheses `()`.
  - ✅ Use comparison operators like `==`, `!=`, `>`, `<`, `>=`, `<=`.
  - ❌ Don't forget the `{}` curly braces even for single-line logic.

```lu
if (score > 50) {
    log("Winner!");
} else {
    log("Try Again");
}
```

### Loop

Luma features a powerful, single-keyword loop that handles all iteration needs.

- **How it works:** The `loop` takes three parts separated by semicolons: `initialization; condition; increment`.
- **Usage Tips:**
  - ✅ Be careful with infinite loops; ensure the condition eventually becomes `false`.
  - ✅ Use `let` for the iterator variable inside the loop initialization.
  - ❌ Don't omit the semicolons inside the `loop` parentheses.

```lu
loop (let i = 0; i < 5; i = i + 1) {
    log("Iteration:", i);
}
```

### Function

Functions are "First-Class" in Luma, meaning they can be assigned to variables and passed around.

- **How it works:** Use the `fn` keyword followed by the name, parameters, and a block of code. Use `return` to send a value back.
- **Usage Tips:**
  - ✅ Give your functions descriptive names like `calculateTotal`.
  - ✅ Use functions to avoid repeating code.
  - ❌ Functions without a `return` statement will return the value of the last expression executed.

```lu
fn greet(name) {
    return "Hello, " + name;
}
log(greet("User"));
```

### Math

Luma handles complex mathematical expressions with ease.

- **How it works:** It follows standard math precedence (Multiplication/Division before Addition/Subtraction). Parentheses `()` can be used to override this.
- **Usage Tips:**
  - ✅ Use `()` to make your math logic easier to read.
  - ✅ Both integers and decimals are supported.
  - ❌ Don't mix incompatible types (like `10 + "string"`) unless you intend to concatenate.

```lu
let result = (10 + 2) * 5 / 2;
log(result); // Output: 30
```

### Built-in Functions

Standard tools that come pre-installed in Luma.

- **How it works:**
  - `log(...)`: Prints one or more values to the console. It works for arrays and nested objects too!
  - `len(val)`: Returns the number of characters in a string or items in an array.
- **Usage Tips:**
  - ✅ Use `log()` with multiple arguments to debug values easily.
  - ✅ Use `len()` to check if a string is empty (`len(str) == 0`).

```lu
log("Items:", [1, 2, 3]);
let size = len("Luma"); // Returns 4
```

## 📖 Deep Dive

To understand how Luma translates text into logic, check out the implementation details:

- [**The Luma Engine (LANGUAGE_DETAILS.md)**](file:///home/gazi/Projects/my-projects/Luma/LANGUAGE_DETAILS.md)
