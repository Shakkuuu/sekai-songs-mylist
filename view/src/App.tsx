import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { ArtistPage } from "./pages/ArtistPage";
import { SingerPage } from "./pages/SingerPage";
import { UnitPage } from "./pages/UnitPage";
import { SongPage } from "./pages/SongPage";
import { ChartPage } from "./pages/ChartPage";
import { SignupPage } from "./pages/SignupPage";
import { LoginPage } from "./pages/LoginPage";
import { TopPage } from "./pages/TopPage";
import { UserPage } from "./pages/UserPage";
import { MyListPage } from "./pages/MyListPage";
import { MyListDetailPage } from "./pages/MyListDetailPage";
import { MyListEditPage } from "./pages/MyListEditPage";
import { MyListChartDetailPage } from "./pages/MyListChartDetailPage";
import { MasterPage } from "./pages/MasterPage";

export const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<TopPage />} />
        <Route path="/master" element={<MasterPage />} />
        <Route path="/master/artist" element={<ArtistPage />} />
        <Route path="/master/singer" element={<SingerPage />} />
        <Route path="/master/unit" element={<UnitPage />} />
        <Route path="/master/song" element={<SongPage />} />
        <Route path="/master/chart" element={<ChartPage />} />
        <Route path="/signup" element={<SignupPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/user" element={<UserPage />} />
        <Route path="/mylist" element={<MyListPage />} />
        <Route path="/mylist/:myListId" element={<MyListDetailPage />} />
        <Route path="/mylist/:myListId/edit" element={<MyListEditPage />} />
        <Route
          path="/mylist/:myListId/chart/:myListChartId"
          element={<MyListChartDetailPage />}
        />
        <Route path="*" element={<div>Not Found</div>} />
      </Routes>
    </Router>
  );
};

export default App;
