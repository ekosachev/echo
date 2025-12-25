import React from "react";
import { Card } from "@/components/ui/Card/Card";
import { FaRegUser, FaRegCalendar } from "react-icons/fa6";
import CalendarList from "./CalendarList";
import { authFetch } from "@/lib/authFetch";

export default async function Sidebar() {
  const res = await authFetch("/api/v1/auth/me");
  const data = await res!.json();
  console.log(data);
  return (
    <div className='flex flex-col h-full gap-4 m-2 ' >

      <Card variant='outline' hoverable>
        <div className="flex flex-row justify-between items-center gap-2">
          <p className="text-(--text) text-lg">{data.email}</p>
          <FaRegUser></FaRegUser>
        </div>
      </Card >

      <Card variant="outline" className="grow" >
        <h1 className="text-(--text) text-xl mb-2">Мои календари</h1>
        <CalendarList></CalendarList>
      </Card>
    </div>
  )
}

