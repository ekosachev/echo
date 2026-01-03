"use server";

import { authFetch } from "@/lib/authFetch";

export async function createCalendarAction(name: string, color: string) {
  const res = await authFetch("/api/v1/calendars", {
    method: "POST",
    body: JSON.stringify({ name, color }),
  });

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || "Failed to create calendar");
  }

  return await res.json();
}

export async function fetchCalendarsAction() {
  const res = await authFetch("/api/v1/calendars");
  
  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || "Failed to fetch calendars");
  }

  console.log("Fetch calendars response:", res);
  return await res.json();
}
