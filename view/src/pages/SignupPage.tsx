import { useState } from "react";
import { authClient } from "../lib/grpcClient";

export const SignupPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [checkPassword, setCheckPassword] = useState("");
  const [message, setMessage] = useState("");

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await authClient.signup({ email, password, checkPassword });
      setMessage("サインアップ成功！");
    } catch (error) {
      setMessage("サインアップ失敗: " + (error.message || String(error)));
    }
  };

  return (
    <div>
      <h2>Signup</h2>
      <form onSubmit={handleSignup}>
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
        <div>
          <label>
            パスワードの確認:{" "}
            <input
              type="password"
              value={checkPassword}
              onChange={(e) => setCheckPassword(e.target.value)}
              required
            />
          </label>
        </div>
        <button type="submit">Signup</button>
      </form>
      {message && <div>{message}</div>}
    </div>
  );
};
