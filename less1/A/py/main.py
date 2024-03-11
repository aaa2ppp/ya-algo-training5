from sys import stdin, stdout


def main():
    p, v = map(int, stdin.readline().strip().split())
    q, m = map(int, stdin.readline().strip().split())

    p1 = p - v
    p2 = p + v

    q1 = q - m
    q2 = q + m

    res = 0
    if p2 >= q1 and q2 >= p1:
        res = max(p2, q2) - min(p1, q1) + 1
    else:
        res = (p2 - p1) + (q2 - q1) + 2

    stdout.write(f"{res}\n")


if __name__ == "__main__":
    main()
