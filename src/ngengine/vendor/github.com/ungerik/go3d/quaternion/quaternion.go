// Package quaternion contains a float32 unit-quaternion type T and functions.
package quaternion

import (
	"fmt"

	math "github.com/barnex/fmath"
	"github.com/ungerik/go3d/vec3"
	"github.com/ungerik/go3d/vec4"
)

var (
	// Zero holds a zero quaternion.
	Zero = T{}

	// Ident holds an ident quaternion.
	Ident = T{0, 0, 0, 1}
)

// T represents a orientatin/rotation as a unit quaternion.
// See http://en.wikipedia.org/wiki/Quaternions_and_spatial_rotation
type T [4]float32

// FromAxisAngle returns a quaternion representing a rotation around and axis.
func FromAxisAngle(axis *vec3.T, angle float32) T {
	angle *= 0.5
	sin := math.Sin(angle)
	q := T{axis[0] * sin, axis[1] * sin, axis[2] * sin, math.Cos(angle)}
	return q.Normalized()
}

// FromXAxisAngle returns a quaternion representing a rotation around the x axis.
func FromXAxisAngle(angle float32) T {
	angle *= 0.5
	return T{math.Sin(angle), 0, 0, math.Cos(angle)}
}

// FromYAxisAngle returns a quaternion representing a rotation around the y axis.
func FromYAxisAngle(angle float32) T {
	angle *= 0.5
	return T{0, math.Sin(angle), 0, math.Cos(angle)}
}

// FromZAxisAngle returns a quaternion representing a rotation around the z axis.
func FromZAxisAngle(angle float32) T {
	angle *= 0.5
	return T{0, 0, math.Sin(angle), math.Cos(angle)}
}

// FromEulerAngles returns a quaternion representing Euler angle rotations.
func FromEulerAngles(yHead, xPitch, zRoll float32) T {
	qy := FromYAxisAngle(yHead)
	qx := FromXAxisAngle(xPitch)
	qz := FromZAxisAngle(zRoll)
	return Mul3(&qy, &qx, &qz)
}

// FromVec4 converts a vec4.T into a quaternion.
func FromVec4(v *vec4.T) T {
	return T(*v)
}

// Vec4 converts the quaternion into a vec4.T.
func (quat *T) Vec4() vec4.T {
	return vec4.T(*quat)
}

// Parse parses T from a string. See also String()
func Parse(s string) (r T, err error) {
	_, err = fmt.Sscan(s, &r[0], &r[1], &r[2], &r[3])
	return r, err
}

// String formats T as string. See also Parse().
func (quat *T) String() string {
	return fmt.Sprint(quat[0], quat[1], quat[2], quat[3])
}

// AxisAngle extracts the rotation in form of an axis and a rotation angle.
func (quat *T) AxisAngle() (axis vec3.T, angle float32) {
	cos := quat[3]
	sin := math.Sqrt(1 - cos*cos)
	angle = math.Acos(cos)

	var ooSin float32
	if math.Abs(sin) < 0.0005 {
		ooSin = 1
	} else {
		ooSin = 1 / sin
	}
	axis[0] = quat[0] * ooSin
	axis[1] = quat[1] * ooSin
	axis[2] = quat[2] * ooSin

	return axis, angle
}

// Norm returns the norm value of the quaternion.
func (quat *T) Norm() float32 {
	return quat[0]*quat[0] + quat[1]*quat[1] + quat[2]*quat[2] + quat[3]*quat[3]
}

// Normalize normalizes to a unit quaternation.
func (quat *T) Normalize() *T {
	norm := quat.Norm()
	if norm != 1 && norm != 0 {
		ool := 1 / math.Sqrt(norm)
		quat[0] *= ool
		quat[1] *= ool
		quat[2] *= ool
		quat[3] *= ool
	}
	return quat
}

// Normalized returns a copy normalized to a unit quaternation.
func (quat *T) Normalized() T {
	norm := quat.Norm()
	if norm != 1 && norm != 0 {
		ool := 1 / math.Sqrt(norm)
		return T{
			quat[0] * ool,
			quat[1] * ool,
			quat[2] * ool,
			quat[3] * ool,
		}
	} else {
		return *quat
	}
}

// Negate negates the quaternion.
func (quat *T) Negate() *T {
	quat[0] = -quat[0]
	quat[1] = -quat[1]
	quat[2] = -quat[2]
	quat[3] = -quat[3]
	return quat
}

// Negated returns a negated copy of the quaternion.
func (quat *T) Negated() T {
	return T{-quat[0], -quat[1], -quat[2], -quat[3]}
}

// Invert inverts the quaterion.
func (quat *T) Invert() *T {
	quat[0] = -quat[0]
	quat[1] = -quat[1]
	quat[2] = -quat[2]
	return quat
}

// Inverted returns an inverted copy of the quaternion.
func (quat *T) Inverted() T {
	return T{-quat[0], -quat[1], -quat[2], quat[3]}
}

// SetShortestRotation negates the quaternion if it does not represent the shortest rotation from quat to the orientation of other.
// (there are two directions to rotate from the orientation of quat to the orientation of other)
// See IsShortestRotation()
func (quat *T) SetShortestRotation(other *T) *T {
	if !IsShortestRotation(quat, other) {
		quat.Negate()
	}
	return quat
}

// IsShortestRotation returns if the rotation from a to b is the shortest possible rotation.
// (there are two directions to rotate from the orientation of quat to the orientation of other)
// See T.SetShortestRotation
func IsShortestRotation(a, b *T) bool {
	return Dot(a, b) >= 0
}

// IsUnitQuat returns if the quaternion is within tolerance of the unit quaternion.
func (quat *T) IsUnitQuat(tolerance float32) bool {
	norm := quat.Norm()
	return norm >= (1.0-tolerance) && norm <= (1.0+tolerance)
}

// RotateVec3 rotates v by the rotation represented by the quaternion.
func (quat *T) RotateVec3(v *vec3.T) {
	qv := T{v[0], v[1], v[2], 0}
	inv := quat.Inverted()
	q := Mul3(quat, &qv, &inv)
	v[0] = q[0]
	v[1] = q[1]
	v[2] = q[2]
}

// RotatedVec3 returns a rotated copy of v.
func (quat *T) RotatedVec3(v *vec3.T) vec3.T {
	qv := T{v[0], v[1], v[2], 0}
	inv := quat.Inverted()
	q := Mul3(quat, &qv, &inv)
	return vec3.T{q[0], q[1], q[2]}
}

// Dot returns the dot product of two quaternions.
func Dot(a, b *T) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3]
}

// Mul multiplies two quaternions.
func Mul(a, b *T) T {
	q := T{
		a[3]*b[0] + a[0]*b[3] + a[1]*b[2] - a[2]*b[1],
		a[3]*b[1] + a[1]*b[3] + a[2]*b[0] - a[0]*b[2],
		a[3]*b[2] + a[2]*b[3] + a[0]*b[1] - a[1]*b[0],
		a[3]*b[3] - a[0]*b[0] - a[1]*b[1] - a[2]*b[2],
	}
	return q.Normalized()
}

// Mul3 multiplies three quaternions.
func Mul3(a, b, c *T) T {
	q := Mul(a, b)
	return Mul(&q, c)
}

// Mul4 multiplies four quaternions.
func Mul4(a, b, c, d *T) T {
	q := Mul(a, b)
	q = Mul(&q, c)
	return Mul(&q, d)
}

// Slerp returns the spherical linear interpolation quaternion between a and b at t (0,1).
// See http://en.wikipedia.org/wiki/Slerp
func Slerp(a, b *T, t float32) T {
	d := math.Acos(a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3])
	ooSinD := 1 / math.Sin(d)

	t1 := math.Sin(d*(1-t)) * ooSinD
	t2 := math.Sin(d*t) * ooSinD

	q := T{
		a[0]*t1 + b[0]*t2,
		a[1]*t1 + b[1]*t2,
		a[2]*t1 + b[2]*t2,
		a[3]*t1 + b[3]*t2,
	}

	return q.Normalized()
}

// Vec3Diff returns the rotation quaternion between two vectors.
func Vec3Diff(a, b *vec3.T) T {
	cr := vec3.Cross(a, b)
	sr := math.Sqrt(2 * (1 + vec3.Dot(a, b)))
	oosr := 1 / sr

	q := T{cr[0] * oosr, cr[1] * oosr, cr[2] * oosr, sr * 0.5}
	return q.Normalized()
}
