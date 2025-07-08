import { UnitList } from "../components/ListUnit";
import { CreateUnit } from "../components/CreateUnit";
import { HamburgerMenu } from "../components/HamburgerMenu";
import { Link } from "react-router-dom";
import { useAdminGuard } from "../hooks/useAdminGuard";

export const UnitPage = () => {
  const checked = useAdminGuard();
  if (!checked) return null;
  return (
    <div className="container">
      <HamburgerMenu />
      <h1>Sekai Songs Mylist - Master - Unit</h1>
      <UnitList />
      <CreateUnit />
      <li>
        <Link to="/master">master</Link>
      </li>
    </div>
  );
};
