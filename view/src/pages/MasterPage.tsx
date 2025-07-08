import { HamburgerMenu } from "../components/HamburgerMenu";
import { Link } from "react-router-dom";
import { useAdminGuard } from "../hooks/useAdminGuard";

export const MasterPage = () => {
  const checked = useAdminGuard();
  if (!checked) return null;
  return (
    <div className="container">
      <HamburgerMenu />
      <h1>Sekai Songs Mylist - Master</h1>
      <ul>
        <li>
          <Link to="/master/artist">artist</Link>
        </li>
        <li>
          <Link to="/master/singer">singer</Link>
        </li>
        <li>
          <Link to="/master/unit">unit</Link>
        </li>
        <li>
          <Link to="/master/song">song</Link>
        </li>
        <li>
          <Link to="/master/chart">chart</Link>
        </li>
      </ul>
    </div>
  );
};
