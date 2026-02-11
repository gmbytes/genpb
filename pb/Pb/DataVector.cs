using System;

namespace Pb
{
    public sealed partial class Vector
    {
        private static readonly double Deg2Rad = Math.PI / 180.0;
        private static readonly double Rad2Deg = 180.0 / Math.PI;
        private static readonly Random RandomSource = new Random();

        public static readonly Vector ZeroVector = new Vector { X = 0, Y = 0, Z = 0 };
        public static readonly Vector ForwardVector = new Vector { X = 1, Y = 0, Z = 0 };
        public static readonly Vector OneVector = new Vector { X = 1, Y = 1, Z = 1 };

        public static Vector NewVector(double x, double y, double z)
        {
            return new Vector { X = x, Y = y, Z = z };
        }

        public string StringF()
        {
            return string.Format("({0:F5}, {1:F5}, {2:F5})", X, Y, Z);
        }

        public double ToAngle2D()
        {
            return Angle2D(ForwardVector);
        }

        public double Angle2D(Vector v)
        {
            return Radian2D(v) * Rad2Deg;
        }

        public double ToRadian2D()
        {
            return Radian2D(ForwardVector);
        }

        public double Radian2D(Vector v)
        {
            double sin = X * v.Y - v.X * Y;
            double cos = X * v.X + Y * v.Y;
            return -Math.Atan2(sin, cos);
        }

        public Vector Rotate2D(double alpha)
        {
            double sinA = Math.Sin(alpha);
            double cosA = Math.Cos(alpha);
            return new Vector
            {
                X = X * cosA - Y * sinA,
                Y = X * sinA + Y * cosA,
                Z = Z
            };
        }

        public Vector RotateAngle2D(double alphaDeg)
        {
            return Rotate2D(alphaDeg * Deg2Rad);
        }

        public double Dot2D(Vector v)
        {
            return X * v.X + Y * v.Y;
        }

        public double Dot(Vector v)
        {
            return X * v.X + Y * v.Y + Z * v.Z;
        }

        public Vector Cross(Vector v)
        {
            return new Vector
            {
                X = Y * v.Z - Z * v.Y,
                Y = Z * v.X - X * v.Z,
                Z = X * v.Y - Y * v.X
            };
        }

        public double LengthSq2D()
        {
            return X * X + Y * Y;
        }

        public double LengthSq()
        {
            return X * X + Y * Y + Z * Z;
        }

        public double Length2D()
        {
            return Math.Sqrt(LengthSq2D());
        }

        public double Length()
        {
            return Math.Sqrt(LengthSq());
        }

        public double DistanceSq2D(Vector v)
        {
            double dx = X - v.X;
            double dy = Y - v.Y;
            return dx * dx + dy * dy;
        }

        public double DistanceSq(Vector v)
        {
            double dx = X - v.X;
            double dy = Y - v.Y;
            double dz = Z - v.Z;
            return dx * dx + dy * dy + dz * dz;
        }

        public double Distance2D(Vector v)
        {
            return Math.Sqrt(DistanceSq2D(v));
        }

        public double Distance(Vector v)
        {
            return Math.Sqrt(DistanceSq(v));
        }

        public bool Equal2D(Vector v)
        {
            return X == v.X && Y == v.Y;
        }

        public bool Equal(Vector v)
        {
            return X == v.X && Y == v.Y && Z == v.Z;
        }

        public bool ApproximatelyEqual2D(Vector v)
        {
            return Math.Abs(X - v.X) < 1e-7 && Math.Abs(Y - v.Y) < 1e-7;
        }

        public bool ApproximatelyEqual(Vector v)
        {
            return Math.Abs(X - v.X) < 1e-7 && Math.Abs(Y - v.Y) < 1e-7 && Math.Abs(Z - v.Z) < 1e-7;
        }

        public Vector Orthogonal2D()
        {
            return new Vector
            {
                X = -Y,
                Y = X,
                Z = Z
            };
        }

        public Vector Copy()
        {
            return new Vector
            {
                X = X,
                Y = Y,
                Z = Z
            };
        }

        public Vector CopyNewZ(double z)
        {
            return new Vector
            {
                X = X,
                Y = Y,
                Z = z
            };
        }

        public void CopyTo(Vector dst)
        {
            dst.X = X;
            dst.Y = Y;
            dst.Z = Z;
        }

        public Vector Reverse2D()
        {
            return new Vector
            {
                X = -X,
                Y = -Y,
                Z = Z
            };
        }

        public Vector Reverse()
        {
            return new Vector
            {
                X = -X,
                Y = -Y,
                Z = -Z
            };
        }

        public Vector Add2D(Vector v)
        {
            return new Vector
            {
                X = X + v.X,
                Y = Y + v.Y,
                Z = Z
            };
        }

        public Vector Add(Vector v)
        {
            return new Vector
            {
                X = X + v.X,
                Y = Y + v.Y,
                Z = Z + v.Z
            };
        }

        public Vector Sub2D(Vector v)
        {
            return new Vector
            {
                X = X - v.X,
                Y = Y - v.Y,
                Z = Z
            };
        }

        public Vector Sub(Vector v)
        {
            return new Vector
            {
                X = X - v.X,
                Y = Y - v.Y,
                Z = Z - v.Z
            };
        }

        public Vector Mul2D(double v)
        {
            return new Vector
            {
                X = X * v,
                Y = Y * v,
                Z = Z
            };
        }

        public Vector Mul(double v)
        {
            return new Vector
            {
                X = X * v,
                Y = Y * v,
                Z = Z * v
            };
        }

        public Vector Div2D(double v)
        {
            if (v == 0)
            {
                return Copy();
            }
            double inv = 1.0 / v;
            return new Vector
            {
                X = X * inv,
                Y = Y * inv,
                Z = Z
            };
        }

        public Vector Div(double v)
        {
            if (v == 0)
            {
                return Copy();
            }
            double inv = 1.0 / v;
            return new Vector
            {
                X = X * inv,
                Y = Y * inv,
                Z = Z * inv
            };
        }

        public Vector Norm2D()
        {
            double lenSq = LengthSq2D();
            if (lenSq == 0)
            {
                return ForwardVector.Copy();
            }
            double l = 1.0 / Math.Sqrt(lenSq);
            return new Vector
            {
                X = X * l,
                Y = Y * l,
                Z = Z
            };
        }

        public Vector Norm()
        {
            double lenSq = LengthSq();
            if (lenSq == 0)
            {
                return Copy();
            }
            if (lenSq == 1)
            {
                return Copy();
            }
            double l = 1.0 / Math.Sqrt(lenSq);
            return new Vector
            {
                X = X * l,
                Y = Y * l,
                Z = Z * l
            };
        }

        public static Vector GenerateRandomVector(Vector min, Vector max)
        {
            return new Vector
            {
                X = RandomSource.NextDouble() * (max.X - min.X) + min.X,
                Y = RandomSource.NextDouble() * (max.Y - min.Y) + min.Y,
                Z = RandomSource.NextDouble() * (max.Z - min.Z) + min.Z
            };
        }
    }
}
