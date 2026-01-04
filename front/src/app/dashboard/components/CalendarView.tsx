"use client";

import React from "react";
import { Card } from "@/components/ui/Card/Card";
import { FaChevronLeft, FaChevronRight, FaPlus } from "react-icons/fa6";
import EventCard from "./EventCard";
import { fetchEventsAction, fetchCalendarsAction, createEventAction, updateEventAction, deleteEventAction, Event } from "./actions";
import { colorVariants } from "./CalendarList";

const WEEK_DAYS = ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"];

function startOfMonth(date: Date) {
  return new Date(date.getFullYear(), date.getMonth(), 1);
}

function endOfMonth(date: Date) {
  return new Date(date.getFullYear(), date.getMonth() + 1, 0);
}

function addDays(d: Date, days: number) {
  const r = new Date(d);
  r.setDate(r.getDate() + days);
  return r;
}

function isSameDay(a?: Date | null, b?: Date | null) {
  if (!a || !b) return false;
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate();
}

function parsePostgresTimestamp(timestamp: string): Date {
  // Handle PostgreSQL format: YYYY-MM-DD HH:MM:SS
  // Also handles ISO format for backward compatibility
  return new Date(timestamp);
}

function getDateKey(timestamp: string): string {
  // Extract date in YYYY-MM-DD format from PostgreSQL or ISO timestamp
  // PostgreSQL: "2024-01-04 15:30:45" -> "2024-01-04"
  // ISO: "2024-01-04T15:30:45Z" -> "2024-01-04"
  return timestamp.split(/[T ]/, 1)[0];
}

function getLocalDateKey(date: Date): string {
  // Convert a Date object to YYYY-MM-DD using local timezone
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

export default function CalendarView() {
  const today = new Date();
  const [current, setCurrent] = React.useState<Date>(() => startOfMonth(new Date()));
  const [events, setEvents] = React.useState<Event[]>([]);
  const [calendars, setCalendars] = React.useState<any[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [showModal, setShowModal] = React.useState(false);
  const [selectedDate, setSelectedDate] = React.useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [selectedEvent, setSelectedEvent] = React.useState<Event | null>(null);
  const [isEditingEvent, setIsEditingEvent] = React.useState(false);
  const [editingEventData, setEditingEventData] = React.useState<Partial<Event> | null>(null);

  React.useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const [eventsData, calendarsData] = await Promise.all([
          fetchEventsAction(),
          fetchCalendarsAction(),
        ]);
        console.log("Events loaded:", eventsData);
        console.log("Calendars loaded:", calendarsData);
        setEvents(eventsData);
        setCalendars(calendarsData);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load data");
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, []);

  const toIsoFromLocal = (s?: string | null) => {
        if (!s) return "";
        // datetime-local inputs are like "YYYY-MM-DDTHH:MM" or "YYYY-MM-DDTHH:MM:SS"
        // If the string contains a space (previously normalized), replace with 'T'
        let t = s.includes(" ") ? s.replace(" ", "T") : s;
        // Ensure seconds are present
        if (/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/.test(t)) {
          t = t + ":00";
        }
        const d = new Date(t); // parsed as local time
        if (isNaN(d.getTime())) return "";
        const pad = (n: number) => String(n).padStart(2, "0");
        const year = d.getFullYear();
        const month = pad(d.getMonth() + 1);
        const day = pad(d.getDate());
        const hours = pad(d.getHours());
        const minutes = pad(d.getMinutes());
        const seconds = pad(d.getSeconds());
        // timezone offset in minutes: positive for UTC+X
        const offsetMin = -d.getTimezoneOffset();
        const sign = offsetMin >= 0 ? "+" : "-";
        const absOffset = Math.abs(offsetMin);
        const offHours = pad(Math.floor(absOffset / 60));
        const offMinutes = pad(absOffset % 60);
        const tz = `${sign}${offHours}:${offMinutes}`;
        // Return local ISO-like format with timezone offset (e.g. 2026-01-06T09:00:00+07:00)
        return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${tz}`;
      };

  const prevMonth = () => setCurrent((c) => new Date(c.getFullYear(), c.getMonth() - 1, 1));
  const nextMonth = () => setCurrent((c) => new Date(c.getFullYear(), c.getMonth() + 1, 1));

  const handleAddEventClick = (dateKey: string) => {
    setSelectedDate(dateKey);
    setShowModal(true);
  };

  const handleEventCardClick = (event: Event) => {
    setSelectedEvent(event);
    setIsEditingEvent(false);
    setEditingEventData(null);
  };

  const handleEditEvent = () => {
    if (selectedEvent) {
      setEditingEventData({ ...selectedEvent });
      setIsEditingEvent(true);
    }
  };

  const handleSaveEvent = async () => {
    if (!selectedEvent || !editingEventData) return;

    setIsSubmitting(true);
    try {
      const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone || "UTC";

      const startAt = toIsoFromLocal(editingEventData.start_at);
      const endAt = toIsoFromLocal(editingEventData.end_at);

      await updateEventAction(
        selectedEvent.id,
        editingEventData.title || "",
        editingEventData.description || "",
        editingEventData.location || "",
        startAt,
        endAt,
        timezone
      );

      // Reload events
      const updatedEvents = await fetchEventsAction();
      setEvents(updatedEvents);

      // Update the selected event with new data
      const updated = updatedEvents.find(e => e.id === selectedEvent.id);
      if (updated) {
        setSelectedEvent(updated);
      }

      setIsEditingEvent(false);
      setEditingEventData(null);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to update event");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeleteEvent = async () => {
    if (!selectedEvent) return;

    if (!confirm("Are you sure you want to delete this event?")) return;

    setIsSubmitting(true);
    try {
      await deleteEventAction(selectedEvent.id);

      // Reload events
      const updatedEvents = await fetchEventsAction();
      setEvents(updatedEvents);

      setSelectedEvent(null);
      setIsEditingEvent(false);
      setEditingEventData(null);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to delete event");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleFormSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!selectedDate) return;

    setIsSubmitting(true);
    try {
      const formData = new FormData(e.currentTarget);
      const calendarId = formData.get("calendar") as string;
      const title = formData.get("title") as string;
      const description = formData.get("description") as string;
      const location = formData.get("location") as string;
      let startAtRaw = formData.get("startAt") as string;
      let endAtRaw = formData.get("endAt") as string;
      console.log("Raw timestamps:", { startAtRaw, endAtRaw });

      // Convert user's local datetime-local input into an ISO timestamp the backend can parse
      

      const startAt = toIsoFromLocal(startAtRaw);
      const endAt = toIsoFromLocal(endAtRaw);
      console.log("Converted timestamps:", { startAt, endAt });
      // detect user's timezone in browser
      const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone || "UTC";

      await createEventAction(calendarId, title, description, location, startAt, endAt, timezone);

      // Reload events
      const updatedEvents = await fetchEventsAction();
      setEvents(updatedEvents);

      setShowModal(false);
      setSelectedDate(null);
      (e.target as HTMLFormElement).reset();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Failed to create event");
    } finally {
      setIsSubmitting(false);
    }
  };

  // build range from the Monday on/before the first to the Saturday on/after the last
  const first = startOfMonth(current);
  const last = endOfMonth(current);

  // compute start: go back to Monday (1). JS getDay: 0 Sun .. 6 Sat
  const start = new Date(first);
  const startDay = start.getDay();
  const diffToMonday = startDay === 0 ? -6 : 1 - startDay;
  start.setDate(start.getDate() + diffToMonday);

  // compute end: move forward until Sunday (0)
  const end = new Date(last);
  while (end.getDay() !== 0) {
    end.setDate(end.getDate() + 1);
  }

  // iterate days from start..end, grouping Mon..Sun rows
  const weeks: Date[][] = [];
  let cur = new Date(start);
  let week: Date[] = [];
  while (cur <= end) {
    week.push(new Date(cur));

    if (week.length === 7) {
      weeks.push(week);
      week = [];
    }

    cur = addDays(cur, 1);
  }
  if (week.length) weeks.push(week);

  const monthLabel = current.toLocaleString("default", { month: "long", year: "numeric" });

  // group events by date
  const eventsByDate = new Map<string, Event[]>();
  events.forEach((event) => {
    const key = getDateKey(event.start_at);
    if (!eventsByDate.has(key)) {
      eventsByDate.set(key, []);
    }
    eventsByDate.get(key)!.push(event);
  });

  return (
    <Card variant="outline" className="m-2 p-2">
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <button
            aria-label="Previous month"
            onClick={prevMonth}
            className="p-2 rounded hover:bg-gray-100"
          >
            <FaChevronLeft />
          </button>
          <h3 className="text-lg font-medium capitalize">{monthLabel}</h3>
          <button
            aria-label="Next month"
            onClick={nextMonth}
            className="p-2 rounded hover:bg-gray-100"
          >
            <FaChevronRight />
          </button>
        </div>
        <div className="text-sm text-gray-500">Пн–Вс</div>
      </div>

      {isLoading && <div className="text-sm text-gray-400 mt-2">Загрузка событий...</div>}
      {error && <div className="text-sm text-red-500 mt-2">{error}</div>}

      <div className="grid grid-cols-7 gap-1">
        {WEEK_DAYS.map((d) => (
          <div key={d} className="text-xs font-medium text-center text-gray-400">
            {d}
          </div>
        ))}
      </div>

      <div className="flex flex-col mt-2 gap-1">
        {weeks.map((w, i) => (
          <div key={i} className="grid grid-cols-7 gap-1">
            {w.map((d) => {
              const inMonth = d.getMonth() === current.getMonth();
              const isToday = isSameDay(d, today);
              return (
                <div
                  key={d.toISOString()}
                  className={`rounded overflow-hidden border flex flex-col min-h-32 ${
                    isToday ? "bg-indigo-600 border-indigo-400 text-white" : "bg-slate-800 border-slate-700 text-slate-100"
                  } ${inMonth ? "" : "opacity-70"}`}
                >
                  <div className="flex items-end justify-center pb-1 h-8 flex-shrink-0">
                    <div className="text-sm font-medium">{d.getDate()}</div>
                  </div>
                  <div className="flex-1 overflow-y-auto px-1 pb-1 space-y-1 min-h-0 group/cell">
                    {(() => {
                      const key = getLocalDateKey(d);
                      const dayEvents = eventsByDate.get(key) || [];
                      return (
                        <>
                          {dayEvents.map((event) => (
                            <div key={event.id} className="text-xs">
                              <EventCard
                                title={event.title}
                                location={event.location}
                                start={event.start_at}
                                end={event.end_at}
                                color={event.color as keyof typeof colorVariants}
                                onClick={() => handleEventCardClick(event)}
                              />
                            </div>
                          ))}
                          <button
                            className="w-full flex items-center justify-center gap-3 p-2 rounded-md bg-slate-800 border border-slate-700 text-slate-300 text-xs opacity-0 group-hover/cell:opacity-100 transition-opacity hover:bg-slate-700"
                            onClick={() => handleAddEventClick(getLocalDateKey(d))}
                          >
                            <FaPlus size={12} />
                            <span>Add event</span>
                          </button>
                        </>
                      );
                    })()}
                  </div>
                </div>
              );
            })}
          </div>
        ))}
      </div>

      {/* Event Creation Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <Card variant="outline" className="w-full max-w-md bg-slate-900 border-slate-700">
            <div className="p-6">
              <h2 className="text-xl font-semibold text-white mb-4">Create Event</h2>
              <p className="text-sm text-slate-400 mb-4">Date: {selectedDate}</p>

              <form onSubmit={handleFormSubmit} className="space-y-4">
                {/* Calendar Selection */}
                <div>
                  <label className="block text-sm font-medium text-slate-200 mb-2">
                    Calendar
                  </label>
                  <select
                    name="calendar"
                    required
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500"
                  >
                    <option value="">Select a calendar</option>
                    {calendars.map((cal) => (
                      <option key={cal.id} value={cal.id}>
                        {cal.name}
                      </option>
                    ))}
                  </select>
                </div>

                {/* Title */}
                <div>
                  <label className="block text-sm font-medium text-slate-200 mb-2">
                    Title
                  </label>
                  <input
                    type="text"
                    name="title"
                    required
                    placeholder="Event title"
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500 placeholder-slate-500"
                  />
                </div>

                {/* Description */}
                <div>
                  <label className="block text-sm font-medium text-slate-200 mb-2">
                    Description
                  </label>
                  <textarea
                    name="description"
                    placeholder="Event description"
                    rows={3}
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500 placeholder-slate-500"
                  />
                </div>

                {/* Location */}
                <div>
                  <label className="block text-sm font-medium text-slate-200 mb-2">
                    Location
                  </label>
                  <input
                    type="text"
                    name="location"
                    placeholder="Event location"
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500 placeholder-slate-500"
                  />
                </div>

                {/* Start Timestamp */}
                <div>
                  <label className="block text-sm font-medium text-slate-200 mb-2">
                    Start Time
                  </label>
                  <input
                    type="datetime-local"
                    name="startAt"
                    required
                    defaultValue={selectedDate ? `${selectedDate}T09:00` : ""}
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500"
                  />
                </div>

                {/* End Timestamp */}
                <div>
                  <label className="block text-sm font-medium text-slate-200 mb-2">
                    End Time
                  </label>
                  <input
                    type="datetime-local"
                    name="endAt"
                    required
                    defaultValue={selectedDate ? `${selectedDate}T10:00` : ""}
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500"
                  />
                </div>

                {/* Form Actions */}
                <div className="flex gap-3 pt-4">
                  <button
                    type="button"
                    onClick={() => {
                      setShowModal(false);
                      setSelectedDate(null);
                    }}
                    className="flex-1 px-4 py-2 bg-slate-700 text-slate-200 rounded-md hover:bg-slate-600 transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={isSubmitting}
                    className="flex-1 px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                  >
                    {isSubmitting ? "Creating..." : "Create Event"}
                  </button>
                </div>
              </form>
            </div>
          </Card>
        </div>
      )}

      {/* Event Details Modal */}
      {selectedEvent && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <Card 
            variant="outline" 
            className={`w-full max-w-md bg-slate-900 border-2`}
            style={{
              borderColor: getColorForEvent(selectedEvent.color as keyof typeof colorVariants, true)
            }}
          >
            <div className="p-6">
              <div className="flex items-start justify-between mb-4">
                {isEditingEvent && editingEventData ? (
                  <input
                    type="text"
                    value={editingEventData.title || ""}
                    onChange={(e) => setEditingEventData({ ...editingEventData, title: e.target.value })}
                    className="flex-1 px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500 text-lg font-semibold"
                  />
                ) : (
                  <h2 className={`text-xl font-semibold`} style={{color: getColorForEvent(selectedEvent.color as keyof typeof colorVariants)}}>
                    {selectedEvent.title}
                  </h2>
                )}
                <button
                  onClick={() => {
                    setSelectedEvent(null);
                    setIsEditingEvent(false);
                    setEditingEventData(null);
                  }}
                  className="text-slate-400 hover:text-slate-200 transition-colors text-lg ml-2"
                >
                  ✕
                </button>
              </div>

              {isEditingEvent && editingEventData ? (
                <div className="space-y-3">
                  {/* Description */}
                  <div>
                    <label className="block text-xs font-medium text-slate-400 mb-1">Description</label>
                    <textarea
                      value={editingEventData.description || ""}
                      onChange={(e) => setEditingEventData({ ...editingEventData, description: e.target.value })}
                      rows={3}
                      className="w-full px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500 placeholder-slate-500"
                    />
                  </div>

                  {/* Location */}
                  <div>
                    <label className="block text-xs font-medium text-slate-400 mb-1">Location</label>
                    <input
                      type="text"
                      value={editingEventData.location || ""}
                      onChange={(e) => setEditingEventData({ ...editingEventData, location: e.target.value })}
                      className="w-full px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500 placeholder-slate-500"
                    />
                  </div>

                  {/* Start Time */}
                  <div>
                    <label className="block text-xs font-medium text-slate-400 mb-1">Start Time</label>
                    <input
                      type="datetime-local"
                      value={convertToDatetimeLocal(editingEventData.start_at)}
                      onChange={(e) => setEditingEventData({ ...editingEventData, start_at: e.target.value })}
                      className="w-full px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500"
                    />
                  </div>

                  {/* End Time */}
                  <div>
                    <label className="block text-xs font-medium text-slate-400 mb-1">End Time</label>
                    <input
                      type="datetime-local"
                      value={convertToDatetimeLocal(editingEventData.end_at)}
                      onChange={(e) => setEditingEventData({ ...editingEventData, end_at: e.target.value })}
                      className="w-full px-3 py-2 bg-slate-800 border border-slate-700 text-white rounded-md focus:outline-none focus:border-indigo-500"
                    />
                  </div>

                  {/* Action Buttons */}
                  <div className="flex gap-3 pt-4">
                    <button
                      onClick={() => {
                        setIsEditingEvent(false);
                        setEditingEventData(null);
                      }}
                      className="flex-1 px-4 py-2 bg-slate-700 text-slate-200 rounded-md hover:bg-slate-600 transition-colors"
                    >
                      Cancel
                    </button>
                    <button
                      onClick={handleSaveEvent}
                      disabled={isSubmitting}
                      className="flex-1 px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                      {isSubmitting ? "Saving..." : "Save"}
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  <div className="space-y-3">
                    {selectedEvent.description && (
                      <div>
                        <p className="text-xs font-medium text-slate-400 mb-1">Description</p>
                        <p className="text-sm text-slate-200">{selectedEvent.description}</p>
                      </div>
                    )}

                    {selectedEvent.location && (
                      <div>
                        <p className="text-xs font-medium text-slate-400 mb-1">Location</p>
                        <p className="text-sm text-slate-200">{selectedEvent.location}</p>
                      </div>
                    )}

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <p className="text-xs font-medium text-slate-400 mb-1">Start</p>
                        <p className="text-sm text-slate-200">{formatEventDateTime(selectedEvent.start_at)}</p>
                      </div>
                      <div>
                        <p className="text-xs font-medium text-slate-400 mb-1">End</p>
                        <p className="text-sm text-slate-200">{formatEventDateTime(selectedEvent.end_at)}</p>
                      </div>
                    </div>

                    {selectedEvent.calendarId && (
                      <div>
                        <p className="text-xs font-medium text-slate-400 mb-1">Calendar</p>
                        <p className="text-sm text-slate-200">
                          {calendars.find(c => c.id === selectedEvent.calendarId)?.name || "Unknown"}
                        </p>
                      </div>
                    )}
                  </div>

                  <div className="flex gap-3 mt-6">
                    <button
                      onClick={handleDeleteEvent}
                      disabled={isSubmitting}
                      className="flex-1 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                      {isSubmitting ? "Deleting..." : "Delete"}
                    </button>
                    <button
                      onClick={handleEditEvent}
                      className="flex-1 px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-500 transition-colors"
                    >
                      Edit
                    </button>
                  </div>
                </>
              )}
            </div>
          </Card>
        </div>
      )}
    </Card>
  );
}

function getColorForEvent(colorKey: keyof typeof colorVariants, desaturate: boolean = false): string {
  const colorClasses: Record<string, string> = {
    red: desaturate ? "#a85555" : "#ef4444",
    orange: desaturate ? "#a87755" : "#f97316",
    yellow: desaturate ? "#aaa855" : "#eab308",
    green: desaturate ? "#5aa855" : "#22c55e",
    blue: desaturate ? "#5581a8" : "#3b82f6",
    indigo: desaturate ? "#6b5ba8" : "#6366f1",
    purple: desaturate ? "#8b5ba8" : "#a855f7",
    pink: desaturate ? "#a85582" : "#ec4899",
  };
  return colorClasses[colorKey as string] || (desaturate ? "#888888" : "#3b82f6");
}

function formatEventDateTime(timestamp: string): string {
  try {
    // Handle both ISO format and PostgreSQL format
    let d: Date;
    if (timestamp.includes(" ")) {
      // PostgreSQL format: YYYY-MM-DD HH:MM:SS
      const [datePart, timePart] = timestamp.split(" ");
      d = new Date(`${datePart}T${timePart}`);
    } else {
      d = new Date(timestamp);
    }
    
    if (isNaN(d.getTime())) return timestamp;
    
    return d.toLocaleString([], {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch {
    return timestamp;
  }
}

function convertToDatetimeLocal(timestamp: string | undefined): string {
  if (!timestamp) return "";
  
  try {
    let d: Date;
    if (timestamp.includes(" ")) {
      // PostgreSQL format: YYYY-MM-DD HH:MM:SS
      const [datePart, timePart] = timestamp.split(" ");
      d = new Date(`${datePart}T${timePart}`);
    } else {
      d = new Date(timestamp);
    }
    
    if (isNaN(d.getTime())) return "";
    
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, "0");
    const day = String(d.getDate()).padStart(2, "0");
    const hours = String(d.getHours()).padStart(2, "0");
    const minutes = String(d.getMinutes()).padStart(2, "0");
    
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  } catch {
    return "";
  }
}
