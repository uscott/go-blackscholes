# go-blackscholes

Go implementation of the basic Black Scholes formulas for European option prices, greeks and implied volatility.

Edge cases such as e.g. zero volatility implemented exactly.

Use at your own risk.

### Install
```shell
go get github.com/uscott/go-blackscholes@latest
```

### Usage

> Refer to test cases for more examples

```go
package main

import (
    "fmt"
    bs "github.com/uscott/go-blackscholes"
)

func main() {

    vol := 0.2
    timeToExpiry := 1.0
    spot := 100.0
    strike := 100.0
    interestRate := 0.02
    dividendYield := 0.01
    optionType := blackscholes.Straddle

    price, err := blackscholes.Price(vol, timeToExpiry, spot, strike, interestRate, dividendYield, optionType)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Price: %.2f\n", price)
}
```
