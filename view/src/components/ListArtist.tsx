import { useEffect, useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { Artist } from "../gen/master/artist_pb"; // 型情報（optional）

export const ArtistList = () => {
  const [artists, setArtists] = useState<Artist[]>([]);

  useEffect(() => {
    const fetchArtists = async () => {
      const response = await masterClient.getArtists({});
      setArtists(response.artists);
    };

    fetchArtists().catch(console.error);
  }, []);

  return (
    <div>
      <h2>アーティスト一覧</h2>
      <ul>
        {artists.map((artist) => (
          <li key={artist.id}>{artist.id} : {artist.name}</li>
        ))}
      </ul>
    </div>
  );
};
