import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { HamburgerMenu } from "../components/HamburgerMenu";
import "../styles/common.css";

export const SignupPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      alert("パスワードが一致しません");
      return;
    }
    // TODO: 新規登録処理の実装
    navigate("/login");
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
          <button type="submit" className="button">登録</button>
        </form>
      </div>
    </div>
  );
};
