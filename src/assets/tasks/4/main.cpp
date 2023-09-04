#include "bits/stdc++.h"
using namespace std;
#define int long long
mt19937_64 rng((int) std::chrono::steady_clock::now().time_since_epoch().count());

int rnd(int x, int y) {
  int u = uniform_int_distribution<int>(x, y)(rng);
  return u;
}

int n, d;
int ask(int l, int r) {
  cout <<"? " << l << " " << r << endl;
  int x;
  cin >>x ;
  return x;
}

int eval(int x) {
  int ans = 0;
  for(int i=max(1ll, x-d); i<=min(n, x+d); i++) {
    for(int j=0; j<d; j++) { // position i to position x, brightness j
      ans += max(0ll, j - abs(i - x));
    }
    for(int j=d; j>0; j--) {
      ans += max(0ll, j - abs(i - x));
    }
  }
  return ans;
}
void solve(int tc) {
  
  cin >> n >> d;
  int mid = (n + 1) / 2 - d;
  int re = 0;
  for(int i=0; i<2*d; i++) re += ask(1, mid);
  int ee = 0;
  for(int i=1; i<=mid; i++) ee += eval(i);
  int lb, rb, mode;
  if(re == ee) lb = mid + 1, rb = n, mode = 1;
  else lb = 1, rb = mid + 2*d, mode = 2;
  if(mode == 2) {
    while(lb < rb) {
      int mid = (lb + rb + 1) >> 1;
      int re = 0;
      for(int i=0; i<2*d; i++) re += ask(mid, n);
      int ee = 0;
      for(int i=mid; i<=n; i++) {
        ee += eval(i);
      }
      if(re == ee) rb = mid - 1;
      else lb = mid;
    }
    cout << "! " << lb - d + 1 << endl;
  }
  else {
    while(lb < rb) {
      int mid = (lb + rb) >> 1;
      int re = 0;
      for(int i=0; i<2*d; i++) re += ask(1, mid);
      int ee = 0;
      for(int i=1; i<=mid; i++) {
        ee += eval(i);
      }
      if(re == ee) lb = mid + 1;
      else rb = mid;
    }
    cout << "! " << lb + d - 1 << endl;
  }
}
int32_t main() {
 // ios::sync_with_stdio(0);
 // cin.tie(0);;
  int t = 1;
  // cin >> t;
  for(int i=1; i<=t; i++) {
    solve(i);
  }
}