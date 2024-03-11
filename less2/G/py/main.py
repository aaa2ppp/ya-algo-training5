from sys import stdin, stdout


def solution(a):
    res = []
    minimum = a.__next__() 
    size = 1

    for v in a:
        size += 1
        minimum = min(minimum, v)

        if minimum < size:
            res.append(size-1)
            minimum = v
            size = 1

    res.append(size)
    
    return res


def task():
    n = int(stdin.readline().strip())
    a = map(int, stdin.readline().strip().split())
    res = solution(a)

    stdout.write(f"{len(res)}\n")
    stdout.write(" ".join(map(str, res)))
    stdout.write("\n")


def main():
    t = int(stdin.readline().strip())
    for _ in range(t):
        task()


if __name__ == "__main__":
    main()
