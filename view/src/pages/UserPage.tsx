import { useEffect, useState } from "react";
import { userClient } from "../lib/grpcClient";
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
      setMessage("メールアドレスを変更しました");
      fetchUserInfo();
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
      setMessage("パスワードを変更しました");
      setOldPassword("");
      setNewPassword("");
      setCheckPassword("");
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
      setMessage("退会しました");
      // ログアウト処理やリダイレクトなど
    } catch (error) {
      console.error(error);
      setMessage("退会に失敗しました");
    }
  };

  // ログアウト
  //   const handleLogout = async () => {
  //     try {
  //       await userClient.logout({});
  //       setMessage("ログアウトしました");
  //       // トークン削除やリダイレクトなど
  //     } catch (error) {
  //       setMessage("ログアウトに失敗しました");
  //     }
  //   };

  return (
    <div>
      <h2>ユーザー情報</h2>
      {user ? (
        <div>
          <div>ID: {user.id}</div>
          <div>メール: {user.email}</div>
          <div>作成日: {user.createdAt}</div>
        </div>
      ) : (
        <div>ユーザー情報取得中...</div>
      )}

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
      {/* <button onClick={handleLogout}>ログアウト</button> */}
      <button onClick={handleDeleteUser} style={{ color: "red" }}>
        退会
      </button>

      {message && <div>{message}</div>}
    </div>
  );
};
