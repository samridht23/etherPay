import { Route, Routes } from "react-router-dom";
import PageLayout from "./Layout";
import Home from "./Home";
import Profile from "./Profile";
import Feed from "../components/Feed";
import Stream from "../components/Stream";
import About from "../components/About";
import Settings from "../components/Settings";


const Router: React.FC = () => {
  return (
    <>
      <Routes>
        <Route element={<PageLayout />}>
          <Route index element={<Home />} />
          <Route path=":user_address" element={<Profile />}>
            <Route index element={<About />} />
            <Route path="feed" element={<Feed />} />
            <Route path="stream" element={<Stream />} />
            <Route path="settings" element={<Settings />} />
          </Route>
        </Route>
      </Routes>
    </>
  );
};
export default Router;
