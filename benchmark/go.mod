module github.com/seeadoog/expr/benchmark

go 1.23

replace github.com/seeadoog/expr => ../

require (
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/antonmedv/expr v1.15.5
	github.com/seeadoog/expr v0.0.0-00010101000000-000000000000
)

require github.com/cespare/xxhash/v2 v2.3.0 // indirect
