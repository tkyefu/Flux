"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";

export default function ResetPage() {
  const [email, setEmail] = useState("");
  const [token, setToken] = useState("");
  const [pwd1, setPwd1] = useState("");
  const [pwd2, setPwd2] = useState("");
  const [message, setMessage] = useState("");

  useEffect(() => {
    const p = new URLSearchParams(window.location.search);
    const t = p.get("token");
    if (t) setToken(t);
  }, []);

  async function requestReset(e: React.FormEvent) {
    e.preventDefault();
    setMessage("");
    try {
      await apiFetch("/auth/forgot-password", { method: "POST", body: JSON.stringify({ email }) });
      setMessage("If the email exists, a reset link has been sent.");
    } catch (err: any) {
      setMessage(err.message);
    }
  }

  async function doReset(e: React.FormEvent) {
    e.preventDefault();
    setMessage("");
    try {
      await apiFetch("/auth/reset-password", {
        method: "POST",
        body: JSON.stringify({ token, new_password: pwd1, confirm_password: pwd2 }),
      });
      setMessage("Password reset successful. You can login now.");
    } catch (err: any) {
      setMessage(err.message);
    }
  }

  return (
    <div className="mx-auto max-w-md p-6 space-y-8">
      <div>
        <h2 className="text-xl font-bold mb-2">Request Reset</h2>
        <form onSubmit={requestReset} className="space-y-3">
          <input className="w-full border p-2" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
          <button className="bg-black text-white px-4 py-2" type="submit">Request</button>
        </form>
      </div>

      <div>
        <h2 className="text-xl font-bold mb-2">Reset Password</h2>
        <form onSubmit={doReset} className="space-y-3">
          <input className="w-full border p-2" placeholder="Token" value={token} onChange={(e) => setToken(e.target.value)} />
          <input className="w-full border p-2" placeholder="New Password" type="password" value={pwd1} onChange={(e) => setPwd1(e.target.value)} />
          <input className="w-full border p-2" placeholder="Confirm Password" type="password" value={pwd2} onChange={(e) => setPwd2(e.target.value)} />
          <button className="bg-black text-white px-4 py-2" type="submit">Reset</button>
        </form>
      </div>

      {message && <p className="text-sm">{message}</p>}
    </div>
  );
}
