#include <vector>
#include <iostream>

using namespace std;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);    

    int p, v, q, m;

    cin >> p >> v >> q >> m;

    int p1 = p - v;
    int p2 = p + v;

    int q1 = q - m;
    int q2 = q + m;

    int cnt;
    if (p1 <= q2 && q1 <= p2) {
        cnt = max(p2, q2) - min(p1, q1) + 1;
    } else {
        cnt = p2 - p1 + q2 - q1 + 2;
    }

    cout << cnt;

    return 0;
}
