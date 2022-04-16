#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
ll ans[100005];
int main()
{
    ll n, k, x;
    cin >> n >> k >> x;
    ll sum = 0;
    for (int i = 1; i <= n; i++)
    {
        ans[i] = i;
        sum += i;
        if (sum > x)
        {
            cout << "-1" << endl;
            return 0;
        }
    }
    ll add = (x - sum) / n;
    for (int i = 1; i <= n; i++)
    {
        ans[i] += add;
        sum += add;
    }
    for (int i = n; i >= 1; i--)
    {
        if (sum == x)
            break;
        ans[i]++;
        sum++;
    }
    if (ans[n] > k)
    {
        cout << "-1" << endl;
        return 0;
    }
    for (int i = 1; i <= n; i++)
    {
        printf("%lld ", ans[i]);
    }
    puts("");
    return 0;
}