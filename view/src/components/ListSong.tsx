import { useEffect, useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { Song } from "../gen/master/song_pb";
import { getMusicVideoTypeDisplayName } from "../utils/enumDisplay";


export const SongList = () => {
  const [songs, setSongs] = useState<Song[]>([]);

  useEffect(() => {
    const fetchSongs = async () => {
      const response = await masterClient.getSongs({});
      setSongs(response.songs);
    };

    fetchSongs().catch(console.error);
  }, []);

  console.log(songs);

  return (
    <div>
      <h2>楽曲一覧</h2>
      <ul>
        {songs.map((song) => (
          <li key={song.id}>
            <p><strong>ID:</strong> {song.id}</p>
            <p><strong>Name:</strong> {song.name}</p>
            <p><strong>Kana:</strong> {song.kana}</p>
            <p><strong>Lyrics:</strong> {song.lyrics?.name}</p>
            <p><strong>Music:</strong> {song.music?.name}</p>
            <p><strong>Arrangement:</strong> {song.arrangement?.name}</p>
            <p><strong>Thumbnail:</strong> {song.thumbnail}</p>
            <p><strong>OriginalVideo:</strong> {song.originalVideo}</p>
            <p><strong>Deleted:</strong> {song.deleted ? "Yes" : "No"}</p>
            <p><strong>Release Time:</strong> {song.releaseTime?.toDate().toLocaleString()}</p>
            <p><strong>Vocal Patterns:</strong></p>
            {song.vocalPatterns && song.vocalPatterns.length > 0 && (
              <ul>
                {song.vocalPatterns.map((pattern, index) => (
                  <li key={index}>
                    <p><strong>Name:</strong> {pattern.name}</p>
                    <p><strong>Singers:</strong> {pattern.singers?.map(singer => singer.name).join(", ")}</p>
                    <p><strong>Units:</strong> {pattern.units?.map(unit => unit.name).join(", ")}</p>
                  </li>
                ))}
              </ul>
            )}
            <p><strong>Music Video Types:</strong></p>
            <ul>
                {song.musicVideoTypes.map((type) =>
                    getMusicVideoTypeDisplayName(type)
                ).join(", ")}
            </ul>
          </li>
        ))}
      </ul>
    </div>
  );
};
