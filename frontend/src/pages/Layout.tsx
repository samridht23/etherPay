import Navbar from "../components/Navbar";
import { Outlet, useLocation } from 'react-router-dom';

const PageLayout: React.FC = () => {
  const location = useLocation();
  const path = location.pathname;

  if (path == '/') {
    return <Outlet />;
  } else {
    return (
      <>
        <Navbar />
        <div className="pt-[55px]">
          <Outlet />
        </div>
      </>
    );
  }
};
export default PageLayout;
