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
$ ./server --token=[YOUR_VERIFICATION_TOKEN] --due=[YOUR_TEAM'S_ LENDING_PERIOD]
```

:heavy\_exclamation\_mark: due(貸出期間)のデフォルト値は14(日間)です。

### 詳しく
#### 全体の流れ

1. 備品一覧を確認
2. 借りたい備品が備品一覧にない場合，**/equipadd**で追加
2. **/equipborrow**で貸出
3. **/equipreturn**で返却

<br>

##### 備品一覧ページ

[備品一覧](http://ik1-323-21780.vs.sakura.ne.jp:3000/equip)から貸出したい物品があるか確認

```
1 | Electronではじめるアプリ開発 | 1 | computer_club | 2019-02-20 | hsm_hx | 1 | 
```
- 1 : 備品ID
- Electronではじめるアプリ開発 : 備品名
- 1 : 種別
- computer\_club : 持ち主
- 2019-02-20 : 返却期日(になる予定)
- hsm\_hx : 貸出者
- 1 : 状態

<br>

**種別について**
1 : 本
2 : コンピュータ
3 : 周辺機器
4 : ケーブル類
5 : その他

<br>

**状態について**
0 : 部室にあります
1 : 借りられています

<br>

##### コマンドの使い方

**借りるとき**

Slackの中で以下のように入力して投稿します。(チャンネルはどこでもOK)
```
/equipborrow [備品ID]

======== 例：備品IDが1の備品を借りたい場合 ==========
/equipborrow 1
```
<br><br>

**返すとき**

Slackの中で以下のように入力して投稿します。(チャンネルはどこでもOK)
```
/equipreturn [備品ID]

======== 例：備品IDが2の備品を借りたい場合 ==========
/equipreturn 2
```
<br><br>

**借りたい備品が一覧にないとき**

Slackの中で以下のように入力して投稿します。(チャンネルはどこでもOK)
```
/equipadd [備品の名前] [種別]

======== 例：「Unityの本」という本を追加したい場合 ==========
/equipadd Unityの本 BOOK
```
**種別について**
- 本→BOOK
- コンピュータ(パソコン，マイコンなど)→COMPUTER
- 周辺機器(ヘッドホン，ゲームパッドなど)→SUPPLY
- ケーブル類→CABLE
- その他→OTHER

**持ち主を指定したい場合**
`/equipadd [備品の名前] [種別] [持ち主]`と入力して投稿します。(持ち主を指定しない場合computer_clubになります。

**注意事項**
本の名前にスペースを使わない。各要素は半角スペースで区切ってください。

<br>

## 今後できるようになること

- 貸出中のユーザーに返却を催促する
- 備品の削除は権限を持ったユーザーのみが行えるようになる

## このプロジェクトについて

:ok: みなさんからの不具合報告や追加機能のご提案をお待ちしています<br>
:ok: ぜひ使ってみてください<br>
:ok: ForkやPull Requestが来ると嬉しい<br>
