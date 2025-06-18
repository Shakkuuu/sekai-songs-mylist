import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { HamburgerMenu } from "../components/HamburgerMenu";
import "../styles/common.css";

export const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: ログイン処理の実装
    navigate("/mylist");
  };

  return (
    <div className="container">
      <HamburgerMenu />
      <div className="page-header">
        <h1>ログイン</h1>
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
          <button type="submit" className="button">ログイン</button>
        </form>
      </div>
    </div>
  );
};
