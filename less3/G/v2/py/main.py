from sys import stdin, stdout

X = 0
Y = 1


def solution(points):
    p = points[0]
    res = [
        (p[X], p[Y]+1),
        (p[X]+1, p[Y]),
        (p[X]+1, p[Y]+1)
    ]
    
    point_set = set(points)

    for i, p1 in enumerate(points):
        for j in range(i+1, len(points)):
            p2 = points[j]

            dx = p2[X] - p1[X]
            dy = p2[Y] - p1[Y]

            dy2 = dx + dy
            if dy2 % 2 != 0:
                continue

            dy2 //= 2
            dx2 = (dy-dx)//2

            p3 = (p1[X]-dx2, p1[Y]+dy2)
            p4 = (p2[X]+dx2, p2[Y]-dy2)

            tmp = [p for p in (p3, p4) if p not in point_set]

            if len(tmp) < len(res):
                res = tmp
                if len(res) == 0:
                    return []
    
    return res


def main():
    n = int(stdin.readline().strip())

    points = [None]*n
    for i in range(n):
        x, y = map(int, stdin.readline().strip().split())
        points[i] = (x, y)

    res = solution(points)

    stdout.write(f"{len(res)}\n")
    for p in res:
        stdout.write(" ".join(map(str, p)))
        stdout.write("\n")


if __name__ == "__main__":
    main()
