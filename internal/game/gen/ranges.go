package gen

import "github.com/gurselcakar/arithmego/internal/game"

// Range defines min/max bounds for operand generation.
type Range struct {
	Min, Max int
}

// RangeTable maps difficulty to operand ranges.
type RangeTable map[game.Difficulty]Range

// MultiRange holds multiple operand ranges for multi-operand patterns.
type MultiRange struct {
	Primary   Range // First operand pair
	Secondary Range // Additional operands (3rd, 4th, etc.)
}

// MultiRangeTable maps difficulty to multi-operand ranges.
type MultiRangeTable map[game.Difficulty]MultiRange

// Standard addition ranges per difficulty.
var AdditionRanges = RangeTable{
	game.Beginner: {1, 9},
	game.Easy:     {10, 50},
	game.Medium:   {20, 200},
	game.Hard:     {100, 500},
	game.Expert:   {200, 999},
}

var AdditionMultiRanges = MultiRangeTable{
	game.Easy:   {Range{10, 50}, Range{5, 30}},
	game.Medium: {Range{20, 200}, Range{10, 80}},
	game.Hard:   {Range{100, 500}, Range{30, 200}},
	game.Expert: {Range{200, 999}, Range{50, 400}},
}

// Subtraction ranges
var SubtractionRanges = RangeTable{
	game.Beginner: {2, 9},
	game.Easy:     {10, 99},
	game.Medium:   {50, 300},
	game.Hard:     {100, 999},
	game.Expert:   {200, 9999},
}

// Multiplication ranges
var MultiplicationRanges = map[game.Difficulty][2]Range{
	game.Beginner: {{2, 9}, {2, 9}},
	game.Easy:     {{2, 12}, {10, 20}},
	game.Medium:   {{5, 30}, {5, 15}},
	game.Hard:     {{10, 50}, {10, 30}},
	game.Expert:   {{15, 99}, {15, 50}},
}

var MultiplicationMultiRanges = map[game.Difficulty]Range{
	game.Easy:   {2, 5},
	game.Medium: {2, 8},
	game.Hard:   {3, 10},
	game.Expert: {3, 12},
}

// Division ranges (divisor, quotient)
var DivisionRanges = map[game.Difficulty][2]Range{
	game.Beginner: {{2, 9}, {2, 9}},
	game.Easy:     {{2, 12}, {2, 12}},
	game.Medium:   {{3, 15}, {5, 20}},
	game.Hard:     {{5, 20}, {10, 30}},
	game.Expert:   {{10, 30}, {15, 50}},
}

// Square ranges
var SquareRanges = RangeTable{
	game.Beginner: {2, 10},
	game.Easy:     {5, 15},
	game.Medium:   {10, 20},
	game.Hard:     {15, 30},
	game.Expert:   {20, 50},
}

// Cube ranges
var CubeRanges = RangeTable{
	game.Beginner: {2, 5},
	game.Easy:     {2, 7},
	game.Medium:   {4, 10},
	game.Hard:     {6, 12},
	game.Expert:   {8, 15},
}

// SquareRoot ranges (root value, not the radicand)
var SquareRootRanges = RangeTable{
	game.Beginner: {2, 10},
	game.Easy:     {5, 15},
	game.Medium:   {10, 25},
	game.Hard:     {15, 35},
	game.Expert:   {25, 50},
}

// CubeRoot ranges (root value)
var CubeRootRanges = RangeTable{
	game.Beginner: {2, 5},
	game.Easy:     {3, 7},
	game.Medium:   {5, 10},
	game.Hard:     {7, 15},
	game.Expert:   {10, 20},
}

// Power ranges (base, exponent)
var PowerRanges = map[game.Difficulty][2]Range{
	game.Beginner: {{2, 10}, {2, 2}},
	game.Easy:     {{2, 12}, {2, 3}},
	game.Medium:   {{2, 10}, {2, 4}},
	game.Hard:     {{2, 8}, {3, 5}},
	game.Expert:   {{2, 6}, {4, 6}},
}

// Factorial ranges
var FactorialRanges = RangeTable{
	game.Beginner: {1, 4},
	game.Easy:     {3, 5},
	game.Medium:   {4, 6},
	game.Hard:     {5, 8},
	game.Expert:   {7, 10},
}

// Modulo ranges (divisor, dividend range)
var ModuloRanges = map[game.Difficulty][2]Range{
	game.Beginner: {{2, 9}, {3, 45}},
	game.Easy:     {{2, 12}, {5, 50}},
	game.Medium:   {{3, 15}, {20, 100}},
	game.Hard:     {{5, 25}, {50, 200}},
	game.Expert:   {{10, 50}, {100, 500}},
}

// Percentage easy/medium/hard/expert pools
var PercentEasy = []int{10, 20, 25, 50, 100}
var PercentMedium = []int{5, 10, 15, 20, 25, 30, 40, 50, 75}
var PercentHard = []int{5, 10, 12, 15, 20, 25, 30, 35, 40, 45, 50, 60, 75, 80}
var PercentExpert = []int{2, 3, 4, 6, 7, 8, 9, 11, 12, 13, 14, 16, 17, 18, 19,
	21, 22, 23, 24, 26, 27, 28, 29, 32, 33, 34, 36, 37, 38, 39}

// Percentage value ranges
var PercentValueRanges = map[game.Difficulty]Range{
	game.Beginner: {20, 100},
	game.Easy:     {10, 100},
	game.Medium:   {20, 200},
	game.Hard:     {50, 500},
	game.Expert:   {100, 1000},
}
