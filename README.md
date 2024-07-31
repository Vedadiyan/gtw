# GTW (Gateway) Framework Documentation

## Table of Contents

1. [Introduction](#1-introduction)
2. [Existing Problem](#2-existing-problem)
3. [Solution](#3-solution)
4. [GTW vs. Fiber/Gin](#4-gtw-vs-fibergin)
5. [Getting Started](#5-getting-started)
6. [Features](#6-features)

## 1. Introduction

GTW (Gateway) is a new web development framework for Go, designed to make web application codebases more maintainable and elegant. It introduces an annotation-driven approach that decouples logic from configuration, allowing developers to create more contextual and easily maintainable code.

## 2. Existing Problem

Many popular Go web frameworks, such as Fiber and Gin, use a handler-based approach to create web APIs. While powerful, this approach can lead to several issues in large-scale projects:

- Lack of clear structure in route definitions and handler logic
- Limited contextual information and difficulty in sharing dependencies
- Challenges in implementing clean dependency injection
- Reduced code reusability due to anonymous functions
- Increased difficulty in unit testing isolated handlers

Example of typical route definition in Fiber:

```go
app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
})
```

## 3. Solution

GTW addresses these issues by introducing an annotation-driven, struct-based approach that separates configuration from handler logic. This solution offers:

- Clear structure for defining APIs and their associated handlers
- Improved contextuality and organization of related routes
- Built-in support for dependency injection
- Enhanced code reusability through struct methods
- Better testability due to decoupled components

Example of API definition in GTW:

```go
type Foo struct {
    Metadata `prefix:"api"`
    Bar Service[int] `name:"test"`
    Get Handler `route:"/foo" method:"GET"`
}

func (f *Foo) GetHandler(httpCtx *HttpCtx) (Status, Response) {
    output := map[string]any{
        "Hello": t.Bar.Value(),
    }
    return 200, JSON(output)
}
```

## 4. GTW vs. Fiber/Gin

| Feature | GTW | Fiber/Gin |
|---------|-----|-----------|
| API Definition | Struct-based with annotations | Handler based |
| Route Configuration | Decoupled from logic | Mixed with handler logic |
| Dependency Injection | Built-in support | Requires manual implementation |
| Code Organization | Hierarchical and contextual | Flat and potentially scattered |


## 5. Getting Started

To start using GTW, follow these steps:

1. Install GTW (assuming it's available as a package):
   ```
   go get github.com/vedadiyan/gtw
   ```

2. Define your API structure:
   ```go
   type Foo struct {
       Get gtw.Handler `route:"/get" method:"GET"`
   }
   ```

3. Implement your handlers:
   ```go
   func (f *Foo) GetHandler(httpCtx *gtw.HttpCtx) (gtw.Status, gtw.Response) {
       // Your logic here
   }
   ```

4. Register your API with the default server:
   ```go
   func init() {
       gtw.DefaultServer().Register(new(TestAPI))
   }
   ```

5. Run your server (implementation details may vary):
   ```go
   func main() {
       err := gtw.DefaultServer().ListenAndServe(&http.Server{Addr: ":8080"})
       if err != nil {
           log.Fatal(err)
       }
   }
   ```

## 6. Features

- Annotation-driven API definition
- Struct-based organization for improved maintainability
- Built-in support for dependency injection
- Clear separation of route configuration and handler logic
- Flexible response handling
