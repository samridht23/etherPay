import { MetaMaskInpageProvider } from '@metamask/providers';
import { useEffect, useContext } from 'react';
import clsx from 'clsx';
import { SettingsIcon, UserRoundIcon } from 'lucide-react';
import { Link } from "react-router-dom"
import { zustStore } from "../../App"
import { AuthContext } from '../../AuthWrapper';

import {
  DropdownMenuRoot,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuPortal,
} from "../../ui";

const Logo = () => {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="none" viewBox="0 0 20 20">
      <path fill="#fff" fill-rule="evenodd" d="m7.027 12.414 2.86 5.801 1.47-4.339-4.329-1.462Zm-5.23-2.518 4.331 1.47L7.6 7.03 1.797 9.896Zm5.483 1.84c.068.037.096.058.128.07 1.397.474 2.791.952 4.191 1.415.114.036.275.005.387-.05 1.601-.78 3.198-1.57 4.794-2.358.454-.225.907-.452 1.403-.7-.11-.041-.166-.065-.225-.085-1.352-.459-2.703-.92-4.06-1.367a.635.635 0 0 0-.434.027c-1.701.827-3.397 1.668-5.093 2.506-.354.174-.706.35-1.091.542Zm5.908-3.695c-1.027-2.086-2.042-4.14-3.082-6.25l-3.163 9.333 6.245-3.083ZM0 10c.995-.491 1.933-.958 2.873-1.422C4.47 7.791 6.069 7.002 7.67 6.22A.593.593 0 0 0 8 5.843C8.627 3.976 9.26 2.113 9.895.25c.023-.07.051-.136.095-.25.064.124.111.21.155.3 1.222 2.48 2.443 4.963 3.676 7.44.06.12.218.22.354.265 1.85.639 3.705 1.266 5.559 1.896l.267.094c-.17.087-.297.153-.424.216-2.39 1.179-4.778 2.36-7.17 3.532a.8.8 0 0 0-.449.513c-.617 1.85-1.25 3.693-1.878 5.538-.017.052-.038.1-.079.206-.16-.323-.3-.598-.436-.875-1.113-2.256-2.227-4.512-3.333-6.772a.653.653 0 0 0-.413-.368C3.962 11.364 2.11 10.73.255 10.1c-.068-.024-.134-.052-.255-.098Z" clip-rule="evenodd" />
    </svg>
  )
}

declare global {
  interface Window {
    ethereum?: MetaMaskInpageProvider
  }
}

const Navbar: React.FC = () => {

  const { isAuthenticated, authStatus } = useContext(AuthContext);

  const { address, updateAddress } = zustStore();

  const connectWallet = async () => {
    if (window.ethereum) {
      try {
        const accounts = (await window.ethereum.request({ method: 'eth_requestAccounts' })) as string[] | undefined;
        if (accounts && accounts.length > 0) {
          const address = accounts[0];
          updateAddress(address);

          // Create a message to sign
          const message = `address: ${address}`;

          // Request the user to sign the message
          const signature = await window.ethereum.request({
            method: 'personal_sign',
            params: [message, address],
          });

          const connectUrl = `${import.meta.env.VITE_API_URL}/connect`;
          const response = await fetch(connectUrl, {
            method: 'POST',
            credentials: 'include',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              address,
              message,
              signature,
            }),
          });

          const data = await response.json();
          if (response.ok) {
            console.log('Signature verified successfully:', data);
            authStatus();
          } else {
            console.error('Signature verification failed:', data);
          }
        } else {
          console.error("No accounts found");
        }
      } catch (error) {
        console.error("User denied account access or signing:", error);
      }
    } else {
      console.error("MetaMask is not installed!");
    }
  };


  const handleAccountsChanged = (accounts: string[]) => {
    if (accounts.length === 0) {
      logout();
      console.log("MetaMask is locked or no accounts are available");
    } else {
      updateAddress(accounts[0]);
      console.log("Account switched to:", accounts[0]);
    }
  };

  const logout = () => {
    updateAddress(null)
    console.log("Logged out");
  };

  useEffect(() => {
    if (window.ethereum) {
      updateAddress(window.ethereum?.selectedAddress)
      if (window.ethereum?.selectedAddress) {
      }
    }

    if (window.ethereum) {
      window.ethereum.on('accountsChanged', handleAccountsChanged as (...args: any[]) => void);
      return () => {
        if (window.ethereum) {
          window.ethereum.removeListener('accountsChanged', handleAccountsChanged as (...args: any[]) => void);
        }
      };
    }
  }, []);

  return (
    <div className="w-full bg-[var(--bg)] fixed py-3 md:px-5 px-2 flex justify-between items-center border-b-[1px] border-[var(--border-subtle)]">
      <div className="flex items-center gap-2">
        <div>
          <Logo />
        </div>
        <span className="font-semibold text">
          EtherPay
        </span>
      </div>
      {isAuthenticated
        ?
        <div>
          <DropdownMenuRoot>
            <DropdownMenuTrigger className="flex outline-none">
              <span className='outline-none'>
                <div className='w-8 h-8 rounded-full overflow-hidden'>
                  <img src="https://i.pravatar.cc/300" />
                </div>
              </span>
            </DropdownMenuTrigger>
            <DropdownMenuPortal>
              <DropdownMenuContent align="end" className="z-[999] w-44">
                <div className="p-2 text-xs text-[var(--text-subtle)] font-semibold">Account</div>
                <DropdownMenuItem>
                  <Link to={'/' + address}
                    className={clsx(
                      "w-full flex gap-2 p-2 bg-[var(--bg-muted)] hover:bg-[var(--bg-subtle)]",
                      "text-[var(--text-subtle)] hover:text-[var(--text)] rounded-[var(--radius-sm)]",
                    )}
                  >
                    <span className="flex items-center justify-center">
                      <UserRoundIcon size={16} />
                    </span>
                    <span className="text-xs font-semibold">
                      Profile
                    </span>
                  </Link>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Link to={`/${address}/settings`}
                    className={clsx(
                      "w-full flex gap-2 p-2 bg-[var(--bg-muted)] hover:bg-[var(--bg-subtle)]",
                      "text-[var(--text-subtle)] hover:text-[var(--text)] rounded-[var(--radius-sm)]",
                    )}
                  >
                    <span className="flex items-center justify-center">
                      <SettingsIcon size={16} />
                    </span>
                    <span className="text-xs font-semibold">
                      Settings
                    </span>
                  </Link>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenuPortal>
          </DropdownMenuRoot>
        </div>
        :
        <button
          className={clsx(
            "bg-[var(--bg-inverted)] rounded-full px-3 py-1 text-[var(--text-inverted)] text-sm",
          )}
          onClick={connectWallet}>
          Connect
        </button>
      }
    </div>
  )
}

export default Navbar
