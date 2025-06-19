import { useEffect, useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { Song } from "../gen/master/song_pb";
import "./ListSong.css";
import { IMAGE_BASE_URL } from "../lib/constants";

interface SongDetailModalProps {
  song: Song;
  onClose: () => void;
}

const SongDetailModal = ({ song, onClose }: SongDetailModalProps) => {
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onClose}>
          ×
        </button>
        <h3>{song.name}</h3>
        <div className="modal-details">
          <p>
            <strong>原曲MV:</strong>{" "}
            <a
              href={song.originalVideo}
              target="_blank"
              rel="noopener noreferrer"
            >
              {song.originalVideo}
            </a>
          </p>
          <p>
            <strong>リリース時間:</strong>{" "}
            {song.releaseTime?.toDate().toLocaleString()}
          </p>
          <p>
            <strong>ボーカルパターン:</strong>
          </p>
          <ul>
            {song.vocalPatterns?.map((pattern, index) => (
              <li key={index}>
                {pattern.name} -{" "}
                {pattern.singers?.map((singer) => singer.name).join(", ")}
              </li>
            ))}
          </ul>
          <p>
            <strong>ユニット:</strong>{" "}
            {song.units?.map((unit) => unit.name).join(", ")}
          </p>
        </div>
      </div>
    </div>
  );
};

export const SongList = () => {
  const [songs, setSongs] = useState<Song[]>([]);
  const [selectedSong, setSelectedSong] = useState<Song | null>(null);

  useEffect(() => {
    const fetchSongs = async () => {
      const response = await masterClient.getSongs({});
      setSongs(response.songs);
    };

    fetchSongs().catch(console.error);
  }, []);

  return (
    <div className="song-list-container">
      <h2>楽曲一覧</h2>
      <div className="song-list">
        {songs.map((song) => (
          <div key={song.id} className="song-item">
            <div className="song-thumbnail">
              {song.thumbnail ? (
                <img
                  src={
                    song.thumbnail.startsWith("http")
                      ? song.thumbnail
                      : `${IMAGE_BASE_URL}${song.thumbnail}`
                  }
                  alt={song.name}
                />
              ) : (
                <div className="no-thumbnail">No Image</div>
              )}
            </div>
            <div className="song-info">
              <h3 className="song-name">{song.name}</h3>
              <p className="song-creators">
                {song.lyrics?.name} / {song.music?.name} /{" "}
                {song.arrangement?.name}
              </p>
            </div>
            <button
              className="detail-button"
              onClick={() => setSelectedSong(song)}
            >
              詳細
            </button>
          </div>
        ))}
      </div>
      {selectedSong && (
        <SongDetailModal
          song={selectedSong}
          onClose={() => setSelectedSong(null)}
        />
      )}
    </div>
  );
};
