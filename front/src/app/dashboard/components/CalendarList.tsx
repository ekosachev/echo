"use client";

import React from "react";
import { CalendarCard } from "./CalendarCard";
import { fetchCalendarsAction } from "./actions";

export const colorVariants = {
  red: 'text-red-500',
  orange: 'text-orange-500',
  amber: 'text-amber-500',
  yellow: 'text-yellow-500',
  lime: 'text-lime-500',
  green: 'text-green-500',
  emerald: 'text-emerald-500',
  teal: 'text-teal-500',
  cyan: 'text-cyan-500',
  sky: 'text-sky-500',
  blue: 'text-blue-500',
  indigo: 'text-indigo-500',
  violet: 'text-violet-500',
  purple: 'text-purple-500',
  fuchsia: 'text-fuchsia-500',
  pink: 'text-pink-500',
  rose: 'text-rose-500',
}

export default function CalendarList({ reloadTrigger = 0 }: { reloadTrigger?: number } = {}) {
  const [calendars, setCalendars] = React.useState<any[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    const loadCalendars = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const data = await fetchCalendarsAction();
        setCalendars(data.calendars || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load calendars");
      } finally {
        setIsLoading(false);
      }
    };

    loadCalendars();
  }, [reloadTrigger]);

  if (isLoading) {
    return <div className="text-sm text-gray-500">Загрузка календарей...</div>;
  }

  if (error) {
    return <div className="text-sm text-red-500">{error}</div>;
  }

  console.log("Calendars loaded:", calendars);
  return (
    <ul style={{ listStyleType: "none" }} className="flex flex-row gap-2 overflow-auto md:flex-col ">
      {calendars.map((calendar: any) => (
        <li key={calendar.id} className="block">
          <CalendarCard
            iconColor={colorVariants[calendar.color as keyof typeof colorVariants]}
            calendarName={calendar.name}
          ></CalendarCard>
        </li>
      ))}
    </ul>
  );
}
