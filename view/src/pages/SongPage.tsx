import { SongList } from "../components/ListSong";
import { CreateSong } from "../components/CreateSong";
import { CreateVocalPattern } from "../components/CreateVocalPattern";
import { HamburgerMenu } from "../components/HamburgerMenu";
import { Link } from "react-router-dom";
import { useAdminGuard } from "../hooks/useAdminGuard";

export const SongPage = () => {
  const checked = useAdminGuard();
  if (!checked) return null;
  return (
    <div className="container">
      <HamburgerMenu />
      <h1>Sekai Songs Mylist - Master - Song</h1>
      <CreateVocalPattern />
      <SongList />
      <CreateSong />
      <li>
        <Link to="/master">master</Link>
      </li>
    </div>
  );
};
