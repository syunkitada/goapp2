# process_utils

## メモ

- プロセスの thread の確認方法
  - /proc/[pid]/task/.. に pid の一覧が表示される
- プロエスの子プロセスの確認方法
  - /proc/[pid]/task/[task]/children ファイルに pid の一覧が書かれている
  - /proc/[pid]/status ファイルに PPid が書かれているので、全プロセスを探索すれば必然的に子プロセスがわかる
