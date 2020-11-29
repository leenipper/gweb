package main

// colors
const (
	BGColor       = "#ecf0f1"
	BallColor     = "#27ae60"
	PlatformColor = "#2c3e50"
	TextColor     = PlatformColor
	FailColor     = "#c0392b"
	WinColor      = BallColor
)

// platform
const (
	PlatformWidth    = 120
	PlatformHeight   = 20
	PlatformMaxSpeed = 40
	PlatformAura     = 5 // additional invisible bounce space around the platform
)

// ball
const BallSize = 20
const InitialCredits = 2

// bricks
const (
	BrickHeight     = 20
	BrickRows       = 8
	BrickCols       = 14
	BrickMarginLeft = 120 // pixels
	BrickMarginTop  = 10  // pixels
	BrickMarginX    = 5   // pixels
	BrickMarginY    = 5   // pixels
)

// text box
const (
	TextWidth  = 90
	TextHeight = 20
	TextLeft   = 10
	TextTop    = 10
	TextMargin = 5

	TextBottom = TextTop + (TextHeight+TextMargin)*4
	TextRight  = TextLeft + TextWidth
)
