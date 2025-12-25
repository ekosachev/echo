"use client";

import { Card } from "@/components/ui/Card/Card";
import { registerAction } from "./actions";
import { useState } from "react";
import { Button } from "@/components/forms/Button/Button";
import { Field } from "@/components/forms/Field/Field";
import { FaDoorOpen } from "react-icons/fa6";

export default function RegisterPage() {
  const [error, setError] = useState("");

  async function action(formData: FormData) {
    try {
      await registerAction(formData);
    } catch {
      setError("Ошибка регистрации");
    }
  }

  return (

    <div className="flex flex-col items-center justify-center min-h-screen bg-(--background)">
      <Card className="md:min-w-1/3">
        <form action={action}>
          <div className="flex flex-row justify-center items-center gap-2 p-2">
            <h1 className="text-center text-4xl text-sky-500 font-bold">Регистрация</h1>
            <FaDoorOpen className="text-4xl text-sky-500"></FaDoorOpen>
          </div>


          {/* <input name="name" placeholder="Имя" required />
          <input name="email" type="email" placeholder="Email" required />
          <input name="password" type="password" placeholder="Пароль" required /> */}

          <fieldset className="border-2 my-2 p-2 rounded-md border-slate-800 border-dashed">
            <legend className="italic text-slate-400 px-2" >Введите ваши данные для регистрации</legend>
            <Field name="name" type="text" placeholder="Имя" required></Field>
            <Field name="email" type="email" placeholder="Email" required></Field>
            <Field name="password" type="password" placeholder="Пароль" required></Field>
          </fieldset>

          <div className="flex flex-row justify-between">
            <Button type="submit">
              Зарегистрироваться
            </Button>
          </div>

          {error && <p style={{ color: "red" }}>{error}</p>}
        </form>
      </Card>
    </div>

  );
}
