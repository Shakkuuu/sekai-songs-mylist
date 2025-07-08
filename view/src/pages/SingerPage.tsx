import { SingerList } from "../components/ListSinger";
import { CreateSinger } from "../components/CreateSinger";
import { HamburgerMenu } from "../components/HamburgerMenu";
import { Link } from "react-router-dom";
import { useAdminGuard } from "../hooks/useAdminGuard";

export const SingerPage = () => {
  const checked = useAdminGuard();
  if (!checked) return null;
  return (
    <div className="container">
      <HamburgerMenu />
      <h1>Sekai Songs Mylist - Master - Singer</h1>
      <SingerList />
      <CreateSinger />
      <li>
        <Link to="/master">master</Link>
      </li>
    </div>
  );
};
