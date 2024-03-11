from sys import stdin, stdout


def main():
    n = int(stdin.readline().strip())

    f = [0]*n
    j, k = -1, n
    maximum = 0
    idx = -1
    height = 0

    for i in range(1, n+1):
        a, b = map(int, stdin.readline().strip().split())

        if a > b:
            height += a - b
            j += 1
            ii = j
        else:
            k -= 1
            ii = k

        f[ii] = i

        v = min(a, b)
        if v > maximum:
            maximum = v
            idx = ii

    height += maximum

    if idx < j:
        f[idx], f[j] = f[j], f[idx]
    elif idx > k:
        f[idx], f[k] = f[k], f[idx]


    stdout.write(f"{height}\n")
    stdout.write(" ".join(map(str, f)))


if __name__ == "__main__":
    main()
