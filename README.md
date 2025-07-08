# sekai-songs-mylist

## TODO

- 検索機能
  - フロントの検索だけで良さそう？
- 生成したコードをそのままリモートに上げてるけど、ちゃんとどのローカルでも生成できるようにした方が良さそうだよナー
- テスト書く
- フロント綺麗にする
- masterデータ作成
  - masterデータはv999.999.999_n_hoge-data.up.sqlで管理
  - dumpではmasterのデータはとりたくない
  - thumbnailの画像やmy_list_chart_attachmentsのfile_urlの画像とか動画は別のローカルサーバーを立てておく ok
  - そこにアップロードしてそこへのURLをDBに保存 ok
  - CreateSongの時にフロントがそこにPostでファイルをアップロードして、そのレスポンスのURLをバックエンドにTEXTで送ってDBに保存 ok
  - アタッチメントのfile_urlも同様 ok
  - thumbnailを表示するところは文字列じゃなくて画像を表示するように修正 ok
- ファイルの保存先としてGoogleDrive ok
  - 保存時にファイルサイズを軽くする
- バックエンドはRender
  - RESTでhealthCheck用ハンドラ作っておく
- フロントはどこに置こうか

## ok

- Chartの実装と検証 ok
- 存在しないartsitIDとかsingerIDとか選ばれたときのinvalid argエラー返し ok
- badrequest ok
- 登録したらリロードされたい ok
- マスター管理画面整備、選択式とか ok
- enumは選択式にしたが、singerやartistもlistを取得してそれを選択肢にしたい ok
- singerのpositionによるsortがうまくいかない ok
- tokenでやりとりしているから、logoutはフロントだけでいいかも ok
- changeEmail,changePasswordした後はログアウトさせたいね ok
- mylistのフォルダ並び替え機能（フォルダの表示順なので中身は並び替えできなくていい。） ok
- mylistの一覧を取るときは、画面にはマイリスト名だけ並んで表示されてればいい ok
- マイリスト名をクリックしたらそのマイリストのmylistchartsを取ってきて並べて表示（表示には楽曲名、artist、歌唱、chart情報、クリア状況を表示） ok
- mylistchartsをクリックしたら、メモの表示やattachmentの取得&表示 ok
- mylistchartsをDBからとるタイミングでchartをjoinして取得するか、mylistchartsをlist取得したものをコードでforで回して、masterのキャッシュからとってchart情報を埋めていくほうどっちがいいか。
- キャッシュからとったほうがいいらしい ok
- ExistsMyListChartByMyListIDAndChartIDで、mylistchart追加時に重複しないか確認 ok
- ユーザー作成時にメールでメールアドレス認証 ok
- handlerに認証メール再送のエンドポイント作成。クエリパラメータにメールアドレス。認証tokenと認証トークン期限を更新する必要あり。再送URLはメールの下の方に書いておく ok
- go fmt、フロントのURLを定数化 ok
- Masterページにアクセスできる管理者権限。どうやってその人だけ可能にするか。UserにisAdminを追加して、Adminのメアドリストをどこかで読み込ませるとか ok
  - 一旦sql直叩き ok
- masterのアクセス権。フロントのページ表示制限はどうするか。ページ遷移前にisAdminを確認して弾くとか？ ok

## メモ

- export PATH="$(pwd)/../../view/node_modules/.bin:$PATH"
- buf generate
- sqlc generate
- 既存のデータ抽出
  - pg_dump -U db_user_name -h localhost -p 5432 -d dbname -Fp --data-only > dbdb.sql
