import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { HamburgerMenu } from "../components/HamburgerMenu";
import { authClient } from "../lib/grpcClient";
import { LoginRequest } from "../gen/auth/v1/auth_pb";
import "../styles/common.css";

export const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      const req = new LoginRequest({ email, password });
      const res = await authClient.login(req);
      localStorage.setItem("token", res.token);
      navigate("/mylist");
    } catch (err) {
      if (err && typeof err === "object" && "code" in err) {
        const code = (err as { code: number }).code;
        if (code === 401) {
          setError("メールアドレスの認証が完了していません。メールをご確認ください。");
        } else if (code >= 500 && code < 600) {
          setError("処理に問題が発生しました。時間をおいてもう一度お試しください。");
        } else {
          setError("メールアドレスまたはパスワードが正しくありません。");
        }
      } else {
        setError("メールアドレスまたはパスワードが正しくありません。");
      }
    }
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
          {error && <div style={{ color: "red", marginTop: 8 }}>{error}</div>}
        </form>
      </div>
    </div>
  );
};
