import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { HamburgerMenu } from "../components/HamburgerMenu";
import { authClient } from "../lib/grpcClient";
import { SignupRequest } from "../gen/auth/v1/auth_pb";
import "../styles/common.css";

export const SignupPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (password !== confirmPassword) {
      setError("パスワードが一致しません");
      return;
    }
    try {
      const req = new SignupRequest({
        email,
        password,
        checkPassword: confirmPassword,
      });
      await authClient.signup(req);
      alert("登録が完了しました。メールを確認してください。");
      navigate("/login");
    } catch {
      setError(
        "新規登録に失敗しました。既に登録済みのメールアドレスあるいは入力項目に問題があります。"
      );
    }
  };

  return (
    <div className="container">
      <HamburgerMenu />
      <div className="page-header">
        <h1>新規登録</h1>
      </div>
      <div className="card">
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="email">メールアドレス</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="password">パスワード</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="confirmPassword">パスワード（確認）</label>
            <input
              type="password"
              id="confirmPassword"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
            />
          </div>
          <button type="submit" className="button">
            登録
          </button>
          {error && <div style={{ color: "red", marginTop: 8 }}>{error}</div>}
        </form>
      </div>
    </div>
  );
};
