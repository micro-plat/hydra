package types

import "github.com/shopspring/decimal"

//Decimal decimal.Decimal的别名
type Decimal = decimal.Decimal

var NewDecimal = decimal.New

var NewDecimalFromInt = decimal.NewFromInt
var NewDecimalFromInt32 = decimal.NewFromInt32
var NewDecimalFromBigInt = decimal.NewFromBigInt
var NewDecimalFromString = decimal.NewFromString
var NewDecimalFromFloat = decimal.NewFromFloat
var NewDecimalFromFloat32 = decimal.NewFromFloat32
var NewDecimalFromFloatWithExponent = decimal.NewFromFloatWithExponent
