#include <iostream>
#include <string>
#include <vector>
#include <algorithm>

using namespace std;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);    

    string s1, s2;
    int f;

    cin >> s1 >> s2 >> f;

    int p;
    int g11, g12, g21, g22, res;
    
    p = s1.find(':');
    g11 = stoi(s1.substr(0, p));
    g12 = stoi(s1.substr(p+1));

    p = s2.find(':');
    g21 = stoi(s1.substr(0, p));
    g22 = stoi(s1.substr(p+1));

    res = (g12+g22) - (g11+g21+1);

    g21 += res;

    if (res < 0) {
        res = 0;
    }

    if (res > 0) {
        switch (f) {
        case 1:
            if ((g21-1 > g22) && !(g12 > g11)) {
                res--;
            }
            break;
        case 2:
            if ((g11 > g12) && !(g22 > g21-1)) {
                res--;
            }
            break;            
        }
    }

    cout << res;

    return 0;
}
