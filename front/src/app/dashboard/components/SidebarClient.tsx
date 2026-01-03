"use client";
import React from "react";
import { Card } from "@/components/ui/Card/Card";
import { Button } from "@/components/forms/Button/Button";
import { FaPlus } from "react-icons/fa6";
import { Field } from "@/components/forms/Field/Field";
import { createCalendarAction } from "./actions";
import CalendarList from "./CalendarList";
// import { colorVariants } from "./CalendarList";
export const colorVariants = {
  red: 'bg-red-500',
  orange: 'bg-orange-500',
  amber: 'bg-amber-500',
  yellow: 'bg-yellow-500',
  lime: 'bg-lime-500',
  green: 'bg-green-500',
  emerald: 'bg-emerald-500',
  teal: 'bg-teal-500',
  cyan: 'bg-cyan-500',
  sky: 'bg-sky-500',
  blue: 'bg-blue-500',
  indigo: 'bg-indigo-500',
  violet: 'bg-violet-500',
  purple: 'bg-purple-500',
  fuchsia: 'bg-fuchsia-500',
  pink: 'bg-pink-500',
  rose: 'bg-rose-500',
}

export default function SidebarClient() {
  const [isOpen, setIsOpen] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [isLoading, setIsLoading] = React.useState(false);
  const [reloadTrigger, setReloadTrigger] = React.useState(0);
  const _colorKeys = Object.keys(colorVariants) as Array<keyof typeof colorVariants>;
  const [selectedColor, setSelectedColor] = React.useState<string>(_colorKeys[0] ?? 'blue');

  return (
    <>
      <div className="flex flex-row justify-between items-center mb-2">
        <h1 className="text-(--text) text-xl mb-2">Мои календари</h1>
        <Button
          className="w-9 h-9 flex items-center justify-center  p-0 text-2xl"
          onClick={() => setIsOpen(true)}
          variant="secondary"
          aria-label="Добавить календарь"
        >
          <FaPlus></FaPlus>
        </Button>
      </div>

      <CalendarList reloadTrigger={reloadTrigger} />

      {isOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div
            className="absolute inset-0 bg-black/50"
            onClick={() => setIsOpen(false)}
          />

          <Card className="z-10 w-full max-w-md p-6" variant="outline">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg">Создать календарь</h2>
              <button
                aria-label="Закрыть"
                onClick={() => setIsOpen(false)}
                className="text-xl px-2"
              >
                ×
              </button>
            </div>

            <form
              onSubmit={async (e) => {
                e.preventDefault();
                setError(null);
                setIsLoading(true);

                const formData = new FormData(e.currentTarget);
                const name = formData.get("name") as string;
                const color = formData.get("color") as string;

                try {
                  await createCalendarAction(name, color);
                  setIsOpen(false);
                  setReloadTrigger((prev) => prev + 1);
                  (e.target as HTMLFormElement).reset();
                  setSelectedColor(_colorKeys[0] ?? 'blue');
                } catch (err) {
                  setError(
                    err instanceof Error ? err.message : "An error occurred"
                  );
                } finally {
                  setIsLoading(false);
                }
              }}
            >
              {error && (
                <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded-md text-sm">
                  {error}
                </div>
              )}

              <Field type='text' name='name' placeholder="Название" required></Field>
              <p>Выберите цвет:</p>
              <div className="flex flex-row flex-wrap gap-2 mb-4">
                <input type="hidden" name="color" value={selectedColor} />
                {_colorKeys.map((k) => {
                  const bg = (colorVariants as any)[k];
                  const isSelected = selectedColor === k;
                  return (
                    <label key={k} className="inline-flex items-center cursor-pointer" title={String(k)}>
                      <input
                        type="radio"
                        name="color"
                        value={String(k)}
                        className="sr-only"
                        checked={isSelected}
                        onChange={() => setSelectedColor(String(k))}
                      />
                      <div className={`${bg} w-6 h-6 rounded-sm ${isSelected ? 'ring-2 ring-offset-2 ring-indigo-500' : ''}`} />
                    </label>
                  );
                })}
              </div>

              <div className="flex justify-end gap-2">
                <Button type="button" variant="secondary" onClick={() => setIsOpen(false)} disabled={isLoading}>Отмена</Button>
                <Button type="submit" disabled={isLoading}>
                  {isLoading ? "Создание..." : "Создать"}
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </>
  );
}
