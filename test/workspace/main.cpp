#include <stdio.h>
#include <bits/stdc++.h>
using namespace std;
int func(int a, int b)
{
    return a + b - 1;
}
int main()
{
    int t;
    cin >> t;
    while (t--)
    {
        int a, b;
        cin >> a >> b;
        cout << func(a, b) << endl;
        cout.flush();
    }
    return 0;
}
