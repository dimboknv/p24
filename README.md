# `p24` [![Build Status](https://github.com/dimboknv/p24/workflows/build/badge.svg)](https://github.com/dimboknv/p24/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/dimboknv/p24)](https://goreportcard.com/report/github.com/dimboknv/p24) [![Coverage Status](https://coveralls.io/repos/github/dimboknv/p24/badge.svg?branch=master)](https://coveralls.io/github/dimboknv/p24?branch=master) 

This library provides [privat24 marchant information api](https://api.privatbank.ua/#p24/main) client.

**Note:** before using `p24` you need [to register merchant in privat24 system](https://api.privatbank.ua/#p24/registration).

## Install

```shell
go get -u github.com/dimboknv/p24
```

## Usage

Also, you can visit [p24-cli](https://github.com/dimboknv/p24-cli) for an example with retryable, rate limited, logged client.


```go
package main

import (
	"fmt"
	"log"
	"context"
	"net/http"
	"time"
	"github.com/dimboknv/p24"
)

func main() {
	client := p24.NewClient(p24.ClientOpts{
		HTTP: &http.Client{},
		Merchant: p24.Merchant{
			ID:   "merchant id",
			Pass: "merchant pass",
		},
	})

	// get merchant statements list for 02.01.2021 - 02.02.2021 date range
	// and "1234567891234567" card number
	startDate, _ := time.Parse("02.01.2006", "02.01.2021")
	endDate, _ := time.Parse("02.01.2006", "02.02.2021")
	ctx := context.Background()
	statements, err := client.GetStatements(ctx, p24.StatementsOpts{
		StartDate:  startDate,
		EndDate:    endDate,
		CardNumber: "1234567891234567",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(statements)

	// get merchant card balance for "1234567891234567" card number
	cardBalance, err := client.GetCardBalance(ctx, p24.BalanceOpts{
		CardNumber: "1234567891234567",
		Country:    "UA",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cardBalance)
}
```
