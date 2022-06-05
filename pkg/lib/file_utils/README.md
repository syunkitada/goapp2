# file_utiles

## YAML のライブラリについてのメモ

- https://github.com/go-yaml/yaml(gopkg.in/yaml)
  - おそらくこれがデファクトに一番近い
  - 特に理由がないならこれを使うのがよい
- https://github.com/ghodss/yaml
  - go-yaml/yaml のラッパー
  - 読み込む場合は yaml を json に変換してから、書き込むときは json に変換してから yaml にする
    - 余計な処理は発生してしまうが、json.Marshal、json.Unmarshal がそのまま流用できる
    - また、ファイルが yaml でも json でも読み込むことができる
  - 便利なのだが、あまり積極的にメンテナンスはされてないので使うのはやめたほうが良さそう (at 2022/06/05)
