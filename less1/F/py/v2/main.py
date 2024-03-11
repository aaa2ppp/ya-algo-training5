from sys import stdin, stdout


def main():
    n = int(stdin.readline().strip())

    odd_count_is_even = True
    first_odd_pos = -1

    a = map(int, stdin.readline().strip().split())

    for i, v in enumerate(a):
        if v&1:
            odd_count_is_even = not odd_count_is_even
            if first_odd_pos == -1:
                first_odd_pos = i

    res = bytearray(b'+'*(n-1))
    if odd_count_is_even:
        res[first_odd_pos] = ord('x')

    stdout.write(res.decode("utf-8"))


if __name__ == "__main__":
    main()