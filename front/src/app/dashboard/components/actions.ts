"use server";

import { authFetch } from "@/lib/authFetch";

export interface Calendar {
  name: string
  color: string
  id: string
}

export interface Event {
  id: string
  title: string
  description: string
  location: string
  start_at: string
  end_at: string
  allDay: boolean
  calendarId: string
  color: string
}

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

export async function deleteCalendarAction(calendarId: string) {
  const res = await authFetch(`/api/v1/calendars/${calendarId}`, {
    method: "DELETE",
  });

  if (!res.ok) {
    const error = await res.json();
    console.error("Error deleting calendar:", error);
    throw new Error(error.message || "Failed to delete calendar");
  }
}

export async function updateCalendarAction(calendarId: string, name: string, color: string) {
  const res = await authFetch(`/api/v1/calendars/${calendarId}`, {
    method: "PUT",
    body: JSON.stringify({ name, color }),
  });

  if (!res.ok) {
    const error = await res.json();
    console.error("Error updating calendar:", error);
    throw new Error(error.message || "Failed to update calendar");
  }

  return await res.json();
}

export async function fetchCalendarsAction(): Promise<Calendar[]> {
  const res = await authFetch("/api/v1/calendars");
  
  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || "Failed to fetch calendars");
  }


  const calendars = (await res.json()).calendars as Calendar[];

  return calendars;
}

async function fetchEventsForCalendar(calendarId: string, color: string): Promise<Event[]> {
  const res = await authFetch(`/api/v1/calendars/${calendarId}/events`);

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || `Failed to fetch events for calendar ${calendarId}`);
  }

  const events = (await res.json()).events as Event[];
  events.map(event => { event.calendarId = calendarId; event.color = color });
  return events;
}

export async function fetchEventsAction() {
  const calendars = await fetchCalendarsAction();
  const allEvents: Event[] = [];
  for (const calendar of calendars) {
    const events = await fetchEventsForCalendar(calendar.id, calendar.color);
    allEvents.push(...events);
  }
  return allEvents;
}

export async function createEventAction(
  calendarId: string,
  title: string,
  description: string,
  location: string,
  startAt: string,
  endAt: string,
  timezone?: string
) {
  const res = await authFetch(`/api/v1/events`, {
    method: "POST",
    body: JSON.stringify({
      calendar_id: calendarId,
      title,
      description,
      location,
      start_at: startAt,
      end_at: endAt,
      timezone,
    }),
  });

  if (!res.ok) {
    const error = await res.json();
    console.error("Error creating event:", error);
    throw new Error(error.message || "Failed to create event");
  }

  return await res.json();
}

export async function updateEventAction(
  eventId: string,
  title: string,
  description: string,
  location: string,
  startAt: string,
  endAt: string,
  timezone?: string
) {
  const res = await authFetch(`/api/v1/events/${eventId}`, {
    method: "PUT",
    body: JSON.stringify({
      title,
      description,
      location,
      start_at: startAt,
      end_at: endAt,
      timezone,
    }),
  });

  if (!res.ok) {
    const error = await res.json();
    console.error("Error updating event:", error);
    throw new Error(error.message || "Failed to update event");
  }

  return await res.json();
}

export async function deleteEventAction(eventId: string) {
  const res = await authFetch(`/api/v1/events/${eventId}`, {
    method: "DELETE",
  });

  if (!res.ok) {
    const error = await res.json();
    console.error("Error deleting event:", error);
    throw new Error(error.message || "Failed to delete event");
  }
}