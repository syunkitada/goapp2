# process_utils

- OS Process の汎用ライブラリ

## メモ

- Process の Thread の確認方法
  - /proc/[pid]/task/.. に Thread の PID の一覧が表示される
  - /proc/[pid]/task/[tid] もしくは、/proc/[tid] で Thread のプロセス情報を得ることができる
    - /proc/[tid] は、ls /proc などとしても表示されないので注意（直接 ls /proc/[tid] とすると実態があることがわかる）
- Process の子 Process の確認方法
  - /proc/[pid]/task/[task]/children ファイルに PID の一覧が書かれている
    - ただし、Kernel のビルドオプションによっては無効化されているので注意
  - /proc/[pid]/status ファイルに PPid が書かれており、これが親プロセスの PID となる
    - 全 Process を探索すれば必然的に子 Process がわかる
- /proc/[pid]/io は、root 権限がないと見れない
  - プロセスと同一ユーザであれば見れる
