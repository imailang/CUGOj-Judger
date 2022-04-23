#include <stdio.h>
#include <bits/stdc++.h>
using namespace std;
int func(int a, int b)
{
    return a + b;
}
int num[3];
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
    int a = 1, b = 2;
    for (int i = 0; i < 1000000000; i++)
    {
        a = b ^ 3;
        b = a ^ rand() % 10;
    }
    return 0;
}
