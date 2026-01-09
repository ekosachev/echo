"use client";

import React from "react";
import CalendarList from "./CalendarList";

export function CalendarListWrapper({
  reloadTrigger,
}: {
  reloadTrigger: number;
}) {
  const [key, setKey] = React.useState(0);

  React.useEffect(() => {
    setKey((prev) => prev + 1);
  }, [reloadTrigger]);

  return <CalendarList key={key} />;
}
