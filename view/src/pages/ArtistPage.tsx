import { ArtistList } from "../components/ListArtist";
import { CreateArtist } from "../components/CreateArtist";
import { HamburgerMenu } from "../components/HamburgerMenu";
import { Link } from "react-router-dom";
import { useAdminGuard } from "../hooks/useAdminGuard";

export const ArtistPage = () => {
  const checked = useAdminGuard();
  if (!checked) return null;
  return (
    <div className="container">
      <HamburgerMenu />
      <h1>Sekai Songs Mylist - Master - Artist</h1>
      <ArtistList />
      <CreateArtist />
      <li>
        <Link to="/master">master</Link>
      </li>
    </div>
  );
};
