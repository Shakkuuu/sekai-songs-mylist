import { useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { ConnectError } from "@bufbuild/connect";

export const CreateArtist = () => {
  const [name, setName] = useState("");
  const [kana, setKana] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await masterClient.createArtist({ name, kana });
      alert("アーティストを作成しました！");
      setName("");
      setKana("");
      window.location.reload();
    } catch (error) {
      if (
        error instanceof ConnectError ||
        (error &&
          typeof error === "object" &&
          "name" in error &&
          error.name === "ConnectError")
      ) {
        // gRPC エラーの場合
        const code =
          "code" in error ? (error as { code: string }).code : undefined;
        const errorMessage =
          typeof error === "object" && error !== null && "message" in error
            ? (error as { message?: string }).message
            : String(error);
        console.error("gRPC Error:", code, errorMessage);
        alert(`Error: ${errorMessage}`);
      } else {
        // その他のエラーの場合
        console.error("Unexpected Error:", error);
        alert("An unexpected error occurred.");
      }
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>アーティスト作成</h2>
      <div>
        <label>
          名前:
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </label>
      </div>
      <div>
        <label>
          かな:
          <input
            value={kana}
            onChange={(e) => setKana(e.target.value)}
            required
          />
        </label>
      </div>
      <button type="submit">作成</button>
    </form>
  );
};
