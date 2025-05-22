import { ArtistList } from "../components/ListArtist";
import { CreateArtist } from "../components/CreateArtist";
import { SingerList } from "../components/ListSinger";
import { CreateSinger } from "../components/CreateSinger";
import { UnitList } from "../components/ListUnit";
import { CreateUnit } from "../components/CreateUnit";
import { SongList } from "../components/ListSong";
import { CreateSong } from "../components/CreateSong";
import { CreateVocalPattern } from "../components/CreateVocalPattern";
import { ChartList } from "../components/ListChart";
import { CreateChart } from "../components/CreateChart";

export const MasterPage = () => (
  <div>
    <h1>Sekai Songs Mylist - Master</h1>
    <ArtistList />
    <CreateArtist />
    <SingerList />
    <CreateSinger />
    <UnitList />
    <CreateUnit />
    <CreateVocalPattern />
    <SongList />
    <CreateSong />
    <ChartList />
    <CreateChart />
  </div>
);
