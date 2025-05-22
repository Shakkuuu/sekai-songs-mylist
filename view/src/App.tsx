import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import { MasterPage } from "./pages/MasterPage";
import { SignupPage } from "./pages/SignupPage";
import { LoginPage } from "./pages/LoginPage";
import { TopPage } from "./pages/TopPage";
import { UserPage } from "./pages/UserPage";

export const App = () => {
  return (
    <Router>
      <nav>
        <ul>
          <li>
            <Link to="/master">master</Link>
          </li>
          <li>
            <Link to="/signup">signup</Link>
          </li>
          <li>
            <Link to="/login">login</Link>
          </li>
          <li>
            <Link to="/user">user</Link>
          </li>
          <li>
            <Link to="/">top</Link>
          </li>
        </ul>
      </nav>
      <Routes>
        <Route path="/" element={<TopPage />} />
        <Route path="/master" element={<MasterPage />} />
        <Route path="/signup" element={<SignupPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/user" element={<UserPage />} />
        {/* <Route path="*" element={<div>Not Found</div>} /> */}
      </Routes>
    </Router>
  );
};

export default App;
