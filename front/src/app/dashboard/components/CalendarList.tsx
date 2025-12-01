import { CalendarCard } from "./CalendarCard";

export default function CalendarList() {
  return (

    <ul style={{ listStyleType: "none" }} className="flex flex-row gap-2 overflow-scroll md:flex-col">
      <li className="block">
        <CalendarCard iconColor="text-emerald-300" calendarName="Calendar 1"></CalendarCard>
      </li>


      <li className="block">

        <CalendarCard iconColor="text-fuchsia-300" calendarName="Calendar 2"></CalendarCard>
      </li>


      <li className="block">

        <CalendarCard iconColor="text-amber-300" calendarName="Calendar 3"></CalendarCard>
      </li>

    </ul>
  )
}
