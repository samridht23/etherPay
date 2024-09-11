import { Outlet, NavLink, useParams } from "react-router-dom";
import clsx from "clsx";
import { zustStore } from "../../App";
import { useEffect, useContext } from "react";
import { AuthContext } from '../../AuthWrapper'
import { ViewingProfileData } from "../../App"



const Profile = () => {
  const { isAuthenticated, isLoading } = useContext(AuthContext);
  const { address, viewingProfileData, updateViewingProfileData } = zustStore();
  let { user_address } = useParams();
  const fetchProfileData = async () => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_URL}/profile/${user_address}`,
        {
          method: 'GET',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );
      if (response.ok) {
        const data = await response.json();
        const profileData = data.data as ViewingProfileData;
        updateViewingProfileData(profileData);
        console.log(profileData)
      } else {
        console.error('Failed to fetch user data');
      }
    } catch (err) {
      console.error(err);
    }
  };
  useEffect(() => {
    if (!isLoading) {
      fetchProfileData();
    }
  }, [isAuthenticated, isLoading, address, user_address]);

  return (
    <div className="w-full">
      <div className="border-b-[1px] border-[var(--border-muted)]">
        <div className="max-w-5xl m-auto">
          <div>
            <div className={clsx(
              "w-full bg-[var(--bg-subtle)] h-48 rounded-b-[var(--radius-lg)]",
              "overflow-hidden items-center justify-between"
            )}>
              <img
                src="https://images.unsplash.com/photo-1725610588145-a508e5cfe90b?q=80&w=3500&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
                className="w-full h-full object-cover"
              />
            </div>
            <div className="flex flex-col items-center justify-center gap-2 -translate-y-12">
              <div className="w-28 h-28 rounded-full overflow-hidden border-[5px] border-[var(--bg)]">
                <img src="https://i.pravatar.cc/300" />
              </div>
              <div className="flex flex-col gap-1 w-full items-center">
                <div className="font-semibold text-base">
                  {viewingProfileData?.username}
                </div>
                <div className="bg-[var(--bg-subtle)] px-4 py-1 rounded-full">
                  {viewingProfileData?.address}
                </div>
              </div>
            </div>
          </div>
          <div className="flex items-center justify-center gap-10 font-semibold">
            <NavLink
              to="" end
              className={({ isActive }) =>
                clsx(
                  "py-2.5 px-6 leading-normal cursor-pointer text-content flex items-center text-center justify-center",
                  "outline-none border-b-[3px]",
                  isActive ? "border-[var(--border-inverted)]" : "border-transparent",
                  isActive ? "text-[var(--text-emphasis)]" : "text-[var(--text)]",
                  "hover:border-white"
                )
              }
            >
              About
            </NavLink>
            <NavLink
              to={`feed`}
              className={({ isActive }) =>
                clsx(
                  "py-2.5 px-6 leading-normal cursor-pointer text-content flex items-center text-center justify-center",
                  "outline-none border-b-[3px]",
                  isActive ? "border-[var(--border-inverted)]" : "border-transparent",
                  isActive ? "text-[var(--text-emphasis)]" : "text-[var(--text)]",
                  "hover:border-white"
                )
              }
            >
              Feed
            </NavLink>
            <NavLink
              to={`stream`}
              className={({ isActive }) =>
                clsx(
                  "py-2.5 px-6 leading-normal cursor-pointer text-content flex items-center text-center justify-center",
                  "outline-none border-b-[3px]",
                  isActive ? "border-[var(--border-inverted)]" : "border-transparent",
                  isActive ? "text-[var(--text-emphasis)]" : "text-[var(--text)]",
                  "hover:border-white"
                )
              }
            >
              Stream
            </NavLink>
            <NavLink
              to={`settings`}
              className={({ isActive }) =>
                clsx(
                  "py-2.5 px-6 leading-normal cursor-pointer text-content flex items-center text-center justify-center",
                  "outline-none border-b-[3px]",
                  isActive ? "border-[var(--border-inverted)]" : "border-transparent",
                  isActive ? "text-[var(--text-emphasis)]" : "text-[var(--text)]",
                  "hover:border-white"
                )
              }
            >
              Settings
            </NavLink>
          </div>
        </div>
      </div>
      <div className="w-full">
        <div className="max-w-5xl m-auto mt-6">
          <Outlet />
        </div>
      </div>
    </div>
  )
};

export default Profile;
