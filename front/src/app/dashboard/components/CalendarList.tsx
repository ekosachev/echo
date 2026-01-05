"use client";

import React from "react";
import { CalendarCard } from "./CalendarCard";
import { fetchCalendarsAction, deleteCalendarAction, updateCalendarAction } from "./actions";

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
  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [selectedCalendar, setSelectedCalendar] = React.useState<any | null>(null);
  const [isEditing, setIsEditing] = React.useState(false);
  const [formName, setFormName] = React.useState("");
  const [formColor, setFormColor] = React.useState("");

  React.useEffect(() => {
    const loadCalendars = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const data = await fetchCalendarsAction();
        setCalendars(data || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load calendars");
      } finally {
        setIsLoading(false);
      }
    };

    loadCalendars();
  }, [reloadTrigger]);

  const openModal = (calendar: any) => {
    setSelectedCalendar(calendar);
    setFormName(calendar.name || "");
    setFormColor(calendar.color || "");
    setIsEditing(false);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setSelectedCalendar(null);
    setIsEditing(false);
  };

  const handleDelete = async () => {
    if (!selectedCalendar) return;
    try {
      await deleteCalendarAction(selectedCalendar.id);
      setCalendars(prev => prev.filter(c => c.id !== selectedCalendar.id));
      closeModal();
    } catch (err) {
      console.error(err);
    }
  };

  const handleSave = async () => {
    if (!selectedCalendar) return;
    try {
      await updateCalendarAction(selectedCalendar.id, formName, formColor);
      setCalendars(prev => prev.map(c => c.id === selectedCalendar.id ? { ...c, name: formName, color: formColor } : c));
      setIsEditing(false);
      closeModal();
    } catch (err) {
      console.error(err);
    }
  };

  if (isLoading) {
    return <div className="text-sm text-gray-500">Загрузка календарей...</div>;
  }

  if (error) {
    return <div className="text-sm text-red-500">{error}</div>;
  }

  console.log("Calendars loaded:", calendars);
  return (
    <>
      <ul style={{ listStyleType: "none" }} className="flex flex-row gap-2 overflow-auto md:flex-col ">
        {calendars.map((calendar: any) => (
          <li key={calendar.id} className="block">
            <CalendarCard
              onClick={() => openModal(calendar)}
              iconColor={colorVariants[calendar.color as keyof typeof colorVariants]}
              calendarName={calendar.name}
            ></CalendarCard>
          </li>
        ))}
      </ul>

      {isModalOpen && selectedCalendar && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="absolute inset-0 bg-black/40" onClick={closeModal} />
          <div className="relative w-full max-w-md mx-4 bg-white dark:bg-gray-800 rounded shadow p-4">
            <h3 className="text-lg font-semibold mb-2">Calendar details</h3>

            {!isEditing ? (
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <span className={colorVariants[selectedCalendar.color as keyof typeof colorVariants] + " font-bold"}>
                    {selectedCalendar.name}
                  </span>
                </div>
                <div className="text-sm text-gray-500">Color: <span className={colorVariants[selectedCalendar.color as keyof typeof colorVariants]}>{selectedCalendar.color}</span></div>
              </div>
            ) : (
              <div className="space-y-3">
                <div>
                  <label className="block text-sm font-medium mb-1">Name</label>
                  <input value={formName} onChange={e => setFormName(e.target.value)} className="w-full border rounded px-2 py-1" />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">Color</label>
                  <select value={formColor} onChange={e => setFormColor(e.target.value)} className="w-full border rounded px-2 py-1">
                    {Object.keys(colorVariants).map(c => (
                      <option key={c} value={c}>{c}</option>
                    ))}
                  </select>
                </div>
              </div>
            )}

            <div className="mt-4 flex justify-between gap-2">
              <button onClick={handleDelete} className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded">Delete</button>
              <div className="flex gap-2">
                {!isEditing ? (
                  <>
                    <button onClick={() => setIsEditing(true)} className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded">Edit</button>
                    <button onClick={closeModal} className="bg-gray-200 hover:bg-gray-300 px-3 py-1 rounded">Close</button>
                  </>
                ) : (
                  <>
                    <button onClick={handleSave} className="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded">Save</button>
                    <button onClick={() => { setIsEditing(false); setFormName(selectedCalendar.name); setFormColor(selectedCalendar.color); }} className="bg-gray-200 hover:bg-gray-300 px-3 py-1 rounded">Cancel</button>
                  </>
                )}
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
