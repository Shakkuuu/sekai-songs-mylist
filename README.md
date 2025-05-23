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
- mylistのフォルダ並び替え機能（フォルダの表示順なので中身は並び替えできなくていい。）

- 生成したコードをそのままリモートに上げてるけど、ちゃんとどのローカルでも生成できるようにした方が良さそうだよナー

- export PATH="$(pwd)/../../view/node_modules/.bin:$PATH"
- buf generate
- sqlc generate
- 既存のデータ抽出
  - pg_dump -U db_user_name -h localhost -p 5432 -d dbname -Fp --data-only > dbdb.sql

多分実装ok 動作確認したい。
