import { useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { Timestamp } from "@bufbuild/protobuf"; // Timestamp をインポート
import { ConnectError } from "@bufbuild/connect"; // ConnectError をインポート
import { MusicVideoType } from "../gen/enums/master_pb";
import { useArtists } from "../hooks/useMasterLists";

export const CreateSong = () => {
  const artists = useArtists();

  const [name, setName] = useState<string>("");
  const [kana, setKana] = useState<string>("");
  const [lyricsId, setLyricsId] = useState<number>(0);
  const [musicId, setMusicId] = useState<number>(0);
  const [arrangementId, setArrangementId] = useState<number>(0);
  const [thumbnail, setThumbnail] = useState<string>("");
  const [originalVideo, setOriginalVideo] = useState<string>("");
  const [releaseTime, setReleaseTime] = useState<string>(""); // ISO string
  const [deleted, setDeleted] = useState<boolean>(false);
  const [musicVideoTypes, setMusicVideoTypes] = useState<MusicVideoType[]>([]);

  const handleMusicVideoTypeChange = (type: MusicVideoType) => {
    setMusicVideoTypes((prev) =>
      prev.includes(type) ? prev.filter((t) => t !== type) : [...prev, type]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    console.log(Timestamp.fromDate(new Date(releaseTime)));
    try {
      await masterClient.createSong({
        name,
        kana,
        lyricsId,
        musicId,
        arrangementId,
        thumbnail,
        originalVideo,
        releaseTime: Timestamp.fromDate(new Date(releaseTime)),
        deleted,
        musicVideoTypes,
      });

      alert("Song created successfully!");
      setName("");
      setKana("");
      setLyricsId(0);
      setMusicId(0);
      setArrangementId(0);
      setThumbnail("");
      setOriginalVideo("");
      setReleaseTime("");
      setDeleted(false);
      setMusicVideoTypes([]);
    } catch (error) {
      if (error instanceof ConnectError) {
        // gRPC エラーの場合
        console.error("gRPC Error:", error.code, error.message);
        alert(`Error: ${error.message || String(error)}`);
      } else {
        // その他のエラーの場合
        console.error("Unexpected Error:", error);
        alert("An unexpected error occurred.");
      }
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>Create Song</h2>
      <div>
        <label>
          Name:
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </label>
      </div>
      <div>
        <label>
          Kana:
          <input
            type="text"
            value={kana}
            onChange={(e) => setKana(e.target.value)}
            required
          />
        </label>
      </div>
      <div>
        <label>
          Lyrics:
          <select
            value={lyricsId}
            onChange={(e) => setLyricsId(Number(e.target.value))}
            required
          >
            <option value={0}>選択してください</option>
            {artists.map((artist) => (
              <option key={artist.id} value={artist.id}>
                {artist.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div>
        <label>
          Music:
          <select
            value={musicId}
            onChange={(e) => setMusicId(Number(e.target.value))}
            required
          >
            <option value={0}>選択してください</option>
            {artists.map((artist) => (
              <option key={artist.id} value={artist.id}>
                {artist.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div>
        <label>
          Arrangement:
          <select
            value={arrangementId}
            onChange={(e) => setArrangementId(Number(e.target.value))}
            required
          >
            <option value={0}>選択してください</option>
            {artists.map((artist) => (
              <option key={artist.id} value={artist.id}>
                {artist.name}
              </option>
            ))}
          </select>
        </label>
      </div>
      <div>
        <label>
          Thumbnail:
          <input
            type="text"
            value={thumbnail}
            onChange={(e) => setThumbnail(e.target.value)}
          />
        </label>
      </div>
      <div>
        <label>
          Original Video:
          <input
            type="text"
            value={originalVideo}
            onChange={(e) => setOriginalVideo(e.target.value)}
          />
        </label>
      </div>
      <div>
        <label>
          Release Time (ISO format):
          <input
            type="datetime-local"
            value={releaseTime}
            onChange={(e) => setReleaseTime(e.target.value)}
          />
        </label>
      </div>
      <div>
        <label>
          Deleted:
          <input
            type="checkbox"
            checked={deleted}
            onChange={(e) => setDeleted(e.target.checked)}
          />
        </label>
      </div>
      <div>
        <label>
          MusicVideoTypes:
          <div>
            <label>
              <input
                type="checkbox"
                checked={musicVideoTypes.includes(
                  MusicVideoType.MUSIC_VIDEO_TYPE_3D
                )}
                onChange={() =>
                  handleMusicVideoTypeChange(MusicVideoType.MUSIC_VIDEO_TYPE_3D)
                }
              />
              3DMV
            </label>
            <label>
              <input
                type="checkbox"
                checked={musicVideoTypes.includes(
                  MusicVideoType.MUSIC_VIDEO_TYPE_2D
                )}
                onChange={() =>
                  handleMusicVideoTypeChange(MusicVideoType.MUSIC_VIDEO_TYPE_2D)
                }
              />
              2DMV
            </label>
            <label>
              <input
                type="checkbox"
                checked={musicVideoTypes.includes(
                  MusicVideoType.MUSIC_VIDEO_TYPE_ORIGINAL
                )}
                onChange={() =>
                  handleMusicVideoTypeChange(
                    MusicVideoType.MUSIC_VIDEO_TYPE_ORIGINAL
                  )
                }
              />
              原曲MV
            </label>
          </div>
        </label>
      </div>
      <button type="submit">Create</button>
    </form>
  );
};
