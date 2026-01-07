"use client";

import React from "react";
import { colorVariants } from "./CalendarList";

type EventCardProps = {
  title: string;
  location?: string;
  start: string | Date;
  end?: string | Date;
  color?: keyof typeof colorVariants;
  onClick?: () => void;
};

function fmtTime(t: string | Date) {
  let d: Date;
  if (typeof t === "string") {
    // Handle both ISO and PostgreSQL format (YYYY-MM-DD HH:MM:SS)
    d = new Date(t);
    // If the date is invalid and looks like PostgreSQL format, parse it manually
    if (isNaN(d.getTime()) && /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}/.test(t)) {
      const [datePart, timePart] = t.split(" ");
      d = new Date(`${datePart}T${timePart}Z`);
    }
  } else {
    d = t;
  }
  
  try {
    return d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  } catch {
    return String(t);
  }
}

export default function EventCard({ title, location, start, end, color = "blue", onClick }: EventCardProps) {
  const textClass = (colorVariants as any)[color] ?? "text-blue-500";
  // convert text-... to bg-... for a colored bar
  const bgClass = textClass.replace(/^text-/, "bg-");

  return (
    <div 
      onClick={onClick}
      className="flex items-center gap-3 p-2 rounded-md bg-slate-800 border border-slate-700 cursor-pointer hover:bg-slate-700 transition-colors"
    >
      <div className={`w-2 h-12 rounded ${bgClass}`} />

      <div className="flex-1 min-w-0">
        <div className={`text-sm font-semibold truncate ${textClass}`}>{title}</div>
        {location && <div className="text-xs text-slate-300 truncate">{location}</div>}
      </div>

      <div className="text-xs text-slate-300 text-center">
        <div>{fmtTime(start)}</div>
        {end && <div>- {fmtTime(end)}</div>}
      </div>
    </div>
  );
}
