"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";

type Task = { id: number; title: string; description?: string; status?: string; user_id?: number };

export default function TasksPage() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState("pending");
  const [message, setMessage] = useState("");

  async function load() {
    try {
      const data = await apiFetch<Task[]>('/tasks');
      setTasks(data ?? []);
    } catch (err: any) {
      setMessage(err.message);
    }
  }

  useEffect(() => { load(); }, []);

  async function createTask(e: React.FormEvent) {
    e.preventDefault();
    setMessage("");
    try {
      await apiFetch('/tasks', { method: 'POST', body: JSON.stringify({ title, description, status }) });
      setTitle(""); setDescription(""); setStatus("pending");
      await load();
    } catch (err: any) {
      setMessage(err.message);
    }
  }

  return (
    <div className="mx-auto max-w-2xl p-6 space-y-6">
      <h1 className="text-2xl font-bold">Tasks</h1>

      <form onSubmit={createTask} className="space-y-3">
        <input className="w-full border p-2" placeholder="Title" value={title} onChange={(e) => setTitle(e.target.value)} />
        <input className="w-full border p-2" placeholder="Description" value={description} onChange={(e) => setDescription(e.target.value)} />
        <select className="w-full border p-2" value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="pending">pending</option>
          <option value="in_progress">in_progress</option>
          <option value="done">done</option>
        </select>
        <button className="bg-black text-white px-4 py-2" type="submit">Create</button>
      </form>

      {message && <p className="text-sm">{message}</p>}

      <ul className="space-y-2">
        {tasks.map(t => (
          <li key={t.id} className="border p-3 rounded">
            <div className="font-medium">{t.title}</div>
            {t.description && <div className="text-sm text-gray-600">{t.description}</div>}
            <div className="text-xs text-gray-500">{t.status}</div>
          </li>
        ))}
      </ul>
    </div>
  );
}
