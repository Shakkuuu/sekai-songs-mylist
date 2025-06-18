import { HamburgerMenu } from "../components/HamburgerMenu";
import "../styles/common.css";

export const TopPage = () => {
  return (
    <div className="container">
      <HamburgerMenu />
      <div className="page-header">
        <h1>SEKAI Songs MyList</h1>
      </div>
      <div className="card">
        <h2>ようこそ</h2>
        <p>SEKAI Songs MyListは、プロジェクトセカイの楽曲を管理できるサービスです。</p>
        <p>以下の機能をご利用いただけます：</p>
        <ul>
          <li>楽曲のマイリスト作成</li>
          <li>クリア状況の管理</li>
          <li>譜面の詳細情報の閲覧</li>
          <li>添付ファイルの管理</li>
        </ul>
      </div>
      <div className="card">
        <h2>はじめ方</h2>
        <p>1. 新規登録またはログインしてください</p>
        <p>2. マイリストを作成して楽曲を追加</p>
        <p>3. 楽曲のクリア状況やメモを管理</p>
      </div>
    </div>
  );
};
