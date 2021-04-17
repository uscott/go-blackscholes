# go-blackscholes

Go implementation of the basic Black Scholes formulas for European option prices, greeks and implied volatility.

Use at your own risk.

### Install
```shell script
go get github.com/uscott/go-blackscholes
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

    pars := &bs.PriceParams{
        Vol:          0.2,
        TimeToExpiry: 1,
        Underlying:   100,
        Strike:       100,
        Rate:         0.02,
        Divident:     0,
        Type:         bs.Straddle,
    }

    price, err := bs.Price(par)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Price: %.2f\n", price)
}
```
