#include <bits/stdc++.h>
using namespace std;
int main()
{
    int t = 5;
    cout << t << endl;
    cout.flush();
    while (t--)
    {
        int a, b;
        a = rand() % 100000, b = rand() % 100000;
        cout << a << " " << b << endl;
        cout.flush();
        int ans;
        cin >> ans;
        if (ans != (a + b))
        {
            fprintf(stderr, "wa");
            return 0;
        }
    }
    fprintf(stderr, "ac");
    return 0;
}