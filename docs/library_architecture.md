# Architecture

## QA

### Why is there no malloc?

The library was designed to have a simple API.

For example: The `Int(int) int, ptr` function returns an `int` and a pointer to it.

If `Int(int)` returned just a pointer to an `int`,
the user would have to dereference the pointer to get the `int` value. Meaning more complex user-handling.

Since the user has to dereference the pointer after, a `FreeInt(ptr)` function is provided to free the allocated memory
for an integer.