import { ChartList } from "../components/ListChart";
import { CreateChart } from "../components/CreateChart";
import { HamburgerMenu } from "../components/HamburgerMenu";
import { Link } from "react-router-dom";
import { useAdminGuard } from "../hooks/useAdminGuard";

export const ChartPage = () => {
  const checked = useAdminGuard();
  if (!checked) return null;
  return (
    <div className="container">
      <HamburgerMenu />
      <h1>Sekai Songs Mylist - Master - Chart</h1>
      <ChartList />
      <CreateChart />
      <li>
        <Link to="/master">master</Link>
      </li>
    </div>
  );
};
