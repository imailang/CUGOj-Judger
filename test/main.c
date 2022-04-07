#include<stdio.h>
int func(int a,int b)
{
    return a+b;
}
int main()
{
    int a,b;
    scanf("%d%d",&a,&b);
    printf("%d\n",func(a,b));
    return 0;
}
