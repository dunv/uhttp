package models

type ParamRequirement struct {
	AllValues bool
	Date      bool
	ShortDate bool
	Enum      []string
	Int       bool
	Float     bool
	Bool      bool
}