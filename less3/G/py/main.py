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

    for p1 in points:
        for p2 in points:
            if p1 != p2:
                dx = p2[X] - p1[X]
                dy = p2[Y] - p1[Y]

                tmp = [pp for pp in ((p[X]-dy, p[Y]+dx) for p in (p1, p2)) if pp not in point_set]

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
