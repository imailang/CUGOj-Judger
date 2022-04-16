#include <bits/stdc++.h>
using namespace std;
int main(int argc, char *argv[])
{
    ifstream in(argv[1]), out(argv[2]);
    int t = 5;
    cout << t << endl;
    while (t--)
    {
        int a, b;
        a = rand(), b = rand();
        cout << a << " " << b << endl;
        cout.flush();
        int ans;
        cin >> ans;
        if (ans != a + b)
        {
            perror("wa");
            return 0;
        }
    }
    perror("ac");
    return 0;
}