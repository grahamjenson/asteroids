package vector2d

type Edge struct {
	P1     Vector
	P2     Vector
	V      Vector
	Normal Vector
}

type Corner struct {
	p1 Vector
	p2 Vector
	p3 Vector
}

type Polygon struct {
	Matrix
	Debug bool
}

// x1,y1 ... xn,yn
func NewPolygon(points ...float64) *Polygon {
	if len(points)%2 != 0 {
		panic("Must init shape with even number of points")
	}
	pointsN := (len(points) / 2)
	nm := make([]Vector, pointsN)

	for i := 0; i < pointsN; i++ {
		nm[i] = Vector{points[(i * 2)], points[(i*2)+1], 1}
	}

	return &Polygon{Matrix: nm}
}

func (s Polygon) Clone() *Polygon {
	return &s
}

func (s *Polygon) Transform(m Matrix) *Polygon {
	// We always transform about the center of the
	x, y := s.Centroid()
	return s.TransformAbout(m, x, y)
}

func (s *Polygon) TransformAbout(m Matrix, x, y float64) *Polygon {
	// We always transform about the center of the
	aboutMatrix := CenterTransformationAt(m, x, y)
	ns := aboutMatrix.Product(s.Matrix)
	s.Matrix = ns
	return s
}

func (s *Polygon) Rotate(radian float64) *Polygon {
	s.Transform(RotateMatrix(radian))
	return s
}

func (s *Polygon) RotateAbout(radian float64, x, y float64) *Polygon {
	return s.TransformAbout(RotateMatrix(radian), x, y)
}

func (s *Polygon) Translate(dx, dy float64) *Polygon {
	s.Transform(TranslateMatrix(dx, dy))
	return s
}

func (s *Polygon) Scale(sx, sy float64) *Polygon {
	s.Transform(ScaleMatrix(sx, sy))
	return s
}

// also known as getAABB
func (s *Polygon) BoundingBox() (minX, minY, maxX, maxY float64) {
	minX = s.Matrix[0][0]
	minY = s.Matrix[0][1]

	maxX = s.Matrix[0][0]
	maxY = s.Matrix[0][1]

	for _, v := range s.Matrix {
		if v[0] > maxX {
			maxX = v[0]
		}
		if v[0] < minX {
			minX = v[0]
		}
		if v[1] > maxY {
			maxY = v[1]
		}
		if v[1] < minY {
			minY = v[1]
		}
	}
	return
}

func (s *Polygon) Edges() []Edge {
	len := len(s.Matrix)
	edges := make([]Edge, len)

	for i := 0; i < len; i++ {
		p1 := s.Matrix[i]
		p2 := s.Matrix[(i+1)%len] // Loop around if last point
		v := p1.Clone().Sub(p2)
		normal := v.Clone().Perp().Normalize()

		edges[i] = Edge{
			p1,
			p2,
			v,
			normal,
		}
	}

	return edges
}

func (s *Polygon) Corners() []Corner {
	len := len(s.Matrix)
	corners := make([]Corner, len)

	for i := 0; i < len; i++ {
		p1 := s.Matrix[i]
		p2 := s.Matrix[(i+1)%len]
		p3 := s.Matrix[(i+2)%len]

		corners[i] = Corner{p1, p2, p3}
	}

	return corners
}

func (s *Polygon) IsConvex() bool {
	// from https://stackoverflow.com/questions/471962/how-do-i-efficiently-determine-if-a-polygon-is-convex-non-convex-or-complex
	previous := 0
	for _, c := range s.Corners() {
		dx1 := c.p2[0] - c.p1[0]
		dy1 := c.p2[1] - c.p1[1]
		dx2 := c.p3[0] - c.p2[0]
		dy2 := c.p3[1] - c.p2[1]

		zcrossproduct := dx1*dy2 - dy1*dx2

		// zcrossproduct == 0 means
		if zcrossproduct > 0 {
			if previous == 0 || previous == +1 {
				previous = +1
			} else {
				return false
			}
		} else if zcrossproduct < 0 {
			if previous == 0 || previous == -1 {
				previous = -1
			} else {
				return false
			}
		}
	}

	return true
}

func (s *Polygon) Centroid() (float64, float64) {
	cx := 0.0
	cy := 0.0
	ar := 0.0

	edges := s.Edges()

	// handle edge case of it being a line
	if len(edges) == 2 {
		e := edges[0]
		p1 := e.P1
		p2 := e.P2
		return (p1[0] + p2[0]) / 2, (p1[1] + p2[1]) / 2
	}

	for _, e := range edges {
		p1 := e.P1
		p2 := e.P2

		a := (p1[0] * p2[1]) - (p2[0] * p1[1])
		cx += (p1[0] + p2[0]) * a
		cy += (p1[1] + p2[1]) * a
		ar += a
	}

	ar = ar * 3 // we want 1 / 6 the area and we currently have 2*area
	cx = cx / ar
	cy = cy / ar
	return cx, cy
}
