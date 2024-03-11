from sys import stdin, stdout

def check(i, x, z):
    while True:
        if z <= 0:
            return i
        
        x -= min(x, z)
        if x <= 0:
            return -1
        
        z -= min(x, z)
        i += 1


def solution(x, y, z, p):
    res = -1
    i = 1
    while res == -1 or i < res:
        if x >= y:
            i2 = check(i, x, z-(x-y))
            if i2 == i:
                return i2
            
            if i2 != -1 and (res == -1 or i2 < res):
                res = i2

        x2 = x
        if y > 0:
            y -= 1
            x2 -= 1
        
        d = min(x2, z)
        z -= d
        x2 -= d

        y -= min(x2, y)

        if y <= 0 and z <= 0:
            return i
        
        x -= min(x, z)
        if x <= 0:
            break

        if y > 0:
            z += p

        i += 1

    return res


def main():
    x, y, p = (int(s.strip()) for s in stdin.read().strip().split("\n"))
    res = solution(x, y, 0, p)
    stdout.write(str(res))


if __name__ == "__main__":
    main()