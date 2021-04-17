# go-blackscholes

Go implementation of the basic Black Scholes formulas for European option prices, greeks and implied volatility.

Edge cases such as e.g. zero volatility implemented exactly.

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
        Dividend:     0,
        Type:         bs.Straddle,
    }

    price, err := bs.Price(par)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Price: %.2f\n", price)
    
    // Version with argument list and no error checking:
    
    price = bs.BSPriceNoErrorCheck(0.2, 1, 100, 100, 0.02, 0, bs.Straddle)
    
    fmt.Printf("Price: %.2f\n", price)
}
```
