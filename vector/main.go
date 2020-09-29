package main

// Notes ...
// This implementation shows a few different ways of initializing new vector
// and modifying already existing ones.
// Obviously when each of these implementation has it's benefits and pitfals.
// This hopefully shows that I am aware

// A basic vector implementation in golang
type Vector struct {
	X, Y, Z float64
}

// A cross product implementation where the vector it self is changed
func (a *Vector) Add(b *Vector) {
	a.X += b.X
	a.Y += b.Y
	a.Z += b.Z
}

// A subtraction implementation that returns a new vector
func (a *Vector) Sub(b *Vector) (c Vector) {
	c.X = a.X - b.X
	c.Y = a.Y - b.Y
	c.Z = a.Z - b.Z
	return
}

// A dot product implementation that modifies the current vector
func (a *Vector) Dot(b *Vector) {
	a.X = a.X * b.X
	a.Y = a.Y * b.Y
	a.Z = a.Z * b.Z
}

// A scalar multiply method that returns a pointer to a new vector
func (a *Vector) ScalarMultiply(b float64) *Vector {
	c := new(Vector)
	c.X = a.X * b
	c.Y = a.Y * b
	c.Z = a.Z * b
	return c
}

// A cross product implementation that returns a new vector
func (a *Vector) Cross(b *Vector) Vector {
	return Vector{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}
