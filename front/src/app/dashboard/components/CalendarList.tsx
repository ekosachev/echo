import { authFetch } from "@/lib/authFetch";
import { CalendarCard } from "./CalendarCard";

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

export default async function CalendarList() {
  const res = await authFetch("/api/v1/calendars");
  const calendars = await res!.json();
  return (

    <ul style={{ listStyleType: "none" }} className="flex flex-row gap-2 overflow-auto md:flex-col ">
      {calendars.calendars.map((calendar: any) => (
        <li key={calendar.id} className="block">
          <CalendarCard  iconColor={colorVariants[calendar.color as keyof typeof colorVariants]} calendarName={calendar.name}></CalendarCard>
        </li>
      ))}

    </ul>
  )
}
