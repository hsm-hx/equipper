# equipper - 部活動のためのSlackコマンド・Webアプリ

## What is this?

部活動のための備品管理用Slackコマンドと備品管理状況がひと目でわかる一覧サイト。

### できること

- 備品を登録・貸出・返却・削除(Slackから)
- 備品の貸出状況をどこからでも確認(Webサイト)

### なにがうれしいのか

- ログインなどの煩雑な処理なく手軽に備品管理
- 普段使うツールによって貸出状況が通知される
- 学校からでも家からでも備品の貸出状況を確認

## 使い方

``` shell
$ git clone git@github.com:hsm-hx/equipper.git
$ cd equipper
$ dep ensure
$ go build ./src
$ ./server
$ ./server --token=[YOUR_VERIFICATION_TOKEN]
```

詳しくは[ブログ](https://mwc922-hsm.hatenablog.com/entry/2019/02/18/171421)にて

## 今後できるようになること

- 貸出中のユーザーに返却を催促する
- 備品の削除は権限を持ったユーザーのみが行えるようになる

## このプロジェクトについて

:ok: みなさんからの不具合報告や追加機能のご提案をお待ちしています<br>
:ok: ぜひ使ってみてください<br>
:ok: ForkやPull Requestが来ると嬉しい<br>
