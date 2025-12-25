"use server";

import { redirect } from "next/navigation";

const API_URL = "http://localhost:8080";

export async function registerAction(formData: FormData) {
  const name = formData.get("name");
  const email = formData.get("email");
  const password = formData.get("password");

  const res = await fetch(`${API_URL}/api/v1/auth/register`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ name, email, password }),
  });

  if (!res.ok) {
    throw new Error("Ошибка регистрации");
  }

  redirect("/login");
}
