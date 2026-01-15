# 🛠️ Go Refactoring: Legacy to Clean Code Guide

This document explains the legacy design of the `report` and `legacy_user_processor` packages, highlights the key blockers for testing and maintainability, and provides a roadmap for refactoring to a clean, testable design.

## 🚀 How to Use

- **To run all tests:**  
  ```bash
  go test ./...
  ```

---

## 📊 Report Package: Legacy Analysis & Refactoring Guide

### 📄 Legacy Code: `report/report.go`

#### 🚩 Key Legacy Blockers

1. **Hard-coded time check:** Uses `time.Now()` directly.
2. **Direct SQL dependency:** Uses a real `*sql.DB` and SQL queries.
3. **Direct AWS S3 usage:** Instantiates and uses a real S3 client.
4. **Direct SMTP usage:** Sends real emails via SMTP.
5. **No dependency injection:** All dependencies are created inside the struct or method.

#### 🧩 Why Is This Hard to Test?

- You cannot control the current time in tests.
- You need a real database, AWS credentials, and SMTP server to run the code.
- Tests may send real emails or upload to S3.
- Logic is mixed with side effects, making unit testing nearly impossible.

#### 🛠️ Refactoring Roadmap

| Legacy Symptom         | Refactoring Technique         | Goal                                 |
|------------------------|------------------------------|--------------------------------------|
| Hard-coded time check  | Function Field               | Make time injectable for tests       |
| Hard-coded SQL         | Interface Substitution       | Allow fake DB for tests              |
| Hard-coded S3          | Interface Substitution       | Allow fake uploader for tests        |
| Hard-coded SMTP        | Function Field               | Allow fake mailer for tests          |
| Logic in method        | Primitivize Parameter        | Extract pure logic for easy testing  |

#### 🧑‍💻 See the Refactored Version

- [`report_fixed/report_fixed.go`](report_fixed/report_fixed.go): Fully refactored, testable version using dependency injection and pure functions.
- [`report_fixed/report_fixed_test.go`](report_fixed/report_fixed_test.go): Example tests using fakes and dependency injection.

- **To write your own tests:**  
  Use the patterns in `report_fixed_test.go`—inject fakes/mocks for all dependencies.

---

## 👤 Legacy User Processor: Analysis & Refactoring Guide

### 📄 Legacy Code: `legacy_user_processor/legacy_user_processor.go`

#### 🚩 Key Legacy Dependencies

1. **Global State:**  
   - `ProcessCount` is a global variable, making tests order-dependent and hard to parallelize.
2. **File System Dependency:**  
   - Uses `os.Open` directly, so tests require real files on disk.
3. **Tightly Coupled JSON Logic:**  
   - JSON decoding is embedded, making it hard to test logic in isolation.
4. **External Network Call:**  
   - Calls a real HTTP API (`http.Get`), making tests slow, flaky, and dependent on external services.
5. **Process Exit:**  
   - Calls `os.Exit`, which will kill the test runner if triggered.

#### 🧩 Why Is This Hard to Test?

- You can't inject fakes or mocks for the file system, network, or global state.
- Tests may have side effects (e.g., incrementing a global, making real HTTP calls, or exiting the process).
- Logic is not isolated—it's mixed with I/O and side effects.

#### 🛠️ Refactoring Roadmap

| Legacy Symptom         | Refactoring Technique         | Goal                                 |
|------------------------|------------------------------|--------------------------------------|
| Global variable        | Function Field or Interface  | Allow test to control state          |
| File system access     | Interface Substitution       | Allow fake file access in tests      |
| Network call           | Interface Substitution       | Allow fake HTTP client in tests      |
| os.Exit                | Function Field               | Prevent test runner from exiting     |
| Logic in method        | Primitivize Parameter        | Extract pure logic for easy testing  |

#### 🧑‍💻 See the Refactored Version

- [`processor_fixed/processor_fixed.go`](processor_fixed/processor_fixed.go): Refactored, testable version using interfaces and function fields.
- [`processor_fixed/processor_fixed_test.go`](processor_fixed/processor_fixed_test.go): Example tests using mocks and dependency injection.

---

## 📝 Next Steps

- Study the differences between the legacy and refactored files for both packages.
- Try writing a test for the legacy versions to see the pain points.
- Use the refactored versions as templates for breaking dependencies in production code