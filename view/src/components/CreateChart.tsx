import { useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { ConnectError } from "@bufbuild/connect";
import { DifficultyType } from "../gen/enums/master_pb";
import { useSongs } from "../hooks/useMasterLists";

export const CreateChart = () => {
  const songs = useSongs();

  const [songId, setSongId] = useState<number>(0);
  const [difficultyType, setDifficultyType] = useState<DifficultyType>(
    DifficultyType.UNSPECIFIED
  );
  const [level, setLevel] = useState<number>(0);
  const [chartViewLink, setChartViewLink] = useState<string>("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      await masterClient.createChart({
        songId,
        difficultyType,
        level,
        chartViewLink,
      });

      alert("Chart created successfully!");
      setSongId(0);
      setDifficultyType(DifficultyType.UNSPECIFIED);
      setLevel(0);
      setChartViewLink("");
    } catch (error) {
      if (
        error instanceof ConnectError ||
        (error &&
          typeof error === "object" &&
          "name" in error &&
          error.name === "ConnectError")
      ) {
        // gRPC エラーの場合
        if ("code" in error) {
          console.error(
            "gRPC Error:",
            (error as ConnectError).code,
            (error as ConnectError).message
          );
        } else {
          console.error("gRPC Error:", (error as ConnectError).message);
        }
        const errorMessage =
          typeof error === "object" && error !== null && "message" in error
            ? String((error as { message: unknown }).message)
            : String(error);
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
      <h2>Create Chart</h2>
      <div>
        <label>
          Song:
          <select
            value={songId}
            onChange={(e) => setSongId(Number(e.target.value))}
            required
          >
            <option value={0}>選択してください</option>
            {songs.map((song) => (
              <option key={song.id} value={song.id}>
                {song.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div>
        <label>
          DifficultyType:
          <select
            value={difficultyType}
            onChange={(e) =>
              setDifficultyType(Number(e.target.value) as DifficultyType)
            }
          >
            <option value={DifficultyType.UNSPECIFIED}>Unspecified</option>
            <option value={DifficultyType.EASY}>Easy</option>
            <option value={DifficultyType.NORMAL}>Normal</option>
            <option value={DifficultyType.HARD}>Hard</option>
            <option value={DifficultyType.EXPERT}>Expert</option>
            <option value={DifficultyType.MASTER}>Master</option>
            <option value={DifficultyType.APPEND}>Append</option>
          </select>
        </label>
      </div>
      <div>
        <label>
          Level:
          <input
            type="number"
            value={level}
            onChange={(e) => setLevel(Number(e.target.value))}
          />
        </label>
      </div>
      <div>
        <label>
          ChartViewLink:
          <input
            type="string"
            value={chartViewLink}
            onChange={(e) => setChartViewLink(String(e.target.value))}
          />
        </label>
      </div>
      <button type="submit">Create</button>
    </form>
  );
};
