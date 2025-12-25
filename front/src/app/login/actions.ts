"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";

const API_URL = "http://localhost:8080";

export async function loginAction(formData: FormData) {
  const email = formData.get("email");
  const password = formData.get("password");

  const res = await fetch(`${API_URL}/api/v1/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });

  if (!res.ok) {
    throw new Error("Неверный email или пароль");
  }

  const data = await res.json();

  (await cookies()).set("token", data.token, {
    httpOnly: true,
    secure: false, // true в prod
    sameSite: "lax",
    path: "/",
  });

  redirect("/dashboard");
}
