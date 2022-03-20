# process_utils

- OS Process の汎用ライブラリ

## メモ

- Process の Thread の確認方法
  - /proc/[pid]/task/.. に PID の一覧が表示される
- Process の子 Process の確認方法
  - /proc/[pid]/task/[task]/children ファイルに PID の一覧が書かれている
  - /proc/[pid]/status ファイルに PPid が書かれているので、全 Process を探索すれば必然的に子 Process がわかる
