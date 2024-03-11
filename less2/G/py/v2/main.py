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


def main():
    lines = stdin.readlines()

    t = int(lines[0].strip())
    for i in range(2, len(lines), 2):
        #skip n = int(lines[i-1].strip())
        res = solution(map(int, lines[i].strip().split()))

        stdout.write(f"{len(res)}\n")
        stdout.write(" ".join(map(str, res)))
        stdout.write("\n")


if __name__ == "__main__":
    main()
