import clsx from "clsx";
import { zustStore } from "../../App"
import { useParams } from "react-router-dom";
import { useState } from "react";
import { ethers } from "ethers";
import { toast } from "sonner"

const About = () => {

  const [formData, setFormData] = useState({
    name: "",
    amount: "",
    message: "",
    transactionHash: "",
  })

  function convertToWei(ethAmount: string) {
    return ethers.parseEther(ethAmount);
  }

  const { address, viewingProfileData } = zustStore(); // address of the user
  const { user_address } = useParams(); // address of the user whose profile is being viewed

  let recipientAddress = user_address;

  let senderAddress = address;

  const submitTransaction = async (
    transactionHash: string,
    amount: string,
    senderAddress: string,
    receiverAddress: string,
    message: string
  ) => {

    const requestData = {
      hash: transactionHash,
      sender_address: senderAddress,
      receiver_address: receiverAddress,
      amount: amount,
      message: message,
    };

    try {
      const response = await fetch(
        `${import.meta.env.VITE_API_URL}/transaction`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        throw new Error('Failed to fetch transaction status');
      }

      const data = await response.json();
      console.log('data ->', data)
    } catch (err) {
      console.error("Failed to submit transaction", err)
    }
  };

  const initTransaction = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    try {
      if (window.ethereum) {
        let params = [
          {
            from: senderAddress, // sender address
            to: recipientAddress, // recipient address
            value: Number(convertToWei(formData.amount)).toString(16),
            gasLimit: '0x5028',
            maxPriorityFeePerGas: '0x3b9aca00',
            maxFeePerGas: '0x2540be400',
          }
        ];
        await window.ethereum.request({
          method: "eth_sendTransaction", params
        })
          .then((txHash) => {
            if (typeof txHash === 'string') {
              setFormData({ ...formData, transactionHash: txHash as string });
              console.log("transactionHash", txHash)
              toast.success("Transaction submitted successfully")
              if (senderAddress && recipientAddress) {
                console.log("triggering submit transaction")
                submitTransaction(txHash, formData.amount, senderAddress, recipientAddress, formData.message)
              }
            } else {
              console.error('Transaction hash is not a string:', txHash);
            }
          })
          .catch((error) => console.error(error));
      }
    } catch (err) {
      console.log(err)
    }
  }

  return (
    <div className="w-full flex gap-4">
      <div className="w-2/5 h-full bg-[var(--bg-muted)] rounded-[var(--radius-md)] p-4">
        <div className="flex flex-col gap-4">
          <span className="font-semibold text-md">About</span>
          <p>
            hi there.
            {viewingProfileData?.about}
          </p>
        </div>
      </div>
      <div className="w-3/5 bg-[var(--bg-muted)] rounded-[var(--radius-md)] p-4">
        <span className="font-semibold text-md">Donate</span>
        <form className="py-4 flex flex-col gap-5">
          <div className="flex flex-col gap-1.5">
            <label className="font-semibold text-xs text-[var(--text-subtle)]">Name</label>
            <input
              id="name"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="Enter your Name"
              className={clsx(
                'outline-none p-2 w-full bg-[var(--bg-subtle)] placeholder:font-[Geist] placeholder:text-sm',
                'text-base border border-[var(--border-subtle)] rounded-[var(--radius-sm)]',
              )}
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <label className="font-semibold text-xs text-[var(--text-subtle)]">Amount</label>
            <input
              required
              id="amount"
              onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
              placeholder="0"
              type="number"
              className={clsx(
                'outline-none p-2 w-full bg-[var(--bg-subtle)] placeholder:font-[Geist] placeholder:text-sm',
                'text-base border border-[var(--border-subtle)] rounded-[var(--radius-sm)]',
              )}
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <label className="font-semibold text-xs text-[var(--text-subtle)]">Message</label>
            <textarea
              id="message"
              value={formData.message}
              onChange={(e) => setFormData({ ...formData, message: e.target.value })}
              placeholder="Add a message"
              className={clsx(
                'outline-none p-2 w-full bg-[var(--bg-subtle)] placeholder:font-[Geist] placeholder:text-sm',
                'text-base border border-[var(--border-subtle)] rounded-[var(--radius-sm)]',
              )}
            />
          </div>
          <div className="">
            <button
              onClick={initTransaction}
              className="p-3 font-semibold text-sm w-full border rounded-md bg-[var(--bg-inverted)] text-[var(--text-inverted)]">
              Send Donation
            </button>
          </div>
        </form>
      </div>
    </div >
  )
};

export default About;
