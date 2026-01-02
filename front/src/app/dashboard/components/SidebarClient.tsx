"use client";
import React from "react";
import { Card } from "@/components/ui/Card/Card";
import { Button } from "@/components/forms/Button/Button";
import { FaPlus } from "react-icons/fa6";
import { Field } from "@/components/forms/Field/Field";
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

            <form onSubmit={(e) => { e.preventDefault(); setIsOpen(false); }}>
              {/* <input
                className="w-full mb-4 rounded-md border px-3 py-2"
                placeholder="Новый календарь"
              /> */}

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
                <Button type="button" variant="secondary" onClick={() => setIsOpen(false)}>Отмена</Button>
                <Button type="submit">Создать</Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </>
  );
}
