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

    i = 0
    if odd_count_is_even:
        for _ in range(first_odd_pos):
            stdout.write("+")
        stdout.write("x")
        i = first_odd_pos + 1

    for _ in range(i, n-1):
        stdout.write("+")


if __name__ == "__main__":
    main()