package postgres

import "fmt"

var ErrNotEnoughMoney = fmt.Errorf("not enough money")
var ErrUserNotFound = fmt.Errorf("user not found")
var ErrTransactionNotFound = fmt.Errorf("transaction not found")
