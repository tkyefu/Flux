"use client";

import { useState } from "react";
import { apiFetch } from "@/lib/api";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState<string>("");

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setMessage("");
    try {
      const res = await apiFetch<{ token: string; user: { name: string } }>("/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });
      if (res?.token) {
        localStorage.setItem("token", res.token);
        setMessage(`Logged in as ${res.user?.name ?? "user"}`);
      } else {
        setMessage("Login succeeded (no token returned)");
      }
    } catch (err: any) {
      setMessage(err.message ?? "Login failed");
    }
  }

  return (
    <div className="mx-auto max-w-md p-6">
      <h1 className="text-2xl font-bold mb-4">Login</h1>
      <form onSubmit={onSubmit} className="space-y-3">
        <input className="w-full border p-2" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
        <input className="w-full border p-2" placeholder="Password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <button className="bg-black text-white px-4 py-2" type="submit">Login</button>
      </form>
      {message && <p className="mt-3 text-sm">{message}</p>}
    </div>
  );
}
