import clsx from "clsx";
import { AuthContext } from "../../AuthWrapper"
import { zustStore } from "../../App"
import { useContext, useEffect } from "react"
import { useNavigate } from "react-router-dom"


const Logo = () => {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
      <path fill="#fff" d="M14.86 5.88a2.88 2.88 0 1 1-5.76 0 2.88 2.88 0 0 1 5.76 0Zm0 6.12a2.88 2.88 0 1 1-5.76 0 2.88 2.88 0 0 1 5.76 0Zm-2.88 9a2.88 2.88 0 1 0 0-5.76 2.88 2.88 0 0 0 0 5.76Zm6.74-9.566a2.88 2.88 0 1 1-2.88-4.988 2.88 2.88 0 0 1 2.88 4.988Z" />
      <path fill="#fff" d="M4.186 16.5a2.88 2.88 0 1 0 4.989-2.88 2.88 2.88 0 0 0-4.989 2.88Zm11.654 1.054a2.88 2.88 0 1 1 2.88-4.988 2.88 2.88 0 0 1-2.88 4.988Z" />
      <path fill="#fff" d="M4.186 7.5a2.88 2.88 0 1 0 4.989 2.88A2.88 2.88 0 0 0 4.186 7.5Z" />
    </svg>
  )
}

const Home = () => {
  const { authStatus, isAuthenticated, user } = useContext(AuthContext);

  const { updateAddress } = zustStore();

  const navigate = useNavigate();

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
            navigate('/' + address);
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

  useEffect(() => {
    if (isAuthenticated) {
      navigate(`/${user?.address}`)
    }
  }, [isAuthenticated, navigate])

  return (
    <>
      <div className="relative w-full h-[100vh]">
        {/* Background image */}
        <div className="absolute inset-0 bg-[url('/glimmer.svg')] bg-cover bg-center z-[1]" />

        {/* Video layer */}
        <div className="absolute inset-0 z-[2] opacity-50 invert">
          <video autoPlay loop muted className="w-full h-full object-cover">
            <source src="https://res.cloudinary.com/dhgck7ebz/video/upload/f_auto:video,q_auto/v1/Non%20Beta%20Release/Animations/Wallet_01" />
          </video>
        </div>

        {/* Text content (above the video) */}
        <div className="relative z-[3] flex flex-col w-full h-full items-center justify-center gap-4">
          <h1 className="text-7xl font-bold w-[640px] text-center text-white">
            Empower Your Craft with Crypto
          </h1>
          <p className="text-2xl text-white">
            Let the Blockchain Fund Your Next Big Idea
          </p>
          <div className="mt-6">
            <button
              onClick={connectWallet}
              className={clsx(
                "px-12 py-6 rounded-full bg-[var(--bg-inverted)] text-[var(--text-inverted)] text-2xl font-semibold",
                "hover:bg-[var(--bg-inverted-emphasis)] transition-all duration-300"
              )}>
              Start My Page
            </button>
          </div>
          <div className="text-md text-white">
            Completely free and just takes a moment!
          </div>
        </div>

        {/* Fixed header */}
        <div className="fixed top-0 left-0 px-6 py-6 w-full z-[4]">
          <div className="flex items-center gap-2">
            <div>
              <Logo />
            </div>
            <span className="text-lg font-semibold text-white">EtherPay</span>
          </div>
        </div>
      </div>
    </>
  )
};

export default Home;
