# sekai-songs-mylist

- Chartの実装と検証 ok
- 存在しないartsitIDとかsingerIDとか選ばれたときのinvalid argエラー返し ok
- badrequest ok
- 登録したらリロードされたい
- マスター管理画面整備、選択式とか ok
- enumは選択式にしたが、singerやartistもlistを取得してそれを選択肢にしたい ok
- singerのpositionによるsortがうまくいかない ok
- 検索機能
  - フロントの検索だけで良さそう？
- tokenでやりとりしているから、logoutはフロントだけでいいかも ok
- changeEmail,changePasswordした後はログアウトさせたいね ok
- mylistのフォルダ並び替え機能（フォルダの表示順なので中身は並び替えできなくていい。） ok
- mylistの一覧を取るときは、画面にはマイリスト名だけ並んで表示されてればいい ok
- マイリスト名をクリックしたらそのマイリストのmylistchartsを取ってきて並べて表示（表示には楽曲名、artist、歌唱、chart情報、クリア状況を表示） ok
- mylistchartsをクリックしたら、メモの表示やattachmentの取得&表示 ok
- mylistchartsをDBからとるタイミングでchartをjoinして取得するか、mylistchartsをlist取得したものをコードでforで回して、masterのキャッシュからとってchart情報を埋めていくほうどっちがいいか。
- キャッシュからとったほうがいいらしい ok
- ExistsMyListChartByMyListIDAndChartIDで、mylistchart追加時に重複しないか確認 ok

- 生成したコードをそのままリモートに上げてるけど、ちゃんとどのローカルでも生成できるようにした方が良さそうだよナー
- テスト書く

- export PATH="$(pwd)/../../view/node_modules/.bin:$PATH"
- buf generate
- sqlc generate
- 既存のデータ抽出
  - pg_dump -U db_user_name -h localhost -p 5432 -d dbname -Fp --data-only > dbdb.sql
