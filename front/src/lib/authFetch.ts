import { cookies } from "next/headers";
import { redirect } from "next/navigation";

const API_URL = "http://localhost:8080";

export async function authFetch(
  path: string,
  init: RequestInit = {}
) {
  const token = (await cookies()).get("token")?.value;

  if (!token) {
    redirect("/login");
  }

  const res = await fetch(`${API_URL}${path}`, {
    ...init,
    headers: {
      ...(init.headers || {}),
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  if (res.status === 401) {
    redirect("/login");
  }

  return res;
}
