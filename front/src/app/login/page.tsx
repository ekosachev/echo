"use client";

import { Geist, Geist_Mono } from "next/font/google";
import { loginAction } from "./actions";
import { useState } from "react";
import { Card } from "@/components/ui/Card/Card";
import { Field } from "@/components/forms/Field/Field";
import { Button } from "@/components/forms/Button/Button";
import { FaDoorOpen } from "react-icons/fa6";

export default function LoginPage() {
  const [error, setError] = useState("");

  async function action(formData: FormData) {
    try {
      await loginAction(formData);
    } catch (e) {
      setError("Ошибка входа");
    }
  }



  return (

    <div className="flex flex-col items-center justify-center min-h-screen bg-(--background)">
      <Card className="md:min-w-1/3">
        <form action={action}>
          <div className="flex flex-row justify-center items-center gap-2 p-2">
            <h1 className="text-center text-4xl text-sky-500 font-bold">Вход</h1>
            <FaDoorOpen className="text-4xl text-sky-500"></FaDoorOpen>
          </div>
          
          <fieldset className="border-2 my-2 p-2 rounded-md border-slate-800 border-dashed">
            <legend className="italic text-slate-400 px-2" >Введите ваши данные для входа</legend>
            <Field name="email" type="email" placeholder="Email" required></Field>
            <Field name="password" type="password" placeholder="Пароль" required></Field>
          </fieldset>

          <div className="flex flex-row justify-between">
            <Button type="submit">Войти</Button>
            <a href="/register" className="text-sky-500 hover:underline self-center">Регистрация</a>
          </div>

          {error && <p style={{ color: "red" }}>{error}</p>}
        </form>
      </Card>
    </div>

  );
}
