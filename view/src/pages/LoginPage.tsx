import { useState } from "react";
import { authClient } from "../lib/grpcClient";

export const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await authClient.login({ email, password });
      if (res.token) {
        localStorage.setItem("token", res.token);
        setMessage("ログイン成功！");
      } else {
        setMessage("トークンが取得できませんでした");
      }
    } catch (error) {
      setMessage("ログイン失敗: " + (error.message || String(error)));
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
