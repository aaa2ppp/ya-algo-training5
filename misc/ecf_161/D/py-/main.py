from sys import stdin, stdout

def main():
    t = int(stdin.readline().strip())

    for i in range(t):
        task()


def task():
    n = int(stdin.readline().strip())
    a = list(map(int, stdin.readline().strip().split()))
    d = list(map(int, stdin.readline().strip().split()))
    print(*solve(n, a, d))


def solve(n, a, d):
    res = [0]*n
    m = len(a)
    b = [0]*(m+2)

    for ii in range(n):
        for i in range(m+2):
            b[i] = 0

        for i in range(m):
            b[i]   += a[i]
            b[i+2] += a[i]
        
        i, j = 0, 0
        while i < m:
            if d[i] >= b[i+1]:
                a[j] = a[i]
                d[j] = d[i]
                j += 1
            i += 1

        if i-j == 0:
            break

        res[ii] = i-j
        m = j

    return res


if __name__ == "__main__":
    main()
