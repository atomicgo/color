<h1 align="center">AtomicGo | color</h1>

<p align="center">
<img src="https://img.shields.io/endpoint?url=https%3A%2F%2Fatomicgo.dev%2Fapi%2Fshields%2Fcolor&style=flat-square" alt="Downloads">

<a href="https://github.com/atomicgo/color/releases">
<img src="https://img.shields.io/github/v/release/atomicgo/color?style=flat-square" alt="Latest Release">
</a>

<a href="https://codecov.io/gh/atomicgo/color" target="_blank">
<img src="https://img.shields.io/github/actions/workflow/status/atomicgo/color/go.yml?style=flat-square" alt="Tests">
</a>

<a href="https://codecov.io/gh/atomicgo/color" target="_blank">
<img src="https://img.shields.io/codecov/c/gh/atomicgo/color?color=magenta&logo=codecov&style=flat-square" alt="Coverage">
</a>

<a href="https://codecov.io/gh/atomicgo/color">
<!-- unittestcount:start --><img src="https://img.shields.io/badge/Unit_Tests-0-magenta?style=flat-square" alt="Unit test count"><!-- unittestcount:end -->
</a>

<a href="https://opensource.org/licenses/MIT" target="_blank">
<img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
</a>
  
<a href="https://goreportcard.com/report/github.com/atomicgo/color" target="_blank">
<img src="https://goreportcard.com/badge/github.com/atomicgo/color?style=flat-square" alt="Go report">
</a>   

</p>

---

<p align="center">
<strong><a href="https://pkg.go.dev/atomicgo.dev/color#section-documentation" target="_blank">Documentation</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CONTRIBUTING.md" target="_blank">Contributing</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CODE_OF_CONDUCT.md" target="_blank">Code of Conduct</a></strong>
</p>

---

<p align="center">
  <img src="https://raw.githubusercontent.com/atomicgo/atomicgo/main/assets/header.png" alt="AtomicGo">
</p>

<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>
<h3  align="center"><pre>go get atomicgo.dev/color</pre></h3>
<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# color

```go
import "atomicgo.dev/color"
```

Package color is used to generate new AtomicGo repositories.

Write the description of the module here. You can use \*\*markdown\*\*\! This description should clearly explain what the package does.

Example description: https://golang.org/src/encoding/gob/doc.go

<details><summary>Example (Demo)</summary>
<p>



```go
package main

import (
	"fmt"

	"atomicgo.dev/color"
)

func main() {
	// Simple coloring
	fmt.Println("Hello, " + color.Red("World") + "!")

	// Theme colors - can be customized in init() function if needed
	fmt.Println(color.Success("Success message"))
	fmt.Println(color.Info("Info message"))
	fmt.Println(color.Warning("Warning message"))
	fmt.Println(color.Error("Error message"))
	fmt.Println(color.Fatal("Fatal message"))
}
```

</p>
</details>

## Index

- [Variables](<#variables>)
- [func Sprint\(a ...any\) string](<#Sprint>)
- [type ANSI256Color](<#ANSI256Color>)
  - [func \(c ANSI256Color\) Sequence\(background bool\) string](<#ANSI256Color.Sequence>)
  - [func \(c ANSI256Color\) String\(\) string](<#ANSI256Color.String>)
- [type ANSIColor](<#ANSIColor>)
  - [func \(c ANSIColor\) Sequence\(background bool\) string](<#ANSIColor.Sequence>)
- [type Color](<#Color>)
  - [func NewColorFromHex\(hex string\) Color](<#NewColorFromHex>)
- [type Mode](<#Mode>)
  - [func \(m Mode\) String\(\) string](<#Mode.String>)
- [type Modifier](<#Modifier>)
  - [func \(m Modifier\) Sequence\(\) string](<#Modifier.Sequence>)
- [type RGBColor](<#RGBColor>)
  - [func NewColorFromRGB\(r, g, b uint8\) RGBColor](<#NewColorFromRGB>)
  - [func \(c RGBColor\) Hex\(\) string](<#RGBColor.Hex>)
  - [func \(c RGBColor\) Sequence\(background bool\) string](<#RGBColor.Sequence>)
- [type Style](<#Style>)
  - [func NewStyle\(foregroundColor, backgroundColor Color, modifiers ...Modifier\) Style](<#NewStyle>)
  - [func \(s Style\) Sequence\(\) string](<#Style.Sequence>)
  - [func \(s Style\) Sprint\(a ...any\) string](<#Style.Sprint>)


## Variables

<a name="Black"></a>

```go
var (
    // ANSI colors
    Black       = NewStyle(ANSIBlack, nil).Sprint
    BrightBlack = NewStyle(ANSIBrightBlack, nil).Sprint

    Red       = NewStyle(ANSIRed, nil).Sprint
    BrightRed = NewStyle(ANSIBrightRed, nil).Sprint

    Green       = NewStyle(ANSIGreen, nil).Sprint
    BrightGreen = NewStyle(ANSIBrightGreen, nil).Sprint

    Yellow       = NewStyle(ANSIYellow, nil).Sprint
    BrightYellow = NewStyle(ANSIBrightYellow, nil).Sprint

    Blue       = NewStyle(ANSIBlue, nil).Sprint
    BrigthBlue = NewStyle(ANSIBrightBlue, nil).Sprint

    Magenta       = NewStyle(ANSIMagenta, nil).Sprint
    BrightMagenta = NewStyle(ANSIBrightMagenta, nil).Sprint

    Cyan       = NewStyle(ANSICyan, nil).Sprint
    BrightCyan = NewStyle(ANSIBrightCyan, nil).Sprint

    White       = NewStyle(ANSIWhite, nil).Sprint
    BrightWhite = NewStyle(ANSIBrightWhite, nil).Sprint

    // Special colors
    Success = NewStyle(ANSIBrightGreen, nil).Sprint
    Info    = NewStyle(ANSIBrightBlue, nil).Sprint
    Warning = NewStyle(ANSIBrightYellow, nil).Sprint
    Error   = NewStyle(ANSIBrightRed, nil).Sprint
    Fatal   = NewStyle(ANSIBrightRed, nil, Bold).Sprint
)
```

<a name="Sprint"></a>
## func [Sprint](<https://github.com/atomicgo/color/blob/main/color.go#L11>)

```go
func Sprint(a ...any) string
```



<a name="ANSI256Color"></a>
## type [ANSI256Color](<https://github.com/atomicgo/color/blob/main/color-ansi256.go#L3>)



```go
type ANSI256Color uint8
```

<a name="ANSI256Color.Sequence"></a>
### func \(ANSI256Color\) [Sequence](<https://github.com/atomicgo/color/blob/main/color-ansi256.go#L9>)

```go
func (c ANSI256Color) Sequence(background bool) string
```



<a name="ANSI256Color.String"></a>
### func \(ANSI256Color\) [String](<https://github.com/atomicgo/color/blob/main/color-ansi256.go#L5>)

```go
func (c ANSI256Color) String() string
```



<a name="ANSIColor"></a>
## type [ANSIColor](<https://github.com/atomicgo/color/blob/main/color-ansi.go#L24>)



```go
type ANSIColor int
```

<a name="ANSIBlack"></a>

```go
const (
    ANSIBlack ANSIColor = iota
    ANSIRed
    ANSIGreen
    ANSIYellow
    ANSIBlue
    ANSIMagenta
    ANSICyan
    ANSIWhite
    ANSIBrightBlack
    ANSIBrightRed
    ANSIBrightGreen
    ANSIBrightYellow
    ANSIBrightBlue
    ANSIBrightMagenta
    ANSIBrightCyan
    ANSIBrightWhite
)
```

<a name="ANSIColor.Sequence"></a>
### func \(ANSIColor\) [Sequence](<https://github.com/atomicgo/color/blob/main/color-ansi.go#L26>)

```go
func (c ANSIColor) Sequence(background bool) string
```



<a name="Color"></a>
## type [Color](<https://github.com/atomicgo/color/blob/main/color.go#L7-L9>)



```go
type Color interface {
    Sequence(background bool) string
}
```

<a name="NoColor"></a>

```go
var NoColor Color = noColor{}
```

<a name="NewColorFromHex"></a>
### func [NewColorFromHex](<https://github.com/atomicgo/color/blob/main/color-rgb.go#L25>)

```go
func NewColorFromHex(hex string) Color
```

NewColorFromHex creates a new Color from a hex string. If the hex string is invalid, NoColor is returned.

<a name="Mode"></a>
## type [Mode](<https://github.com/atomicgo/color/blob/main/mode.go#L3>)



```go
type Mode int
```

<a name="TrueColor"></a>

```go
const (
    TrueColor Mode = iota
    ANSI256
    ANSI

    Disabled
)
```

<a name="Mode.String"></a>
### func \(Mode\) [String](<https://github.com/atomicgo/color/blob/main/mode.go#L13>)

```go
func (m Mode) String() string
```



<a name="Modifier"></a>
## type [Modifier](<https://github.com/atomicgo/color/blob/main/modifier.go#L6>)

Modifier type for text modifiers

```go
type Modifier int
```

<a name="Reset"></a>Modifiers

```go
const (
    Reset Modifier = iota
    Bold
    Italic
    Underline
)
```

<a name="Modifier.Sequence"></a>
### func \(Modifier\) [Sequence](<https://github.com/atomicgo/color/blob/main/modifier.go#L18>)

```go
func (m Modifier) Sequence() string
```



<a name="RGBColor"></a>
## type [RGBColor](<https://github.com/atomicgo/color/blob/main/color-rgb.go#L9-L11>)



```go
type RGBColor struct {
    R, G, B uint8
}
```

<a name="NewColorFromRGB"></a>
### func [NewColorFromRGB](<https://github.com/atomicgo/color/blob/main/color-rgb.go#L19>)

```go
func NewColorFromRGB(r, g, b uint8) RGBColor
```

NewColorFromRGB creates a new RGBColor.

<a name="RGBColor.Hex"></a>
### func \(RGBColor\) [Hex](<https://github.com/atomicgo/color/blob/main/color-rgb.go#L14>)

```go
func (c RGBColor) Hex() string
```

Hex returns the hex representation of the color.

<a name="RGBColor.Sequence"></a>
### func \(RGBColor\) [Sequence](<https://github.com/atomicgo/color/blob/main/color-rgb.go#L47>)

```go
func (c RGBColor) Sequence(background bool) string
```



<a name="Style"></a>
## type [Style](<https://github.com/atomicgo/color/blob/main/style.go#L8-L13>)



```go
type Style struct {
    Foreground Color
    Background Color

    Modifiers []Modifier
}
```

<a name="NewStyle"></a>
### func [NewStyle](<https://github.com/atomicgo/color/blob/main/style.go#L15>)

```go
func NewStyle(foregroundColor, backgroundColor Color, modifiers ...Modifier) Style
```



<a name="Style.Sequence"></a>
### func \(Style\) [Sequence](<https://github.com/atomicgo/color/blob/main/style.go#L31>)

```go
func (s Style) Sequence() string
```



<a name="Style.Sprint"></a>
### func \(Style\) [Sprint](<https://github.com/atomicgo/color/blob/main/style.go#L49>)

```go
func (s Style) Sprint(a ...any) string
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->

---

> [AtomicGo.dev](https://atomicgo.dev) &nbsp;&middot;&nbsp;
> with ❤️ by [@MarvinJWendt](https://github.com/MarvinJWendt) |
> [MarvinJWendt.com](https://marvinjwendt.com)
