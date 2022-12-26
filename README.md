# The «dp» programming language

The language uses imperative style.  Its features are designed to work together
so that memory management can be automatic, simple and fast.

The language has null (nil) pointers.  Memory access operations don't require
elaborate syntax.

The `dpfmt` formatting tool is very lenient.  Therefore source code files stay
consistently formatted, but it's easy to write them.


## Details

- Type names start with an upper-case letter and consist of upper- and
  lower-case letters.

- Function and variable names consist of lower-case letters and non-consecutive
  underscores.

- Only fixed-size primitive types: `Bool`, `I8`, `I16`, `I32`, `I64`, `U8`,
  `U16`, `U32`, `U64`, `F32`, `F64`.  Expressions can yield unnamed integer
  types with arbitrary bit width; they must be converted to a named type when
  assigning to an explicitly typed variable, parameter or return value.

- Generic variable-length array type: `[T]`.

- Zero value literals: `0`, `{}`, `nil`, and `_`.

- Struct types with methods.  Field accessibility modes: hidden (default),
  `visible`, `mutable` (also visible), and `assignable` (also visible and
  mutable).

- Struct types can overload array index operator (unnamed method `(=StructType)
  (IndexType) =ItemType`).

- Packages are imported using URI strings.  Imported packages may be accessed
  using namespaces (`package::Symbol`) or specific symbols can be imported into
  scope (`Symbol`).

- Assignment syntax supports multiple targets and values.

- Function or method call expressions can appear in assignment target list.

- Expression `clone x`.


## Type modifier reference

```
 Type   Description

    T   non-null                          mutable  owned value
   *T   nullable           pointer   to   mutable  owned value
   =T   non-null temporary reference to   mutable        value
  =*T   nullable temporary pointer   to   mutable        value
   &T   non-null           reference to immutable        value
  *&T   nullable           pointer   to immutable        value
   #T   non-null                        immutable shared value
  *#T   nullable           pointer   to immutable shared value
  &#T   non-null           reference to immutable shared value
 *&#T   nullable           pointer   to immutable shared value
```
