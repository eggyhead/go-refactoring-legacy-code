# Rewriting the "Legacy Code" Manual to Go

Michael Feathers’ "Working Effectively with Legacy Code" is a classic guide of software maintenance. It teaches us that Legacy Code is problematic mostly because it lacks tests, and is thus dangerous to update since we can't assert the current state to compare changes. 

The book's framework relies on finding Seams: points in a codebase where you can alter behavior without editing source code directly. However, for a Go developer, the book can feel like it was written in a foreign tongue. Its primary weapons of Class Inheritance and Subclass/Override simply don't exist in our world.

To understand the problematic legacy patterns in Go services, we rewrote the examples and we are sharing some of the details. 


## 1. The Death of "Subclass and Override"

In Java or C++, if a method hits a real database, you subclass the object and override that method in your test. In Go, this is a dead end.

We move from **Inheritance** to **Function Fields**.

Instead of a hard-coded side effect, we turn the problematic call into a field on the struct. This creates a "Remote Control" for our tests.

### Before: The "Toxic" Method

```go
func (rm *ReportManager) Generate() error {
    // Hard-coded global call! Impossible to mock.
    return smtp.SendMail("smtp.example.com", auth, from, to, msg)
}

```

### After: The Function Field Seam

```go
type ReportManager struct {
    // We wrap the global call in a function field
    SendEmail func(to string, body string) error
}

func (rm *ReportManager) Generate() error {
    return rm.SendEmail("admin@test.com", "Report Data")
}

```

In our tests, we simply swap `SendEmail` for a function that does nothing, allowing us to test the surrounding logic without actually sending mail.



## 2. Implicit Interfaces: The Silent Decoupler

The book frequently suggests **Extract Interface**. In Java, this is an operation requiring `implements` keywords and a rigid hierarchy. In Go, interfaces are **implicit**.

Practically speaking, this means the code *calling* the dependency gets to decide the interface. This is a powerful decoupling tool unique to Go's philosophy.

### The "Provider" Pattern

By using the **NewXXX** pattern with an interface, we create a composable system where the production code and the test code are identical to the compiler, but worlds apart in behavior:

```go
type RevenueReader interface {
    GetMonthlyRevenue(month time.Month) (float64, error)
}

// Our Service doesn't care if it's SQL, a CSV, or a Mock
func NewReportManager(r RevenueReader) *ReportManager {
    return &ReportManager{reader: r}
}

```

## 3. Primitivizing the Hard-to-Test

One of the most powerful translations we discovered wasn't in the book's mechanics, but in Go's ergonomics: **Primitivize Parameter.**

Legacy structs often suck the "entire world" into a test. If a struct has 50 fields and a database handle, testing one small method becomes a nightmare. We found that extracting the math into a **Pure Function** that speaks only in Go primitives (`int`, `string`, `bool`) is often better than any complex mocking.


## The Verdict: Learning by Translation

Translating these 20-year-old patterns into Go did more than just improve our test coverage; it deepened our understanding of Go’s unique ergonomics.

We’ve compiled these "translations" into a Learning Repository. Inside, you’ll find "Before and After" examples for:

* **Hard-coded SQL:** Moving from `sql.DB` to interfaces.
* **Time Clings:** Replacing `time.Now()` with function fields.
* **Global Singletons:** Replacing package-level variables with getters.

By spotting these patterns in your production code, you can use this repository as a field guide for safely breaking dependencies and finally writing that first test.