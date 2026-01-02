"use server";
import React from "react";
import { Card } from "@/components/ui/Card/Card";

import CalendarList from "./CalendarList";
import SidebarClient from "./SidebarClient";

import { authFetch } from "@/lib/authFetch";
import { FaRegUser } from "react-icons/fa6";
import { Button } from "@/components/forms/Button/Button";

export default async function Sidebar() {
  const res = await authFetch("/api/v1/auth/me");
  const data = await res!.json();
  return (
    <div className='flex flex-col h-full gap-4 m-2 ' >

      <Card variant='outline' hoverable>
        <div className="flex flex-row justify-between items-center gap-2">
          <p className="text-(--text) text-lg">{data.email}</p>
          <FaRegUser></FaRegUser>
        </div>
      </Card>

      <Card variant="outline" className="grow" >
        <SidebarClient />
        <CalendarList></CalendarList>
      </Card>
    </div>
  )
}

