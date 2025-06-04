import { useState } from "react";
import { authClient } from "../lib/grpcClient";
import { useNavigate } from "react-router-dom";

export const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await authClient.login({ email, password });
      if (res.token) {
        localStorage.setItem("token", res.token);
        setMessage("ログイン成功！ユーザー画面に移行します。");
        setTimeout(() => navigate("/user"), 1000);
      } else {
        setMessage("トークンが取得できませんでした");
      }
    } catch (error) {
      setMessage(
        "ログイン失敗: " +
          (error instanceof Error ? error.message : String(error))
      );
    }
  };

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleLogin}>
        <div>
          <label>
            Email:{" "}
            <input
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </label>
        </div>
        <div>
          <label>
            Password:{" "}
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </label>
        </div>
        <button type="submit">Login</button>
      </form>
      {message && <div>{message}</div>}
    </div>
  );
};
