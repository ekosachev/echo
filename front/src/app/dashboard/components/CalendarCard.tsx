import { Card, } from "@/components/ui/Card/Card";
import React from "react";
import { FaRegCalendar } from "react-icons/fa6";

export interface CalendarCardProps extends React.HTMLAttributes<HTMLDivElement> {
  iconColor: string,
  calendarName: string,
}

const CalendarCard = React.forwardRef<HTMLDivElement, CalendarCardProps>(({
  iconColor,
  calendarName,
  ...props
}, ref) => {
  return (

    <Card variant="default" hoverable {...props}>
      <div ref={ref} className="flex flex-row justify-start items-center gap-2 text-lg">
        <FaRegCalendar className={iconColor + ' text-2xl'}></FaRegCalendar>
        <p className={iconColor + ' font-bold saturate-[.30]'}>{calendarName}</p>
      </div>
    </Card>
  )
})

CalendarCard.displayName = 'CalendarCard'

export { CalendarCard }

