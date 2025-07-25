import { useEffect, useState } from "react";
import { userClient } from "../lib/grpcClient";
import { useNavigate } from "react-router-dom";
import { HamburgerMenu } from "../components/HamburgerMenu";
import "../styles/common.css";
// import { ConnectError } from "@bufbuild/connect";

export const UserPage = () => {
  const [user, setUser] = useState<{
    id: string;
    email: string;
    createdAt: string;
  } | null>(null);
  const [email, setEmail] = useState("");
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [checkPassword, setCheckPassword] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  // ユーザー情報取得
  const fetchUserInfo = async () => {
    try {
      const res = await userClient.userInfo({});
      setUser({
        id: res.id,
        email: res.email,
        createdAt: res.createdAt?.toDate().toLocaleString() ?? "",
      });
      setEmail(res.email);
    } catch (error) {
      console.error(error);
      setMessage("ユーザー情報取得に失敗しました");
    }
  };

  useEffect(() => {
    fetchUserInfo();
  }, []);

  // メールアドレス変更
  const handleChangeEmail = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await userClient.changeEmail({ email });
      setMessage("メールアドレスを変更しました。ログアウトします。");
      localStorage.removeItem("token");
      setTimeout(() => navigate("/login"), 1000);
    } catch (error) {
      console.error(error);
      setMessage("メールアドレス変更に失敗しました");
    }
  };

  // パスワード変更
  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await userClient.changePassword({
        oldPassword,
        newPassword,
        newCheckPassword: checkPassword,
      });
      setMessage("パスワードを変更しました。ログアウトします。");
      localStorage.removeItem("token");
      setTimeout(() => navigate("/login"), 1000);
    } catch (error) {
      console.error(error);
      setMessage("パスワード変更に失敗しました");
    }
  };

  // 退会
  const handleDeleteUser = async () => {
    if (!window.confirm("本当に退会しますか？")) return;
    try {
      await userClient.deleteUser({});
      setMessage("退会しました。サインアップ画面に移動します。");
      localStorage.removeItem("token");
      setTimeout(() => navigate("/signup"), 1000);
    } catch (error) {
      console.error(error);
      setMessage("退会に失敗しました");
    }
  };

  // ログアウト
  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  return (
    <div className="container">
      <HamburgerMenu />
      <div className="page-header">
        <h1>ユーザー情報</h1>
      </div>
      <div className="card">
        <h2>プロフィール</h2>
        <div className="form-group">
          <label>メールアドレス</label>
          <p>{user ? user.email : "ユーザー情報取得中..."}</p>
        </div>
        <div className="form-group">
          <label>登録日</label>
          <p>{user ? user.createdAt : "ユーザー情報取得中..."}</p>
        </div>
      </div>
      <h3>メールアドレス変更</h3>
      <form onSubmit={handleChangeEmail}>
        <input
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <button type="submit">変更</button>
      </form>

      <h3>パスワード変更</h3>
      <form onSubmit={handleChangePassword}>
        <input
          type="password"
          placeholder="現在のパスワード"
          value={oldPassword}
          onChange={(e) => setOldPassword(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="新しいパスワード"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="新しいパスワード（確認）"
          value={checkPassword}
          onChange={(e) => setCheckPassword(e.target.value)}
          required
        />
        <button type="submit">変更</button>
      </form>

      <h3>アカウント操作</h3>
      <button onClick={handleLogout}>ログアウト</button>
      <button onClick={handleDeleteUser} style={{ color: "red" }}>
        退会
      </button>

      {message && <div>{message}</div>}
    </div>
  );
};
