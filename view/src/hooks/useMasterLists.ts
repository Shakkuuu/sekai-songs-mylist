import { useEffect, useState } from "react";
import { masterClient } from "../lib/grpcClient";
import { Artist } from "../gen/master/artist_pb";
import { Singer } from "../gen/master/singer_pb";
import { Unit } from "../gen/master/unit_pb";
import { Song } from "../gen/master/song_pb";

export function useArtists() {
  const [artists, setArtists] = useState<Artist[]>([]);
  useEffect(() => {
    masterClient.getArtists({}).then(res => setArtists(res.artists));
  }, []);
  return artists;
}

export function useSingers() {
  const [singers, setSingers] = useState<Singer[]>([]);
  useEffect(() => {
    masterClient.getSingers({}).then(res => setSingers(res.singers));
  }, []);
  return singers;
}

export function useUnits() {
  const [units, setUnits] = useState<Unit[]>([]);
  useEffect(() => {
    masterClient.getUnits({}).then(res => setUnits(res.units));
  }, []);
  return units;
}

export function useSongs() {
  const [songs, setSongs] = useState<Song[]>([]);
  useEffect(() => {
    masterClient.getSongs({}).then(res => setSongs(res.songs));
  }, []);
  return songs;
}
