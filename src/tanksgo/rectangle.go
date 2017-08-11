package main

type Rectangle struct {
	Points [4]Point // LeftUP, RightUP, RightDOWN, LeftDOWN
}

func (rectangle *Rectangle) intersects(r *Rectangle) bool {
	if rectangle.Points[RIGHTDOWN].X < r.Points[LEFTUP].X ||
		r.Points[RIGHTDOWN].X < rectangle.Points[LEFTUP].X ||
		rectangle.Points[RIGHTDOWN].Y < r.Points[LEFTUP].Y ||
		r.Points[RIGHTDOWN].Y < rectangle.Points[LEFTUP].Y {
		return false
	}
	return true
}

func (rectangle *Rectangle) pointInside(point Point) bool {
	if rectangle.Points[LEFTUP].X <= point.X && rectangle.Points[LEFTUP].Y <= point.Y &&
		rectangle.Points[RIGHTDOWN].X >= point.X && rectangle.Points[RIGHTDOWN].Y >= point.Y {
		return true
	}
	return false
}
